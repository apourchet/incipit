package credentials

import (
	"fmt"
	"strings"

	"github.com/apourchet/incipit/lib/etcd"
	"github.com/golang/glog"
	uuid "github.com/satori/go.uuid"
)

type etcdCredentialClient struct {
	store etcd.EtcdClient
}

func NewEtcdCredentialsClient() (CredentialsClient, error) {
	glog.Infof("Creds: Using EtcdCredentialClient")
	c := &etcdCredentialClient{}
	c.store = etcd.GetDefaultClient()
	return c, nil
}

func NewCredentialsClient(store etcd.EtcdClient) CredentialsClient {
	glog.Infof("Using EtcdCredentialClient2")
	c := &etcdCredentialClient{}
	c.store = store
	return c
}

func (c *etcdCredentialClient) UserExists(key string) (exists bool, err error) {
	glog.Infof("Creds: UserExists: %s", key)
	storeKey := fmt.Sprintf("/creds/%s", key)
	_, found, err := c.store.Get(storeKey)
	if err != nil {
		return true, err
	}
	return found, nil
}

func (c *etcdCredentialClient) Register(key string, pass string) (ok bool, err error) {
	glog.Infof("Creds: Register %s", key)

	exists, err := c.UserExists(key)
	if err != nil || exists {
		return !exists, err
	}

	storeKey := fmt.Sprintf("/creds/%s", key)
	id := uuid.NewV4().String()
	storeVal := fmt.Sprintf("%s;%s", id, pass)

	err = c.store.Set(storeKey, storeVal)
	if err != nil {
		glog.Errorf("Failed to set key in etcd: %v", err)
		return false, err
	}

	glog.Infof("Creds: Successfully set key in etcd: %s", storeKey)
	return true, nil
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
	glog.Infof("Creds: CheckCredentials %s", key)
	storeKey := fmt.Sprintf("/creds/%s", key)
	val, found, err := c.store.Get(storeKey)
	if err != nil {
		glog.Errorf("Creds: Failed to check credentials: %v", err)
		return false, err
	}
	if !found {
		glog.Infof("Creds: Failed to find user: %v", err)
		return false, nil
	}
	arr := strings.Split(val, ";")
	if len(arr) < 2 {
		return false, fmt.Errorf("Account invalid")
	}
	b := arr[1] == pass
	glog.Infof("Creds: Credential match: %t", b)
	return b, nil
}

func (c *etcdCredentialClient) Deregister(key string) error {
	glog.Infof("Creds: Deregister %s", key)
	storeKey := fmt.Sprintf("/creds/%s", key)
	err := c.store.Delete(storeKey)
	return err
}
