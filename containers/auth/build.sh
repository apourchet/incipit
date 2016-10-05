#! /bin/sh

CGO_ENABLED=0 go build -i -ldflags "-s" -installsuffix cgo -o ./auth/build/auth ./auth/server.go
