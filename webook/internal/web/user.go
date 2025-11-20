package web

import (
	"net/http"

	regexp "github.com/dlclark/regexp2"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler() *UserHandler {
	const (
		emailRegexPattern    = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d!@#$%^&*()_\-+=]{6,8}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
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

func (u UserHandler) SignUp(c *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "请求参数错误: " + err.Error()})
		return
	}

	//邮箱校验
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4, "msg": "系统错误"})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, gin.H{"msg": "你的邮箱格式不正确"})
		return
	}

	//密码校验
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 4, "msg": "系统错误"})
		return
	}
	if !ok {
		c.JSON(http.StatusOK, gin.H{"msg": "你的密码格式不正确"})
		return
	}
	//两次密码校验
	if req.ConfirmPassword != req.Password {
		c.JSON(http.StatusOK, gin.H{"msg": "两次输入的密码不一致"})
		return
	}
	c.String(http.StatusOK, "signup success")
	//数据库操作

}
func (u UserHandler) Login(c *gin.Context)   {}
func (u UserHandler) Edit(c *gin.Context)    {}
func (u UserHandler) Profile(c *gin.Context) {}
