package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	Email    *regexp.Regexp
	Password *regexp.Regexp
}

func NewUserHandler() *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)

	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		Email:    emailExp,
		Password: passwordExp,
	}
}

func (u *UserHandler) RegisterUserRoutes(server *gin.Engine) {
	group := server.Group("/users")
	group.POST("/signup", u.SignUp)
	group.POST("/login", u.Login)
	group.POST("/edit", u.Edit)
	group.GET("/profile", u.Profile)
}

func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignReq

	// 浏览器根据Content-Type解析相应数据到req中, 若没有则返回错误
	if err := ctx.Bind(&req); err != nil {
		return
	}

	// 校验邮箱格式
	emailMatch, err := u.Email.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	if !emailMatch {
		ctx.String(http.StatusOK, "邮箱格式错误")
	}

	// 校验密码格式
	passwordMatch, err := u.Password.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
	}
	if !passwordMatch {
		ctx.String(http.StatusOK, "密码必须大于8位, 且包含数字、特殊字符")
	}

	// 校验两次密码是否一致
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次密码输入不一致")
	}

	// 注册成功
	ctx.String(http.StatusOK, "注册成功")
	fmt.Printf("%v", req)
}

func (u *UserHandler) Login(ctx *gin.Context) {

}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {

}
