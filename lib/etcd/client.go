package etcd

import (
	"time"

	"github.com/apourchet/dummy/lib/logging"
	"github.com/apourchet/dummy/lib/utils"
	etcdClient "github.com/coreos/etcd/client"
)

const (
	CLIENT_SERVICENAME = "etcd-client"
)

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
func GetK8sDefaultClient() (etcdClient.KeysAPI, error) {
	config := GetK8sDefaultConfig()
	client, err := etcdClient.New(config)
	return etcdClient.NewKeysAPI(client), err
}

// Same as GetK8sDefaultClient but log.Fatal
// if there was an error
func GetK8sDefaultClientFatal() etcdClient.KeysAPI {
	c, err := GetK8sDefaultClient()
	logging.FatalError(err)
	return c
}
