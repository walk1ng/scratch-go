package router

import (
	"gin-jwt/controller"
	"gin-jwt/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) *gin.Engine {
	r.GET("/login", controller.Login)
	apiv1 := r.Group("/api/v1")
	apiv1.Use(middleware.JwtAuth())
	{
		apiv1.GET("/foo", controller.GetSomething)
	}
	return r
}
