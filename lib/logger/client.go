package logger

import (
	"github.com/apourchet/incipit/lib/utils"
	protos "github.com/apourchet/incipit/protos/go"
	"github.com/golang/glog"
	"google.golang.org/grpc"
)

const (
	ClientServiceName = "logger"
)

func GetK8sDefaultClient() (protos.LoggerClient, error) {
	conn, err := grpc.Dial(utils.GetK8sAddress(ClientServiceName), grpc.WithInsecure())
	if err != nil {
		glog.Errorf("Failed to dial GRPC: %v", err)
		return nil, err
	}
	return protos.NewLoggerClient(conn), err
}

func GetDefaultClient() (protos.LoggerClient, error) {
	if utils.InKubernetes() {
		return GetK8sDefaultClient()
	}
	return NewMockLoggerClient(), nil
}
