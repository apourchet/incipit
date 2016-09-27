/*
A simple RPC service that puts simple things
in and out of etcd
*/
package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	etcd_client "github.com/coreos/etcd/client"

	"github.com/apourchet/dummy/lib/etcd"
	"github.com/apourchet/dummy/lib/healthz"
	"github.com/apourchet/dummy/lib/utils"
	"github.com/apourchet/hermes"
	"github.com/gin-gonic/gin"
)

const (
	ServiceName = "simplerpc"
	RpcPort     = 8080
)

var (
	kapi etcd_client.KeysAPI
)

type SimpleRpc struct{}

// PutKey
type PutKeyIn struct {
	Key   string
	Value string
}
type PutKeyOut struct {
	Ok bool
}

func NewPutKeyIn() interface{}  { return &PutKeyIn{} }
func NewPutKeyOut() interface{} { return &PutKeyOut{} }

func (s *SimpleRpc) PutKey(c *gin.Context, in *PutKeyIn, out *PutKeyOut) (int, error) {
	_, err := kapi.Set(context.Background(), in.Key, in.Value, nil)
	if err != nil {
		utils.Error("Failed to insert into etcd: %v", err)
		return http.StatusInternalServerError, err
	}
	utils.Info("Successfully inserted into etcd")
	out.Ok = true
	return http.StatusOK, nil
}

// GetKey
type GetKeyIn struct {
	Key string
}
type GetKeyOut struct {
	Value string
}

func NewGetKeyIn() interface{}  { return &GetKeyIn{} }
func NewGetKeyOut() interface{} { return &GetKeyOut{} }

func (s *SimpleRpc) GetKey(c *gin.Context, in *GetKeyIn, out *GetKeyOut) (int, error) {
	resp, err := kapi.Get(context.Background(), in.Key, nil)
	if err != nil {
		utils.Error("Failed to get from etcd: %v", err)
		return http.StatusInternalServerError, err
	}
	utils.Info("Successfully retrieved from etcd")
	out.Value = resp.Node.Value
	return http.StatusOK, nil
}

// Implement hermes.Serviceable
func (s *SimpleRpc) Host() string {
	return utils.GetK8sAddress(ServiceName)
}

func (s *SimpleRpc) Endpoints() hermes.EndpointMap {
	return hermes.EndpointMap{
		hermes.Endpoint{"PutKey", "POST", "/rpc/v1/simplerpc", NewPutKeyIn, NewPutKeyOut},
		hermes.Endpoint{"GetKey", "GET", "/rpc/v1/simplerpc", NewGetKeyIn, NewGetKeyOut},
	}
}

func main() {
	kapi = etcd.GetK8sDefaultClientFatal()
	healthz.SpawnHealthCheck(healthz.DefaultPort)

	engine := gin.New()
	hermes.InitService(&SimpleRpc{}).Serve(engine)
	engine.Run(fmt.Sprintf(":%d", RpcPort))
}
