package web

import (
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
	"webook/internal/domain"
	"webook/internal/service"
)

type UserHandler struct {
	svc      *service.UserService
	Email    *regexp.Regexp
	Password *regexp.Regexp
}

func NewUserHandler(service *service.UserService) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)

	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc:      service,
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
		return
	}
	if !emailMatch {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}

	// 校验密码格式
	passwordMatch, err := u.Password.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !passwordMatch {
		ctx.String(http.StatusOK, "密码必须大于8位, 且包含数字、特殊字符")
		return
	}

	// 校验两次密码是否一致
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次密码输入不一致")
		return
	}

	// 数据库操作: handler调用下面的service
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	// 如果error在层与层之间传递的时候用的是fmt.Errorf(), 这里的判断
	// 要记得写成errors.Is(errors.Unwrap(err), 另一个参数)
	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "邮箱重复, 请换一个邮箱")
		return
	} else if err != nil {
		ctx.String(http.StatusOK, "服务器异常, 注册失败")
		return
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
