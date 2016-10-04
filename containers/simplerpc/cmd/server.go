/*
A simple RPC service that puts simple things
in and out of etcd
*/
package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	"github.com/apourchet/hermes"
	"github.com/apourchet/incipit/lib/etcd"
	"github.com/apourchet/incipit/lib/healthz"
	"github.com/apourchet/incipit/lib/logging"
	"github.com/gin-gonic/gin"

	. "github.com/apourchet/incipit/lib/simplerpc"
)

type SimpleRpc struct {
	hermes.Serviceable
	kapi etcd.Client
}

func NewService(kapi etcd.Client) *SimpleRpc {
	return &SimpleRpc{ServiceDefinition, kapi}
}

func (s *SimpleRpc) PutKey(c *gin.Context, in *PutKeyIn, out *PutKeyOut) (int, error) {
	err := s.kapi.Set(context.Background(), in.Key, in.Value, nil)
	if err != nil {
		logging.Error("Failed to insert into etcd: %v", err)
		return http.StatusInternalServerError, err
	}
	logging.Info("Successfully inserted into etcd")
	out.Ok = true
	return http.StatusOK, nil
}

func (s *SimpleRpc) GetKey(c *gin.Context, in *GetKeyIn, out *GetKeyOut) (int, error) {
	val, err := s.kapi.Get(context.Background(), in.Key, nil)
	if err != nil {
		logging.Error("Failed to get from etcd: %v", err)
		return http.StatusInternalServerError, err
	}
	logging.Info("Successfully retrieved from etcd")
	out.Value = val
	return http.StatusOK, nil
}

func main() {
	kapi := etcd.GetK8sDefaultClientFatal()
	healthz.SpawnHealthCheck(healthz.DefaultPort)
	svc := NewService(kapi)

	engine := gin.New()
	hermes.InitService(svc).Serve(engine)
	engine.Run(fmt.Sprintf(":%d", RpcPort))
}
