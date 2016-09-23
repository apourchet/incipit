package healthz

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const (
	DEFAULT_PORT = 10000
)

func healthCheck(ctx *gin.Context) {
	ctx.Status(200)
}

// Spawns a server that returns 200 on host:{port}/{base}/healthz
func SpawnHealthCheck(base string, port int) {
	go func() {
		engine := gin.New()
		engine.GET(base+"/healthz", healthCheck)
		engine.Run(fmt.Sprintf(":%d", port))
	}()
}
