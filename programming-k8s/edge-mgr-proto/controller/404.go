package controller

import "github.com/gin-gonic/gin"

func NotFound(c *gin.Context) {
	ResponseError(c, Code404)
}
