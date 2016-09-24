package healthz

import (
	"fmt"

	"github.com/apourchet/dummy/lib/utils"
	"github.com/gin-gonic/gin"
)

const (
	DEFAULT_PORT     = 10000
	DEFAULT_ENDPOINT = "/healthz"
)

func healthCheck(ctx *gin.Context) {
	ctx.Status(200)
}

// Spawns a server that returns 200 on host:{port}/{base}/healthz
func SpawnHealthCheck(base string, port int) {
	engine := gin.New()
	engine.GET(base+DEFAULT_ENDPOINT, healthCheck)
	go engine.Run(fmt.Sprintf(":%d", port))

	utils.Info("Healthz Ready at %s%s:%d\n", base, DEFAULT_ENDPOINT, port)
}
