#! /bin/bash

echo "TESTING"

go test -v github.com/apourchet/incipit/lib/...
# for file in /tests/*; do
#     echo ">>> $file" && $file
# done

exit 0
