package etcd

import (
	"testing"

	"github.com/apourchet/incipit/lib/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func etcdClients() []EtcdClient {
	if utils.InKubernetes() {
		return []EtcdClient{
			NewMockEtcdClient(),
			GetK8sClientV1Fatal(),
		}
	}
	return []EtcdClient{
		NewMockEtcdClient(),
	}
}

func TestSimpleSetGet(t *testing.T) {
	for _, c := range etcdClients() {
		k := "/" + uuid.NewV4().String()
		_, found, err := c.Get(k)
		assert.Nil(t, err)
		assert.Equal(t, found, false)
		err = c.Set(k, "value")
		assert.Nil(t, err)
		val, found, err := c.Get(k)
		assert.Nil(t, err)
		assert.Equal(t, found, true)
		assert.Equal(t, val, "value")
	}
}

func TestDelete(t *testing.T) {
	for _, c := range etcdClients() {
		k := "/" + uuid.NewV4().String()
		err := c.Delete(k)
		assert.Nil(t, err)

		err = c.Set(k, "value")
		assert.Nil(t, err)

		val, found, err := c.Get(k)
		assert.Nil(t, err)
		assert.Equal(t, found, true)
		assert.Equal(t, val, "value")

		err = c.Delete(k)
		assert.Nil(t, err)

		val, found, err = c.Get(k)
		assert.Nil(t, err)
		assert.Equal(t, found, false)
	}
}
