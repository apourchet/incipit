package main

import (
	"flag"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	"github.com/apourchet/incipit/lib/healthz"
	"github.com/apourchet/incipit/lib/utils"
	pb "github.com/apourchet/incipit/protos/go"
	"github.com/golang/glog"
)

type loggerServer struct{}

func (l *loggerServer) LogLogin(ctx context.Context, req *pb.LogLoginReq) (*pb.LogLoginRes, error) {
	return nil, nil
}

func main() {
	flag.Parse()
	healthz.SpawnHealthCheck(healthz.DefaultPort)

	server := grpc.NewServer()
	pb.RegisterLoggerServer(server, &loggerServer{})

	glog.Fatalf("LoggerServer exited with error: %v", utils.ServeGRPC(10002, server))
}
