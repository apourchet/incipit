package etcd

import (
	"time"

	"github.com/apourchet/dummy/lib/utils"
	etcdClient "github.com/coreos/etcd/client"
)

const (
	CLIENT_SERVICENAME = "etcd-client"
)

func GetClientK8sAddress() string {
	return utils.GetK8sAddress(CLIENT_SERVICENAME)
}

func GetK8sDefaultConfig() etcdClient.Config {
	return etcdClient.Config{
		Endpoints:               []string{"http://" + GetClientK8sAddress()},
		Transport:               etcdClient.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
}

func GetK8sDefaultClient() (etcdClient.KeysAPI, error) {
	config := GetK8sDefaultConfig()
	client, err := etcdClient.New(config)
	return etcdClient.NewKeysAPI(client), err
}
