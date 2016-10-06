package settings

import (
	"testing"

	"github.com/apourchet/incipit/lib/etcd"
	"github.com/stretchr/testify/assert"
)

func settingsClients() []SettingsClient {
	return []SettingsClient{
		NewSimpleSettingsClient(etcd.NewMockEtcdClient()),
	}
}

func TestSimpleSetGet(t *testing.T) {
	for _, c := range settingsClients() {
		_, found, err := c.GetSetting("user1", "field1")
		assert.Nil(t, err)
		assert.Equal(t, found, false)
		err = c.SetSetting("user1", "field1", "value1")
		assert.Nil(t, err)
		val, found, err := c.GetSetting("user1", "field1")
		assert.Nil(t, err)
		assert.Equal(t, found, true)
		assert.Equal(t, val, "value1")
	}
}
