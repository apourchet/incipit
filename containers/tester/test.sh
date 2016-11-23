#! /bin/bash

echo "TESTING"

go test github.com/apourchet/incipit/lib/... || go test -v github.com/apourchet/incipit/lib/... -args -alsologtostderr

exit 0
