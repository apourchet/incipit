package etcd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func etcdClients() []EtcdClient {
	return []EtcdClient{
		NewMockEtcdClient(),
	}
}

func TestSimpleSetGet(t *testing.T) {
	for _, c := range etcdClients() {
		_, found, err := c.Get("/key")
		assert.Nil(t, err)
		assert.Equal(t, found, false)
		err = c.Set("/key", "value")
		assert.Nil(t, err)
		val, found, err := c.Get("/key")
		assert.Nil(t, err)
		assert.Equal(t, found, true)
		assert.Equal(t, val, "value")
	}
}

func TestDelete(t *testing.T) {
	for _, c := range etcdClients() {
		err := c.Delete("/key")
		assert.Nil(t, err)

		err = c.Set("/key", "value")
		assert.Nil(t, err)

		val, found, err := c.Get("/key")
		assert.Nil(t, err)
		assert.Equal(t, found, true)
		assert.Equal(t, val, "value")

		err = c.Delete("/key")
		assert.Nil(t, err)

		val, found, err = c.Get("/key")
		assert.Nil(t, err)
		assert.Equal(t, found, false)
	}
}
