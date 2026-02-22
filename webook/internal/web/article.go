package web

import (
	"net/http"
	"webook/internal/domain"
	"webook/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var _Handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc service.ArticleService
	l   *zap.Logger
}

func NewArticleHandler(svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
	}
}
func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")
	g.POST("/edit", h.Edit)
	g.POST("/publish", h.Publish)
	g.POST("/withdraw", h.Withdraw)
}
func (h *ArticleHandler) Edit(c *gin.Context) {
	type Editreq struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req Editreq
	if err := c.Bind(&req); err != nil {
		return
	}
	//检测输入
	//调用svc的代码
	id, err := h.svc.Save(c, domain.Article{
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.L().Error("保存文章失败", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "成功",
		Data: id,
	})
}

func (h *ArticleHandler) Publish(c *gin.Context) {
	type PublishReq struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	var req PublishReq
	if err := c.Bind(&req); err != nil {
		return
	}
	id, err := h.svc.Publish(c, domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.L().Error("发布文章失败", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "成功",
		Data: id,
	})
}
func (h *ArticleHandler) Withdraw(c *gin.Context) {
	type WithdrawReq struct {
		Id int64 `json:"id"`
	}
	var req WithdrawReq
	if err := c.Bind(&req); err != nil {
		return
	}
	err := h.svc.Withdraw(c, domain.Article{
		Id: req.Id,
	})
	if err != nil {
		c.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.L().Error("撤回文章失败", zap.Error(err))
		return
	}
	c.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "成功",
		Data: req.Id,
	})
}
