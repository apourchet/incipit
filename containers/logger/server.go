package main

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	"github.com/apourchet/incipit/lib/etcd"
	"github.com/apourchet/incipit/lib/healthz"
	"github.com/apourchet/incipit/lib/utils"
	pb "github.com/apourchet/incipit/protos/go"
	"github.com/golang/glog"
)

type loggerService struct {
	store etcd.EtcdClient
}

func (l *loggerService) LogLogin(ctx context.Context, req *pb.LogLoginReq) (*pb.LogLoginRes, error) {
	key := req.UserId
	glog.Infof("LogLogin: %s", key)

	// Get the last login value
	var lastLogin int64
	val, found, err := l.store.Get(key)
	if err != nil {
		return nil, err
	}
	if found {
		lastLogin, err = strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
	}

	res := &pb.LogLoginRes{}
	res.LastLogin = int64(lastLogin)

	// Store the new
	newLogin := time.Now().Unix()
	return res, l.store.Set(key, fmt.Sprintf("%d", newLogin))
}

func main() {
	flag.Parse()
	go healthz.SpawnHealthCheck(healthz.DefaultPort)

	service := &loggerService{etcd.GetDefaultClient()}
	server := grpc.NewServer()
	pb.RegisterLoggerServer(server, service)

	go func() {
		err := utils.ServeGRPC(10001, server)
		glog.Fatalf("LoggerServer exited with error: %v", err)
	}()

	select {} // Sleep forever
}
