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
import random

# Clear existing database
open("database.csv", "w").close()

# Generate and write all records at once
with open("database.csv", "w") as f:
    for i in range(1, $NUM_RECORDS + 1):
        key = f"key_{i}"
        value = json.dumps({"id": i, "name": f"item_{i}"})
        f.write(f"{key},{value}\n")

        # Randomly write duplicate keys (70% chance)
        if random.random() < 0.7:
            dup_id = random.randint(1, i)
            dup_key = f"key_{dup_id}"
            update_count = random.randint(1, 5)
            dup_value = json.dumps({"id": dup_id, "name": f"item_{dup_id}_updated_{update_count}_times"})
            f.write(f"{dup_key},{dup_value}\n")

print(f"Created $NUM_RECORDS records with duplicates in database.csv")
EOF
