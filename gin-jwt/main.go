package main

import (
	"gin-jwt/router"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := router.SetupRouter(gin.Default())
	log.Fatalln(r.Run(":8089"))
}
