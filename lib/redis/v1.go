package redis

import (
	"time"

	"github.com/apourchet/incipit/lib/utils"
	"github.com/golang/glog"
	redis "gopkg.in/redis.v4"
)

const (
	CLIENT_SERVICENAME = "redis"
)

func GetK8sConfigV1() *redis.Options {
	return &redis.Options{
		Addr:     utils.GetK8sAddress(CLIENT_SERVICENAME),
		Password: "", // no password set
		DB:       0,  // use default DB
	}
}

func GetK8sClientV1() RedisClient {
	glog.Infof("Using redis K8sClientV1")
	config := GetK8sConfigV1()
	return NewK8sClientV1(redis.NewClient(config))
}

type k8sClientV1 struct {
	client *redis.Client
}

func NewK8sClientV1(client *redis.Client) RedisClient {
	c := &k8sClientV1{}
	c.client = client
	return c
}

func (c *k8sClientV1) Get(key string) (value string, found bool, err error) {
	found, err = c.client.Exists(key).Result()
	if err != nil || !found {
		return "", found, err
	}
	val, err := c.client.Get(key).Result()
	return val, true, err
}

func (c *k8sClientV1) Set(key string, value string) error {
	return c.client.Set(key, value, 0).Err()
}

func (c *k8sClientV1) SetExpire(key string, value string, exp time.Duration) error {
	return c.client.Set(key, value, exp).Err()
}

func (c *k8sClientV1) Delete(key string) error {
	return c.client.Del(key).Err()
}
