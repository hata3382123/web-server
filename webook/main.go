package main

import (
	"strings"
	"time"
	"webook/internal/web"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	//use 作用于全部路由
	server.Use(cors.New(cors.Config{
		//AllowOrigins:     []string{"https://foo.com"},
		//AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders: []string{"authorization", "content-type"},
		//ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "mycompany.com")
		},
		MaxAge: 12 * time.Hour,
	}))
	u := web.NewUserHandler()
	u.RegisterRoutes(server)
	server.Run(":8080")
}
