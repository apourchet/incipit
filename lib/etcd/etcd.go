package etcd

import "github.com/apourchet/incipit/lib/utils"

type EtcdClient interface {
	// Error for internal only
	Get(key string) (value string, found bool, err error)

	// Error for internal only
	Set(key string, value string) error

	// Error for internal only
	Delete(key string) error
}

func GetDefaultClient() EtcdClient {
	if utils.InKubernetes() {
		return GetK8sClientV1Fatal()
	}
	return NewMockEtcdClient()
}
