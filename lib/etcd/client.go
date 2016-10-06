package etcd

import (
	"fmt"
	"time"

	"github.com/apourchet/incipit/lib/utils"
	etcdClient "github.com/coreos/etcd/client"
	"github.com/golang/glog"
	"golang.org/x/net/context"
)

const (
	CLIENT_SERVICENAME = "etcd-client"
)

// Get the address of the etcd client in the k8s cluster
// from environment variables
func GetClientK8sAddress() string {
	return utils.GetK8sAddress(CLIENT_SERVICENAME)
}

// Gets etcd endpoints based on environment variables
// set by kubernetes
func GetK8sDefaultConfig() etcdClient.Config {
	return etcdClient.Config{
		Endpoints:               []string{"http://" + GetClientK8sAddress()},
		Transport:               etcdClient.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
}

// Returns an etcd client based on K8sDefaultConfig
func GetK8sDefaultClient() (Client, error) {
	config := GetK8sDefaultConfig()
	client, err := etcdClient.New(config)
	kapi := etcdClient.NewKeysAPI(client)
	return NewClient(kapi), err
}

// Same as GetK8sDefaultClient but log.Fatal
// if there was an error
func GetK8sDefaultClientFatal() Client {
	c, err := GetK8sDefaultClient()
	if err != nil {
		glog.Fatal(err)
	}
	return c
}

type Client interface {
	Get(ctx context.Context, key string, opts *etcdClient.GetOptions) (string, error)
	Set(ctx context.Context, key, value string, opts *etcdClient.SetOptions) error

	// TODO Implement those
	// Delete(ctx context.Context, key string, opts *etcdClient.DeleteOptions) (string, error)
	// Update(ctx context.Context, key, value string) (string, error)
	// Create(ctx context.Context, key, value string) (string, error)
	// CreateInOrder(ctx context.Context, dir, value string, opts *etcdClient.CreateInOrderOptions) (string, error)
}

type ClientImpl struct {
	kapi etcdClient.KeysAPI
}

func NewClient(kapi etcdClient.KeysAPI) Client {
	return &ClientImpl{kapi}
}

func (c *ClientImpl) Get(ctx context.Context, key string, opts *etcdClient.GetOptions) (string, error) {
	resp, err := c.kapi.Get(ctx, key, opts)
	if err != nil {
		return "", err
	}
	return resp.Node.Value, nil
}

func (c *ClientImpl) Set(ctx context.Context, key, value string, opts *etcdClient.SetOptions) error {
	_, err := c.kapi.Set(ctx, key, value, opts)
	return err
}

// Mock of client
type MockClient struct {
	values map[string]string
}

func NewMockClient() Client {
	return &MockClient{make(map[string]string)}
}

func (c *MockClient) Get(ctx context.Context, key string, opts *etcdClient.GetOptions) (string, error) {
	val, ok := c.values[key]
	if !ok {
		return val, etcdClient.Error{
			etcdClient.ErrorCodeKeyNotFound,
			fmt.Sprintf("Could not find key '%s'", key),
			fmt.Sprintf("Key '%s' not set", key),
			0,
		}
	}
	return val, nil
}

func (c *MockClient) Set(ctx context.Context, key, value string, opts *etcdClient.SetOptions) error {
	c.values[key] = value
	return nil
}
