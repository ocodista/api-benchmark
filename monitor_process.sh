#!/bin/bash

# Check for correct number of arguments
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <PID> <Label>"
    exit 1
fi

PID=$1
LABEL=$2
FILE="process_stats_${LABEL}.csv"

# Add CSV header
echo "Timestamp, RAM (KB), CPU (%), Threads_Count, FD_Count" > $FILE

# Monitoring function
monitor_process() {
    while true; do
        TIMESTAMP=$(date +%Y-%m-%d\ %H:%M:%S)
        RAM=$(ps -p $PID -o rss=)
        CPU=$(ps -p $PID -o %cpu=)
        THREADS_COUNT=$(ps -p $PID -o nlwp=)
        FD_COUNT=$(ls /proc/$PID/fd | wc -l)

        # Append to CSV
        echo "$TIMESTAMP, $RAM, $CPU, $THREADS_COUNT, $FD_COUNT" >> $FILE

        # Update every second
        sleep 1
    done
}

# Check if process exists
if kill -0 $PID 2>/dev/null; then
    monitor_process
else
    echo "Process with PID $PID not found."
    exit 1
fi
