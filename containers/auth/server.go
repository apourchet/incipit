package main

import (
	"flag"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/apourchet/incipit/lib/auth"
	"github.com/apourchet/incipit/lib/healthz"
	"github.com/apourchet/incipit/lib/logger"
	"github.com/apourchet/incipit/lib/utils"
	protos "github.com/apourchet/incipit/protos/go"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	go healthz.SpawnHealthCheck(healthz.DefaultPort)

	authClient := auth.GetDefaultClient()
	loggerClient, err := logger.GetDefaultClient()
	if err != nil {
		glog.Fatalf("Fatal error: %v", err)
	}
	service := NewAuthService(authClient, loggerClient)

	server := grpc.NewServer()
	protos.RegisterAuthServiceServer(server, service)

	go func() {
		err := utils.ServeGRPC(10001, server)
		glog.Fatalf("AuthServer exited with error: %v", err)
	}()

	go func() {
		err := utils.GrpcReverseProxy(context.Background(), protos.RegisterAuthServiceHandlerFromEndpoint, 10000, 10001)
		glog.Fatalf("AuthServer reverse-proxy exited with error: %v", err)
	}()

	select {} // Sleep forever
}
