#!/bin/bash

# Populate database with sample records
# Usage: ./populate.sh <number_of_records>

if [ $# -eq 0 ]; then
    echo "Usage: $0 <number_of_records>"
    exit 1
fi

NUM_RECORDS=$1

python3 << EOF
import json

# Clear existing database
open("database.csv", "w").close()

# Generate and write all records at once
with open("database.csv", "w") as f:
    for i in range(1, $NUM_RECORDS + 1):
        key = f"key_{i}"
        value = json.dumps({"id": i, "name": f"item_{i}"})
        f.write(f"{key},{value}\n")

print(f"Created $NUM_RECORDS records in database.csv")
EOF
