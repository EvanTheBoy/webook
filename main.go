package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
	"webook/config"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/internal/web/middleware"
	"webook/pkg/ginx/middleware/ratelimit"
)

func main() {
	db := initDB()
	rdb := initRedis()
	user := initUser(db, rdb)
	server := initWebServer(rdb)
	user.RegisterUserRoutes(server)
	if err := server.Run(":8081"); err != nil {
		return
	}
}

func initWebServer(redisClient redis.Cmdable) *gin.Engine {
	server := gin.Default()

	// 引入redis, 基于IP地址进行限流
	server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())

	// 引入CORS的相关中间件解决跨域问题
	server.Use(cors.New(cors.Config{
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
	}))

	server.Use(middleware.NewLoginMiddleWareBuilder().Build())
	return server
}

func initDB() *gorm.DB {
	// 初始化数据库操作需要的组件
	//db, err := gorm.Open(mysql.Open("root:root@tcp(192.168.183.134:13316)/webook"))
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		// 结束goroutine
		// 一旦初始化过程中出错, 应用就不要启动
		panic(err)
	}
	// 建表
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return redisClient
}

func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
	userDao := dao.NewUserDao(db)
	userCache := cache.NewUserCache(rdb)
	repo := repository.NewUserRepository(userDao, userCache)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}
