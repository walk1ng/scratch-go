package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"gin-jwt/pkg/util"
)

var authKey = "Authorization"
var tokenKey = "Bearer"

func JwtAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get(authKey)
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 1,
				"msg":  "no access, token is required",
			})

			c.Abort()
			return
		}

		log.Println("token:", authHeader)

		s := strings.SplitN(authHeader, " ", 2)
		if !(len(s) == 2 && s[0] == tokenKey) {
			c.JSON(http.StatusOK, gin.H{
				"code": 2,
				"msg":  "invalid format auth header",
			})

			c.Abort()
			return
		}
		claims, err := util.ParseToken(s[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "invalid token",
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}
