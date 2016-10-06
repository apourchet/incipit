package settings

import (
	"fmt"

	"github.com/apourchet/incipit/lib/etcd"
)

type SettingsClient interface {
	SetSetting(userId string, field string, value string) error
	GetSetting(userId string, field string) (value string, found bool, err error)
}

type SimpleClient struct {
	store etcd.EtcdClient
}

func NewSimpleSettingsClient(store etcd.EtcdClient) SettingsClient {
	c := &SimpleClient{}
	c.store = store
	return c
}

func (c *SimpleClient) SetSetting(userid string, field string, value string) error {
	key := fmt.Sprintf("/settings/%s/%s", userid, field)
	return c.store.Set(key, value)
}

func (c *SimpleClient) GetSetting(userid string, field string) (value string, found bool, err error) {
	key := fmt.Sprintf("/settings/%s/%s", userid, field)
	return c.store.Get(key)
}
