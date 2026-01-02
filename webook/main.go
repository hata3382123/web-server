package main

func main() {
	server := InitWebServer()
	server.Run(":8888")
}

//func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
//	ud := dao.NewUserDao(db)
//	userCache := cache.NewUserCache(rdb)
//	repo := repository.NewUserRepository(ud, userCache)
//	codeCache := cache.NewCodeCache(rdb)
//	codeRepo := repository.NewCodeRepository(codeCache)
//	svc := service.NewUserService(repo)
//	smsSvc := memory.NewService()
//	codeSvc := service.NewCodeService(codeRepo, smsSvc)
//	u := web.NewUserHandler(svc, codeSvc)
//	return u
//}

//func initWebServer(rdb redis.Cmdable) *gin.Engine {
//	server := gin.Default()
//	//use 作用于全部路由
//	server.Use(cors.New(cors.Config{
//		//AllowOrigins:     []string{"https://foo.com"},
//		//AllowMethods:     []string{"PUT", "PATCH"},
//		AllowHeaders:     []string{"authorization", "content-type"},
//		ExposeHeaders:    []string{"x-jwt-token"},
//		AllowCredentials: true,
//		AllowOriginFunc: func(origin string) bool {
//			if strings.HasPrefix(origin, "http://localhost") {
//				return true
//			}
//			return strings.Contains(origin, "mycompany.com")
//		},
//		MaxAge: 12 * time.Hour,
//	}))
//	server.Use(ratelimit.NewBuilder(rdb, time.Second, 100).Build())
//	server.Use(middleware.NewLoginJWTMiddleBuilder().
//		IgnorePaths("/users/login_sms/code/send").
//		IgnorePaths("/users/login_sms").
//		Build())
//	// store := cookie.NewStore([]byte("secret"))
//	// store := cookie.NewStore([]byte("secret"))
//	// store, err := redis.NewStore(16, "tcp", "localhost:6379", "", "111111", []byte("fbVaSQV8cgR3YIxMBBoUNGoDJ3aFuCjCdDuR7iIUCxzoiSLheCqxIYdkudC9npYK"))
//	// if err != nil {
//	// 	panic(err)
//	// }
//	store := memstore.NewStore([]byte("fbVaSQV8cgR3YIxMBBoUNGoDJ3aFuCjCdDuR7iIUCxzoiSLheCqxIYdkudC9npYK"))
//	server.Use(sessions.Sessions("mysession", store))
//	return server
//}
