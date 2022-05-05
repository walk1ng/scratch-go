package controller

import (
	"gin-jwt/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	if !(username == "admin" && password == "123abc") {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "wrong username or password",
		})
		return
	}

	user := util.User{
		Username: username,
		Password: password,
	}

	token, err := util.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
		"data": gin.H{
			"token": token,
		},
	})
}

func GetSomething(c *gin.Context) {
	claims, ok := c.Get("claims")
	if !ok {
		claims = "nothing set"
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 1,
		"msg":  "success",
		"data": gin.H{
			"claims": claims,
		},
	})
}
