#!/bin/bash
if [[ $# -ne 1 ]]; then
    echo 'Wrong arguments, expecting only one (reqs/s)'
    exit 1
fi

TARGET_FILE="targets.txt"
DURATION="30s"  # Duration of the test, e.g., 60s for 60 seconds
RATE=$1    # Number of requests per second
RESULTS_FILE="results_$RATE.bin"
REPORT_FILE="report_$RATE.txt"
# Set the endpoint to the location of your API
ENDPOINT="http://$SERVER_API_IP:3000/user"

echo "SERVER IP: $SERVER_API_IP"


# Check if Vegeta is installed
if ! command -v vegeta &> /dev/null
then
    echo "Vegeta could not be found, please install it."
    exit 1
fi

> "$TARGET_FILE"  # Clear the file if it already exists

# Assuming body.json exists and contains the correct JSON structure for the POST request
for i in $(seq 1 $RATE); do 
    echo "POST $ENDPOINT" >> "$TARGET_FILE"
    echo "Content-Type: application/json" >> "$TARGET_FILE"
    echo "@body.json" >> "$TARGET_FILE"
    echo "" >> "$TARGET_FILE"
done

echo "Starting Vegeta attack for $DURATION at $RATE requests per second..."
# Run the attack and save the results to a binary file
vegeta attack -rate=$RATE -duration=$DURATION -targets="$TARGET_FILE" > "$RESULTS_FILE"

echo "Load test finished, generating reports..."
# Generate a textual report from the binary results file
vegeta report -type=text "$RESULTS_FILE" > "$REPORT_FILE"
echo "Textual report generated: $REPORT_FILE"

# Jump lines
echo -e 
echo -e

cat $REPORT_FILE
