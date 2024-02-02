package middleware

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
	"webook/internal/web"
)

type LoginMiddleWareBuilder struct {
}

func NewLoginMiddleWareBuilder() *LoginMiddleWareBuilder {
	return &LoginMiddleWareBuilder{}
}

func (l *LoginMiddleWareBuilder) Build() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		// 如果是登录注册就不需要校验
		if ctx.Request.URL.Path == "/users/signup" ||
			ctx.Request.URL.Path == "/users/login" {
			return
		}
		// 使用token进行校验
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segments := strings.Split(tokenHeader, " ")
		if len(segments) != 2 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims := &web.UserClaims{}
		tokenStr := segments[1]
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("MKdBdqsaVyzxj1WM3ZZsDeZrmv0zLDLG"), nil
		})
		if err != nil || token == nil ||
			!token.Valid || claims.Uid == 0 {
			// 未登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			// 严重的安全问题
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		now := time.Now()
		// 每十秒刷新一次, 这里生成一个新的token
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(now.Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("MKdBdqsaVyzxj1WM3ZZsDeZrmv0zLDLG"))
			if err != nil {
				log.Panicln("jwt续约失败", err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}
		// 往ctx中加入userId
		ctx.Set("claims", claims)
	}
}
