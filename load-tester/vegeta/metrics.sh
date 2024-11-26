#!/bin/bash
if [[ $# -ne 1 ]]; then
    echo 'Wrong arguments, expecting only one (reqs/s)'
    exit 1
fi

# Check if Vegeta is installed
if ! command -v vegeta &> /dev/null
then
    echo "Vegeta could not be found, please install it."
    exit 1
fi


TARGET_FILE="targets.txt"
DURATION="10s"  # Duration of the test, e.g., 60s for 60 seconds
RATE=$1    # Number of requests per second
RESULTS_FILE="results_$RATE.bin"
METRICS_FILE="metrics_$RATE.txt"
REPORT_FILE="report_$RATE.txt"
ENDPOINT="http://$SERVER_API_IP:3000/user"

> "$TARGET_FILE"  # Clear the file if it already exists

for i in $(seq 1 $RATE); do 
    echo "POST $ENDPOINT" >> "$TARGET_FILE"
    echo "Content-Type: application/json" >> "$TARGET_FILE"
    echo "@body.json" >> "$TARGET_FILE"
    echo "" >> "$TARGET_FILE"
done

echo "Starting Vegeta attack for $DURATION at $RATE requests per second..."
vegeta attack -rate=$RATE -duration=$DURATION -targets="$TARGET_FILE" | tee $RESULTS_FILE | vegeta encode > $METRICS_FILE
#
# Generate a textual report from the binary results file
vegeta report -type=text "$RESULTS_FILE" > "$REPORT_FILE"
echo "Textual report generated: $REPORT_FILE"

# Jump lines
echo -e 
echo -e

cat $REPORT_FILE
