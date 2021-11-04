package router

import (
	"net/http"

	apiv1 "webhook/api/v1"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.POST("/mutate", apiv1.MutatingAdmission)
			v1.POST("/validate", apiv1.ValidatingAdmission)
		}
	}

	return r
}
