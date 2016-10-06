package etcd

import (
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
func GetK8sDefaultClient() (EtcdClient, error) {
	config := GetK8sDefaultConfig()
	client, err := etcdClient.New(config)
	kapi := etcdClient.NewKeysAPI(client)
	return NewK8sClient(kapi), err
}

// Same as GetK8sDefaultClient but log.Fatal
// if there was an error
func GetK8sDefaultClientFatal() EtcdClient {
	c, err := GetK8sDefaultClient()
	if err != nil {
		glog.Fatal(err)
	}
	return c
}

type K8sClientImpl struct {
	kapi etcdClient.KeysAPI
}

func NewK8sClient(kapi etcdClient.KeysAPI) EtcdClient {
	return &K8sClientImpl{kapi}
}

func (c *K8sClientImpl) Get(key string) (string, bool, error) {
	resp, err := c.kapi.Get(context.Background(), key, nil)
	if err != nil {
		casted, ok := err.(etcdClient.Error)
		if ok && casted.Code == etcdClient.ErrorCodeKeyNotFound {
			return "", false, nil
		}
		return "", false, err
	}
	return resp.Node.Value, true, nil
}

func (c *K8sClientImpl) Set(key string, value string) error {
	_, err := c.kapi.Set(context.Background(), key, value, nil)
	return err
}

func (c *K8sClientImpl) Delete(key string) error {
	_, err := c.kapi.Delete(context.Background(), key, nil)
	if err != nil {
		if casted, ok := err.(etcdClient.Error); ok && casted.Code == etcdClient.ErrorCodeKeyNotFound {
			return nil
		}
	}
	return err
}
