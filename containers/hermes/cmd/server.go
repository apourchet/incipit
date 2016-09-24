/*
A simple RPC service that puts simple things
in and out of etcd
*/
package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"

	etcd_client "github.com/coreos/etcd/client"

	"github.com/apourchet/dummy/lib/etcd"
	"github.com/apourchet/dummy/lib/healthz"
	"github.com/apourchet/dummy/lib/utils"
	"github.com/gin-gonic/gin"
)

const (
	RpcPort = 8080
)

var (
	kapi etcd_client.KeysAPI
)

func helloPost(c *gin.Context) {
	_, err := kapi.Set(context.Background(), "/foo", "bar", nil)
	if err != nil {
		utils.Error("Failed to insert into etcd: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Info("Successfully inserted into etcd")
	c.JSON(http.StatusOK, gin.H{"message": "Hello Hermes!"})
}

func helloGet(c *gin.Context) {
	resp, err := kapi.Get(context.Background(), "/foo", nil)
	if err != nil {
		utils.Error("Failed to get from etcd: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.Info("Successfully retrieved from etcd")
	c.JSON(http.StatusOK, gin.H{"value": fmt.Sprintf("%q", resp.Node.Value)})
}

func main() {
	kapi = etcd.GetK8sDefaultClientFatal()
	healthz.SpawnHealthCheck(healthz.DefaultPort)

	engine := gin.New()
	engine.POST("/rpc/v1/hermes", helloPost)
	engine.GET("/rpc/v1/hermes", helloGet)
	engine.Run(fmt.Sprintf(":%d", RpcPort))
}
