package simplerpc

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/apourchet/hermes"
	"github.com/apourchet/incipit/lib/etcd"
)

type SimpleRpcMock struct {
	hermes.Mockable
	kapi etcd.Client
}

func NewSimpleRpcMock() hermes.Callable {
	mocked := &SimpleRpcMock{ServiceDefinition, etcd.NewMockClient()}
	return hermes.InitMockService(mocked)
}

func (s *SimpleRpcMock) PutKey(in *PutKeyIn, out *PutKeyOut) (int, error) {
	err := s.kapi.Set(context.Background(), in.Key, in.Value, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	out.Ok = true
	return http.StatusOK, nil
}

func (s *SimpleRpcMock) GetKey(in *GetKeyIn, out *GetKeyOut) (int, error) {
	val, err := s.kapi.Get(context.Background(), in.Key, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	out.Value = val
	return http.StatusOK, nil
}
