#! /bin/sh

CGO_ENABLED=0 go build -i -ldflags "-s" -o ./bin/auth ./*.go
echo "Built auth"
