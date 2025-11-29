package web

import (
	"net/http"
	"webook/internal/domain"
	"webook/internal/service"

	regexp "github.com/dlclark/regexp2"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d!@#$%^&*()_\-+=]{6,8}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")

	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
}

func (u UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "请求参数错误: " + err.Error()})
		return
	}

	//邮箱校验
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 4, "msg": "系统错误"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": "你的邮箱格式不正确"})
		return
	}

	//密码校验
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 4, "msg": "系统错误"})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": "你的密码格式不正确"})
		return
	}
	//两次密码校验
	if req.ConfirmPassword != req.Password {
		ctx.JSON(http.StatusOK, gin.H{"msg": "两次输入的密码不一致"})
		return
	}
	//调用一下svc的方法
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.String(http.StatusOK, "signup success")
	//数据库操作

}
func (u UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.JSON(http.StatusOK, gin.H{"code": 4, "msg": "用户名或密码不正确"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 4, "msg": "系统错误"})
		return
	}
	sess := sessions.Default(ctx)
	sess.Set("userId", user.Id)
	sess.Save()
	ctx.JSON(http.StatusOK, gin.H{"code": 4, "msg": "登录成功"})
	return
}
func (u UserHandler) Edit(ctx *gin.Context)    {}
func (u UserHandler) Profile(ctx *gin.Context) {}
