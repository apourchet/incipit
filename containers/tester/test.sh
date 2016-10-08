#! /bin/bash

echo "TESTING"

for file in /tests/*; do
    $file
done

exit 0
