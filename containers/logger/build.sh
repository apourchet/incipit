#! /bin/sh

CGO_ENABLED=0 go build -i -ldflags "-s" -o ./bin/logger ./server.go
echo "Built logger"
