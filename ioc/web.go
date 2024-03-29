package ioc

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
	"webook/internal/web"
	"webook/internal/web/middleware"
	"webook/pkg/ginx/middleware/logger"
	"webook/pkg/ginx/middleware/ratelimit"
	log2 "webook/pkg/logs"
	ratelimit2 "webook/pkg/ratelimit"
)

func InitGin(middlewares []gin.HandlerFunc, uHdl *web.UserHandler, aHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middlewares...)
	uHdl.RegisterUserRoutes(server)
	aHdl.RegisterArticleRoutes(server)
	return server
}

func InitMiddlewares(limiter ratelimit2.Limiter, l log2.Logger) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		// 引入日志模块
		logger.NewBuilder(func(ctx context.Context, al *logger.AccessLog) {
			l.Debug("HTTP 请求", log2.Field{Key: "al", Value: al})
		}).AllowReq().AllowResp().Build(),
		// 引入CORS的相关中间件解决跨域问题
		cors.New(cors.Config{
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			ExposeHeaders:    []string{"x-jwt-token"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://192.168.183.134") {
					return true
				}
				return strings.Contains(origin, "company.com")
			},
			MaxAge: 12 * time.Hour,
		}),
		// 设置白名单
		middleware.NewLoginMiddleWareBuilder().
			IgnorePaths("/users/signup").
			IgnorePaths("/users/login_sms/code/send").
			IgnorePaths("/users/login_sms/code/verify").
			IgnorePaths("/users/login").Build(),
		// 引入redis, 基于IP地址进行限流
		ratelimit.NewBuilder(limiter).Build(),
	}
}
