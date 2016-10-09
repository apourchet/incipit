package redis

import (
	"time"

	"github.com/apourchet/incipit/lib/utils"
)

type RedisClient interface {
	Get(key string) (value string, found bool, err error)
	Set(key string, value string) error
	SetExpire(key string, value string, exp time.Duration) error
	Delete(key string) error
}

func GetDefaultClient() RedisClient {
	if utils.InKubernetes() {
		return GetK8sClientV1()
	}
	return NewMockRedisClient()
}
