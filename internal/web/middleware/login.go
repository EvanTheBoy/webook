package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginMiddleWareBuilder struct {
}

func NewLoginMiddleWareBuilder() *LoginMiddleWareBuilder {
	return &LoginMiddleWareBuilder{}
}

func (l *LoginMiddleWareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if ctx.Request.URL.Path == "/users/signup" ||
			ctx.Request.URL.Path == "/users/login" {
			return
		}
		sess := sessions.Default(ctx)
		if id := sess.Get("userId"); id == nil {
			// 未登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}
