package hermes

import (
	"fmt"

	"github.com/apourchet/dummy/lib/utils"
)

const (
	ServiceName = "hermes"
)

type Hermes interface {
	GetHello(key string) (string, error)
	PostHello(key, value string) error
}

type HelloRequest struct {
	Key string `json:"key"`
}

type HelloResponse struct {
	Value string `json:"value"`
}

type DefaultClient struct{ host string }

func NewDefaultClient() Hermes {
	return &DefaultClient{
		utils.GetK8sAddress(ServiceName),
	}
}

func (c *DefaultClient) GetHello(key string) (string, error) {
	url := fmt.Sprintf("http://%s/%s", c.host, "rpc/v1/hermes")
	data := HelloRequest{key}
	var resp HelloResponse
	err := utils.GetRpc(url, data, &resp)
	return resp.Value, err
}

func (c *DefaultClient) PostHello(key, value string) error {
	return nil
}
