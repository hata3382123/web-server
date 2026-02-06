package ioc

import (
	"strings"
	"time"
	"webook/internal/web"
	"webook/internal/web/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler, OAuth2WechatHandler *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	OAuth2WechatHandler.RegisterRoutes(server)
	return server
}

func InitMiddlewares(rdb redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{corsHdl(),
		middleware.NewLoginJWTMiddleBuilder().
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms").Build(),
		//ratelimit.NewBuilder().Build(),
	}
}
func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins:     []string{"https://foo.com"},
		//AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"authorization", "content-type"},
		ExposeHeaders:    []string{"x-jwt-token"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "mycompany.com")
		},
		MaxAge: 12 * time.Hour,
	})
}
