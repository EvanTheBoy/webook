package web

import (
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
	"webook/internal/domain"
	"webook/internal/service"
)

const biz = "login"

type UserHandler struct {
	svc        service.UserService
	codeSvc    *service.CodeService
	Email      *regexp.Regexp
	Password   *regexp.Regexp
	Birthday   *regexp.Regexp
	Nickname   *regexp.Regexp
	Address    *regexp.Regexp
	BriefIntro *regexp.Regexp
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}

func NewUserHandler(service service.UserService, codeSvc *service.CodeService) *UserHandler {
	const (
		emailRegexPattern      = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern   = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
		birthdayRegexPattern   = `^\d{4}-\d{2}-\d{2}$`
		nicknameRegexPattern   = `^[\u4e00-\u9fa5a-zA-Z0-9]{4,20}$`
		addressRegexPattern    = `^[\u4e00-\u9fa5a-zA-Z0-9]{0,40}$`
		briefIntroRegexPattern = `^^.{0,60}$`
	)

	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	birthdayExp := regexp.MustCompile(birthdayRegexPattern, regexp.None)
	nicknameExp := regexp.MustCompile(nicknameRegexPattern, regexp.None)
	addressExp := regexp.MustCompile(addressRegexPattern, regexp.None)
	briefIntroExp := regexp.MustCompile(briefIntroRegexPattern, regexp.None)
	return &UserHandler{
		svc:        service,
		codeSvc:    codeSvc,
		Email:      emailExp,
		Password:   passwordExp,
		Birthday:   birthdayExp,
		Nickname:   nicknameExp,
		Address:    addressExp,
		BriefIntro: briefIntroExp,
	}
}

func (u *UserHandler) RegisterUserRoutes(server *gin.Engine) {
	group := server.Group("/users")
	group.POST("/signup", u.SignUp)
	group.POST("/login", u.Login)
	group.POST("/edit", u.Edit)
	group.GET("/profile", u.Profile)
	group.POST("/login_sms/code/send", u.SendLoginSmsCode)
	group.POST("/login_sms/code/verify", u.VerifyLoginSmsCode)
}

func (u *UserHandler) SendLoginSmsCode(ctx *gin.Context) {
	type SendCodeReq struct {
		Phone string `json:"phone"`
	}
	var req SendCodeReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if err := u.codeSvc.Send(ctx, biz, req.Phone); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "发送成功")
}

func (u *UserHandler) VerifyLoginSmsCode(ctx *gin.Context) {
	type VerifyCodeReq struct {
		Code  string `json:"code"`
		Phone string `json:"phone"`
	}
	var req VerifyCodeReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := u.codeSvc.Verify(ctx, biz, req.Code, req.Phone)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码错误",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	if err = u.setJWTToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "验证成功",
	})
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
	type LoginReq struct {
		Email    string
		Password string
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "用户名或密码错误")
		return
	} else if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	if err = u.setJWTToken(ctx, user.Id); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	// 登录成功
	ctx.String(http.StatusOK, "登录成功")
	fmt.Printf("%v", req)
}

func (u *UserHandler) setJWTToken(ctx *gin.Context, uid int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置过期时间
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
		},
		UserAgent: ctx.Request.UserAgent(),
		Uid:       uid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("MKdBdqsaVyzxj1WM3ZZsDeZrmv0zLDLG"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	type UserInfo struct {
		Nickname   string
		Birthday   string
		Address    string
		BriefIntro string
	}
	var req UserInfo
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 校验昵称
	nicknameMatch, err := u.Nickname.MatchString(req.Nickname)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !nicknameMatch {
		ctx.String(http.StatusOK, "昵称长度必须大于4位且小于20位")
		return
	}
	// 校验生日
	birthdayMatch, err := u.Birthday.MatchString(req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !birthdayMatch {
		ctx.String(http.StatusOK, "生日应符合格式: YYYY-MM-DD")
		return
	}
	// 校验地区
	matchAddress, err := u.Address.MatchString(req.Address)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !matchAddress {
		ctx.String(http.StatusOK, "地区文本过长, 应在40以内")
		return
	}
	// 校验个人简介
	matchBriefIntro, err := u.BriefIntro.MatchString(req.BriefIntro)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !matchBriefIntro {
		ctx.String(http.StatusOK, "简介文本过长, 应在60以内")
		return
	}
	// 从jwt中获取用户id
	c, _ := ctx.Get("claims")
	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	err = u.svc.UpdateUserInfo(ctx, domain.User{
		Id:         claims.Uid,
		Nickname:   req.Nickname,
		Birthday:   req.Birthday,
		Address:    req.Address,
		BriefIntro: req.BriefIntro,
	})
	if errors.Is(err, service.ErrUserNotFound) {
		ctx.String(http.StatusOK, "用户不存在")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 更新成功
	ctx.String(http.StatusOK, "更新成功")
	fmt.Printf("%v", req)
}

func (u *UserHandler) Profile(ctx *gin.Context) {
	c, _ := ctx.Get("claims")
	claims, ok := c.(*UserClaims)
	if !ok {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	user, err := u.svc.SearchById(ctx, domain.User{
		Id: claims.Uid,
	})
	if errors.Is(err, service.ErrUserNotFound) {
		ctx.String(http.StatusOK, "用户不存在")
		return
	} else if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	// 查询成功
	ctx.JSON(http.StatusOK, gin.H{
		"Email":      user.Email,
		"Nickname":   user.Nickname,
		"Birthday":   user.Birthday,
		"Address":    user.Address,
		"BriefIntro": user.BriefIntro,
	})
}
