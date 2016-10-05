#! /bin/sh

CGO_ENABLED=0 go build -i -ldflags "-s" -o ./build/auth ./server.go
