package middleware

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddleBuilder struct {
}

func NewLoginMiddleBuilder() *LoginMiddleBuilder {
	return &LoginMiddleBuilder{}
}

func (l *LoginMiddleBuilder) Build() gin.HandlerFunc {
	//用go的方式编码解码
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/users/signup" ||
			ctx.Request.URL.Path == "/users/login" {
			return
		}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			//没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		// 与登录时保持一致的选项，防止覆盖后 Cookie 无效
		sess.Options(sessions.Options{
			MaxAge:   30 * 60,
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
		now := time.Now()
		//刚刚登录还没刷新
		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
		//有updateTime的情况
		updateTimeVal, _ := updateTime.(time.Time)
		if now.Sub(updateTimeVal) > time.Second*10 {
			sess.Set("update_time", now)
			sess.Save()
		}
	}
}
