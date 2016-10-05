package healthz

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPort     = 10000
	DefaultEndpoint = "/healthz"
)

func healthCheck(ctx *gin.Context) {
	ctx.Status(200)
}

// Spawns a server that returns 200 on host:{port}/{base}/healthz
func SpawnHealthCheck(port int) {
	engine := gin.New()
	engine.GET(DefaultEndpoint, healthCheck)
	go engine.Run(fmt.Sprintf(":%d", port))

	// TODO LOG
	// logging.Info("Healthz Ready at :%d%s\n", port, DefaultEndpoint)
}
