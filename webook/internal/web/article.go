package web

import (
	"errors"
	"net/http"
	"webook/internal/domain"
	"webook/internal/repository"
	"webook/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var _Handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc service.ArticleService
}

func NewArticleHandler(svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
	}
}
func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", h.Edit)
}
func (h *ArticleHandler) Edit(c *gin.Context) {
	// 从 JWT 中间件写入的 context 获取当前用户 ID
	uidVal, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusOK, Result{Code: 5, Msg: "未登录"})
		return
	}
	userId, ok := uidVal.(int64)
	if !ok {
		c.JSON(http.StatusOK, Result{Code: 5, Msg: "用户ID类型错误"})
		return
	}

	type Editreq struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Editreq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusOK, Result{Code: 1, Msg: "请求参数错误"})
		return
	}

	id, err := h.svc.Save(c, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author:  domain.Author{Id: userId}, // 使用 JWT 中的 userId，不信任前端
	})
	if err != nil {
		if errors.Is(err, repository.ErrArticleNotFound) {
			c.JSON(http.StatusOK, Result{Code: 4, Msg: "文章不存在或无权操作"})
			return
		}
		c.JSON(http.StatusOK, Result{Code: 5, Msg: "系统错误"})
		zap.L().Error("保存文章失败", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, Result{Code: 0, Msg: "成功", Data: id})
}
