package main

import (
	"./controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong\n")
	})

	uc := controllers.UserController{}
	controllers.RegisterController(r, uc)

	g := r.Group("/v1")
	controllers.RegisterController(g, uc)

	r.Run(":8080")
}
