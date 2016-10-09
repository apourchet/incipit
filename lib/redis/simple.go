package redis

import (
	"time"

	"github.com/apourchet/incipit/lib/utils"
	"github.com/golang/glog"
	redis "gopkg.in/redis.v4"
)

const (
	CLIENT_SERVICENAME = "redis-client"
)

func GetK8sDefaultConfig() *redis.Options {
	return &redis.Options{
		Addr:     utils.GetK8sAddress(CLIENT_SERVICENAME),
		Password: "", // no password set
		DB:       0,  // use default DB
	}
}

func GetK8sDefaultClient() RedisClient {
	glog.Infof("Using redis K8sDefaultClient")
	config := GetK8sDefaultConfig()
	return NewK8sClient(redis.NewClient(config))
}

type k8sClientImpl struct {
	client *redis.Client
}

func NewK8sClient(client *redis.Client) RedisClient {
	c := &k8sClientImpl{}
	c.client = client
	return c
}

func (c *k8sClientImpl) Get(key string) (value string, found bool, err error) {
	return "", false, err
}

func (c *k8sClientImpl) Set(key string, value string) error {
	return nil
}

func (c *k8sClientImpl) SetExpire(key string, value string, exp time.Duration) error {
	return nil
}

func (c *k8sClientImpl) Delete(key string) error {
	return nil
}
