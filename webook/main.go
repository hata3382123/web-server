package main

import (
	"net/http"
	"strings"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"
	"webook/pkg/ginx/middleware/ratelimit"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(webook-mysql:3309)/webook"))
	if err != nil {
		//只会在初始化的时候panic
		//panic相当于整个goroutine结束
		//一旦初始化出错 应用就不要启动
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	ud := dao.NewUserDao(db)
	redisClient := redis.NewClient(&redis.Options{
		Addr: "webook-redis:11479",
	})
	userCache := cache.NewUserCache(redisClient)
	repo := repository.NewUserRepository(ud, userCache)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	server := gin.Default()

	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好 大家好")
	})
	//use 作用于全部路由
	server.Use(cors.New(cors.Config{
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
	}))
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
	// store := cookie.NewStore([]byte("secret"))
	// store := cookie.NewStore([]byte("secret"))
	// store, err := redis.NewStore(16, "tcp", "localhost:6379", "", "111111", []byte("fbVaSQV8cgR3YIxMBBoUNGoDJ3aFuCjCdDuR7iIUCxzoiSLheCqxIYdkudC9npYK"))
	// if err != nil {
	// 	panic(err)
	// }
	store := memstore.NewStore([]byte("fbVaSQV8cgR3YIxMBBoUNGoDJ3aFuCjCdDuR7iIUCxzoiSLheCqxIYdkudC9npYK"))
	server.Use(sessions.Sessions("mysession", store))
	server.Use(middleware.NewLoginJWTMiddleBuilder().Build())
	u.RegisterRoutes(server)
	server.Run(":8080")
}
