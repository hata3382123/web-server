package main

import (
	"strings"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/service/sms/memory"
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
	db, err := gorm.Open(mysql.Open("root:123456@tcp(localhost:13316)/webook"))
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
		Addr: "127.0.0.1:6379",
	})
	userCache := cache.NewUserCache(redisClient)
	repo := repository.NewUserRepository(ud, userCache)
	codeCache := cache.NewCodeCache(redisClient)
	codeRepo := repository.NewCodeRepository(codeCache)
	svc := service.NewUserService(repo)
	smsSvc := memory.NewService()
	codeSvc := service.NewCodeService(codeRepo, smsSvc)
	u := web.NewUserHandler(svc, codeSvc)
	server := gin.Default()

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
	server.Use(middleware.NewLoginJWTMiddleBuilder().
		IgnorePaths("/users/login_sms/code/send").
		IgnorePaths("/users/login_sms").
		Build())
	// store := cookie.NewStore([]byte("secret"))
	// store := cookie.NewStore([]byte("secret"))
	// store, err := redis.NewStore(16, "tcp", "localhost:6379", "", "111111", []byte("fbVaSQV8cgR3YIxMBBoUNGoDJ3aFuCjCdDuR7iIUCxzoiSLheCqxIYdkudC9npYK"))
	// if err != nil {
	// 	panic(err)
	// }
	store := memstore.NewStore([]byte("fbVaSQV8cgR3YIxMBBoUNGoDJ3aFuCjCdDuR7iIUCxzoiSLheCqxIYdkudC9npYK"))
	server.Use(sessions.Sessions("mysession", store))
	u.RegisterRoutes(server)
	server.Run(":8080")
}
