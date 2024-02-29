package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
	"webook/internal/web"
	"webook/internal/web/middleware"
	"webook/pkg/ginx/middleware/ratelimit"
)

func InitGin(middlewares []gin.HandlerFunc, handler *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(middlewares...)
	handler.RegisterUserRoutes(server)
	return server
}

func InitMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
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
			IgnorePaths("/login_sms/code/send").
			IgnorePaths("/login_sms/code/verify").
			IgnorePaths("/users/login").Build(),
		// 引入redis, 基于IP地址进行限流
		ratelimit.NewBuilder(redisClient, time.Second, 100).Build(),
	}
}
