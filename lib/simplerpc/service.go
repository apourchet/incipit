package simplerpc

import (
	"github.com/apourchet/hermes"
	"github.com/apourchet/incipit/lib/utils"
)

const (
	ServiceName = "simplerpc"
	RpcPort     = 8080
)

// Implements hermes.Serviceable
type serviceDefinition struct{}

var ServiceDefinition = &serviceDefinition{}

// Implementation of the definition
func (s *serviceDefinition) Host() string {
	return utils.GetK8sAddress(ServiceName)
}

func (s *serviceDefinition) Endpoints() hermes.EndpointMap {
	return hermes.EndpointMap{
		hermes.Endpoint{"PutKey", "POST", "/rpc/v1/simplerpc", NewPutKeyIn, NewPutKeyOut},
		hermes.Endpoint{"GetKey", "GET", "/rpc/v1/simplerpc", NewGetKeyIn, NewGetKeyOut},
	}
}

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
