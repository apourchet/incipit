package simplerpc

import "net/http"

type SimpleRpcMock struct {
	ServiceDefinition
}

func (s *SimpleRpcMock) PutKey(in *PutKeyIn, out *PutKeyOut) (int, error) {
	return http.StatusOK, nil
}

func (s *SimpleRpcMock) GetKey(in *GetKeyIn, out *GetKeyOut) (int, error) {
	return http.StatusOK, nil
}
