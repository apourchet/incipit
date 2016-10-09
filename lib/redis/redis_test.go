package redis

import (
	"testing"
	"time"

	"github.com/apourchet/incipit/lib/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func redisClients(t *testing.T) []RedisClient {
	if utils.InKubernetes() {
		return []RedisClient{
			NewMockRedisClient(),
			GetK8sClientV1(),
		}
	}
	return []RedisClient{
		NewMockRedisClient(),
	}
}

func TestEmptyGet(t *testing.T) {
	for _, c := range redisClients(t) {
		k := uuid.NewV4().String()

		val, found, err := c.Get(k)
		assert.Nil(t, err)
		assert.False(t, found)
		assert.Equal(t, val, "")
	}
}

func TestSetGet(t *testing.T) {
	for _, c := range redisClients(t) {
		k := uuid.NewV4().String()

		err := c.Set(k, "value")
		assert.Nil(t, err)

		val, found, err := c.Get(k)
		assert.Nil(t, err)
		assert.True(t, found)
		assert.Equal(t, val, "value")
	}
}

func TestSetGetDeleteGet(t *testing.T) {
	for _, c := range redisClients(t) {
		k := uuid.NewV4().String()

		err := c.Set(k, "value")
		assert.Nil(t, err)

		val, found, err := c.Get(k)
		assert.Nil(t, err)
		assert.True(t, found)
		assert.Equal(t, val, "value")

		err = c.Delete(k)
		assert.Nil(t, err)

		_, found, err = c.Get(k)
		assert.Nil(t, err)
		assert.False(t, found)
	}
}

func TestSetExpire(t *testing.T) {
	for _, c := range redisClients(t) {
		k := uuid.NewV4().String()

		err := c.SetExpire(k, "value", 10*time.Millisecond)
		assert.Nil(t, err)

		val, found, err := c.Get(k)
		assert.Nil(t, err)
		assert.True(t, found)
		assert.Equal(t, val, "value")

		time.Sleep(20 * time.Millisecond)
		_, found, err = c.Get(k)
		assert.Nil(t, err)
		assert.False(t, found)
	}
}
