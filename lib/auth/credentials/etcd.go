package credentials

import (
	"fmt"
	"strings"

	"github.com/apourchet/incipit/lib/etcd"
	"github.com/golang/glog"
)

type etcdCredentialClient struct {
	store  etcd.EtcdClient
	currId int
}

func NewK8sEtcdCredentialsClient() (CredentialsClient, error) {
	glog.Infof("Using K8sEtcdCredentialClient")
	c := &etcdCredentialClient{}
	store, err := etcd.GetK8sDefaultClient()
	if err != nil {
		return nil, err
	}
	c.store = store
	return c, nil
}

func NewEtcdCredentialsClient(store etcd.EtcdClient) CredentialsClient {
	glog.Infof("Using etcdCredentialClient")
	c := &etcdCredentialClient{}
	c.store = store
	return c
}

func (c *etcdCredentialClient) UserExists(key string) (exists bool, err error) {
	storeKey := fmt.Sprintf("/creds/%s", key)
	_, found, err := c.store.Get(storeKey)
	if err != nil {
		return true, err
	}
	return found, nil
}

func (c *etcdCredentialClient) Register(key string, pass string) (err error) {
	storeKey := fmt.Sprintf("/creds/%s", key)
	storeVal := fmt.Sprintf("%d;%s", c.currId, pass)
	c.currId += 1
	return c.store.Set(storeKey, storeVal)
}

func (c *etcdCredentialClient) GetUserId(key string) (id string, err error) {
	storeKey := fmt.Sprintf("/creds/%s", key)
	val, found, err := c.store.Get(storeKey)
	if err != nil || !found {
		return "", err
	}
	arr := strings.Split(val, ";")
	return arr[0], nil
}

func (c *etcdCredentialClient) CheckCredentials(key, pass string) (check bool, err error) {
	storeKey := fmt.Sprintf("/creds/%s", key)
	val, found, err := c.store.Get(storeKey)
	if err != nil || !found {
		return false, err
	}
	arr := strings.Split(val, ";")
	return arr[1] == pass, nil
}

func (c *etcdCredentialClient) Deregister(key string) error {
	storeKey := fmt.Sprintf("/creds/%s", key)
	err := c.store.Delete(storeKey)
	return err
}
