#! /bin/bash

echo "TESTING"

for file in /tests/*; do
    echo ">>> $file" && $file
done

exit 0
