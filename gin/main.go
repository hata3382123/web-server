package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	server.GET("/hello/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "hello go %s", name)
	})
	server.GET("/order", func(c *gin.Context) {
		oid := c.Query("id")
		c.String(http.StatusOK, "hello go %s", oid)
	})
	server.POST("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "hello post")
	})
	server.Run(":8080")
}
