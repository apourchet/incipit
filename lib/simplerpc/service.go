package simplerpc

import (
	"github.com/apourchet/hermes"
	"github.com/apourchet/incipit/lib/utils"
	"github.com/gin-gonic/gin"
)

const (
	ServiceName = "simplerpc"
	RpcPort     = 8080
)

type SimpleRpcService interface {
	PutKey(*gin.Context, *PutKeyIn, *PutKeyOut) (int, error)
	GetKey(*gin.Context, *GetKeyIn, *GetKeyOut) (int, error)
}

type ServiceDefinition struct{}

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

// GetKey
type GetKeyIn struct {
	Key string
}
type GetKeyOut struct {
	Value string
}

func NewGetKeyIn() interface{}  { return &GetKeyIn{} }
func NewGetKeyOut() interface{} { return &GetKeyOut{} }

// Implementation
func (s *ServiceDefinition) Host() string {
	return utils.GetK8sAddress(ServiceName)
}

func (s *ServiceDefinition) Endpoints() hermes.EndpointMap {
	return hermes.EndpointMap{
		hermes.Endpoint{"PutKey", "POST", "/rpc/v1/simplerpc", NewPutKeyIn, NewPutKeyOut},
		hermes.Endpoint{"GetKey", "GET", "/rpc/v1/simplerpc", NewGetKeyIn, NewGetKeyOut},
	}
}
