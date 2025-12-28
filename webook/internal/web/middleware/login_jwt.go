package middleware

import (
	"log"
	"net/http"
	"strings"
	"time"
	"webook/internal/web"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginJWTMiddleBuilder struct {
}

func NewLoginJWTMiddleBuilder() *LoginJWTMiddleBuilder {
	return &LoginJWTMiddleBuilder{}
}

func (l *LoginJWTMiddleBuilder) Build() gin.HandlerFunc {
	//用go的方式编码解码
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/users/signup" ||
			ctx.Request.URL.Path == "/users/login" {
			return
		}
		// 用 JWT 校验：优先 x-jwt-token，兼容 Authorization: Bearer <token>
		tokenHeader := strings.TrimSpace(ctx.GetHeader("x-jwt-token"))
		if tokenHeader == "" {
			tokenHeader = strings.TrimSpace(ctx.GetHeader("Authorization"))
		}
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if strings.HasPrefix(strings.ToLower(tokenHeader), "bearer ") {
			tokenHeader = strings.TrimSpace(tokenHeader[7:])
		}
		tokenStr := tokenHeader
		claims := &web.UserClaims{}
		// ParseWithClaims 里面一定要传指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("fbVaSQV8cgR3YIxMBBoUNGoDJ3aFuCjCdDuR7iIUCxzoiSLheCqxIYdkudC9npYK"), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			//严重的安全问题
			//加监控
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//设置刷新时间 每10秒刷新一次
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("fbVaSQV8cgR3YIxMBBoUNGoDJ3aFuCjCdDuR7iIUCxzoiSLheCqxIYdkudC9npYK"))
			if err != nil {
				log.Panicln("JWT 续约失败", err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}

		// 将 uid 写入 context，后续 handler 从 context 读取
		ctx.Set("userId", claims.Uid)
		ctx.Set("claims", claims)
	}
}
