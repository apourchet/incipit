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
func GetClientK8sAddressV1() string {
	return utils.GetK8sAddress(CLIENT_SERVICENAME)
}

// Gets etcd endpoints based on environment variables
// set by kubernetes
func GetK8sConfigV1() etcdClient.Config {
	return etcdClient.Config{
		Endpoints:               []string{"http://" + GetClientK8sAddressV1()},
		Transport:               etcdClient.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
}

// Returns an etcd client based on K8sDefaultConfig
func GetK8sClientV1() (EtcdClient, error) {
	glog.Infof("Using etcd K8sDefaultClient")
	config := GetK8sConfigV1()
	client, err := etcdClient.New(config)
	kapi := etcdClient.NewKeysAPI(client)
	return NewK8sClientV1(kapi), err
}

// Same as GetK8sDefaultClient but log.Fatal
// if there was an error
func GetK8sClientV1Fatal() EtcdClient {
	c, err := GetK8sClientV1()
	if err != nil {
		glog.Fatal(err)
	}
	return c
}

type K8sClientV1 struct {
	kapi etcdClient.KeysAPI
}

func NewK8sClientV1(kapi etcdClient.KeysAPI) EtcdClient {
	return &K8sClientV1{kapi}
}

func (c *K8sClientV1) Get(key string) (string, bool, error) {
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

func (c *K8sClientV1) Set(key string, value string) error {
	_, err := c.kapi.Set(context.Background(), key, value, nil)
	return err
}

func (c *K8sClientV1) Delete(key string) error {
	_, err := c.kapi.Delete(context.Background(), key, nil)
	if err != nil {
		if casted, ok := err.(etcdClient.Error); ok && casted.Code == etcdClient.ErrorCodeKeyNotFound {
			return nil
		}
	}
	return err
}
