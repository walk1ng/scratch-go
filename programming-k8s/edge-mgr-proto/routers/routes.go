package routers

import (
	"edge-mgr-proto/controller"
	"edge-mgr-proto/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(mode string) *gin.Engine {
	// release mode
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "PONG!")
	})

	// /api/v1
	v1Group := r.Group("/api/v1")
	initK8sRouter(v1Group)

	r.NoRoute(controller.NotFound)

	return r
}
