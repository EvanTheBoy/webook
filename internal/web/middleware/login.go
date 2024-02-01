package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
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
		tokenStr := segments[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte("MKdBdqsaVyzxj1WM3ZZsDeZrmv0zLDLG"), nil
		})
		fmt.Println("token=", token)
		if err != nil || token == nil || !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
