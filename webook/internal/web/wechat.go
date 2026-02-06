package web

import (
	"net/http"
	"webook/internal/service"
	"webook/internal/service/oauth2/wechat"

	"github.com/gin-gonic/gin"
)

type OAuth2WechatHandler struct {
	svc  wechat.Service
	user service.UserService
}

func NewOAuth2WechatHandler(svc wechat.Service, user service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:  svc,
		user: user,
	}
}

func (h *OAuth2WechatHandler) RegisterRoutes(s *gin.Engine) {
	g := s.Group("/oauth2/wechat")
	g.GET("/authurl", h.Auth2URL)
	g.Any("/callback", h.Callback)
}

func (h *OAuth2WechatHandler) Auth2URL(ctx *gin.Context) {
	url, err := h.svc.Auth2URL(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造扫码登录URL失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: url,
	})
}
func (h *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	return
}
