package web

import (
	"fmt"
	"net/http"
	"time"
	"webook/internal/domain"
	"webook/internal/service"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

const biz = "login"

type UserHandler struct {
	svc         service.UserService
	codeSvc     service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc service.UserService, codeSvc service.CodeService) *UserHandler {
	const (
		emailRegexPattern    = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d!@#$%^&*()_\-+=]{6,8}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		codeSvc:     codeSvc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")

	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/edit", u.Edit)
	ug.GET("/profile", u.ProfileJWT)
	ug.POST("/login_sms/code/send", u.sendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSMS)
}
func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 1,
			Msg:  "请求参数错误: " + err.Error(),
		})
		return
	}

	ok, err := u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		fmt.Printf("Verify error: %v\n", err)
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误: " + err.Error(),
		})
		return
	}
	if !ok {

		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "验证码输入错误",
		})
		return
	}

	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		fmt.Printf("FindOrCreate error: %v\n", err)
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误: " + err.Error(),
		})
		return
	}

	if err = u.setJWTToken(ctx, user); err != nil {
		fmt.Printf("setJWTToken error: %v\n", err)
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误: " + err.Error(),
		})
		return
	}
	fmt.Printf("LoginSMS success, returning code: 0\n")
	ctx.JSON(http.StatusOK, Result{
		Code: 0,
		Msg:  "登录成功",
	})
}
func (u *UserHandler) sendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone" binding:"required"`
	}

	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Result{
			Code: 1,
			Msg:  "请求参数错误: " + err.Error(),
		})
		return
	}
	err = u.codeSvc.Send(ctx, biz, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "发送成功",
	})
}
func (u *UserHandler) SignUp(ctx *gin.Context) {
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
	if err == service.ErrUserDuplicate {
		ctx.JSON(http.StatusOK, gin.H{"msg": "邮箱冲突"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"msg": "系统异常"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "注册成功"})
	//数据库操作

}
func (u *UserHandler) Login(ctx *gin.Context) {
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
	// 本地开发使用 http，禁止 Secure 避免浏览器不带 Cookie；上线 https 再打开
	sess.Options(sessions.Options{
		Secure:   false,
		HttpOnly: true,
		MaxAge:   30 * 60,
		// SameSiteLax 既防 CSRF，又能在同站子域场景工作
		SameSite: http.SameSiteLaxMode,
	})
	sess.Save()
	ctx.JSON(http.StatusOK, gin.H{"code": 4, "msg": "登录成功"})
}
func (u *UserHandler) LoginJWT(ctx *gin.Context) {
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
	if err = u.setJWTToken(ctx, user); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 5, "msg": "系统错误"})
		return
	}

	fmt.Println(user)
	ctx.JSON(http.StatusOK, gin.H{"code": 4, "msg": "登录成功"})
}

func (u *UserHandler) setJWTToken(ctx *gin.Context, user domain.User) error {
	claims := UserClaims{
		Uid:       user.Id,
		UserAgent: ctx.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			// 根据需要设置过期时间、颁发者等
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenstr, err := token.SignedString([]byte("fbVaSQV8cgR3YIxMBBoUNGoDJ3aFuCjCdDuR7iIUCxzoiSLheCqxIYdkudC9npYK"))
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenstr)
	return nil
}
func (u *UserHandler) LoginOut(ctx *gin.Context) {
	sess := sessions.Default(ctx)
	sess.Delete("userId")
	sess.Options(sessions.Options{
		MaxAge: -1,
	})
	sess.Save()
	ctx.JSON(http.StatusOK, gin.H{"code": 4, "msg": "退出登录成功"})
}
func (u *UserHandler) Edit(ctx *gin.Context) {
	type Editreq struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req Editreq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"code": 1, "msg": "请求参数错误: " + err.Error()})
		return
	}
	// 中间件已经验证了登录，可以直接从 context 获取 userId
	userId := ctx.MustGet("userId").(int64)
	err := u.svc.Edit(ctx, domain.User{
		Id:       userId,
		Nickname: req.Nickname,
		Birthday: req.Birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 4, "msg": "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 4, "msg": "修改成功"})
}
func (u *UserHandler) Profile(ctx *gin.Context) {
	// 中间件已验证 JWT，将 uid 放入上下文
	uidVal, ok := ctx.Get("userId")
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userId, ok := uidVal.(int64)
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	user, err := u.svc.Profile(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 4, "msg": "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 4, "msg": "这个profile", "data": user})
}
func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, _ := ctx.Get("claims")
	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	user, err := u.svc.Profile(ctx, claims.Uid)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 4, "msg": "系统错误"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 4, "msg": "这个profile", "data": user})
}

type UserClaims struct {
	jwt.RegisteredClaims
	//声明自己要放进token里的数据
	Uid       int64
	UserAgent string
}
