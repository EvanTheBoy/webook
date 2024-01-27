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
		// 如果是登录注册就不需要校验
		if ctx.Request.URL.Path == "/users/signup" ||
			ctx.Request.URL.Path == "/users/login" {
			return
		}
		sess := sessions.Default(ctx)
		// 其他的页面都需要校验是否登录
		if id := sess.Get("userId"); id == nil {
			// 未登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		updatedTime := sess.Get("update_time")
		if updatedTime == nil {
			sess.Set("update_time", 60)
			if err := sess.Save(); err != nil {
				return
			}
		}

	}
}
