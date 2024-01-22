package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
)

func main() {
	server := gin.Default()

	// 初始化数据库操作需要的组件
	db, err := gorm.Open(mysql.Open("root:root@tcp(192.168.183.132:13316)/webook"))
	if err != nil {
		// 结束goroutine
		// 一旦初始化过程中出错, 应用就不要启动
		panic(err)
	}
	userDao := dao.NewUserDao(db)
	repo := repository.NewUserRepository(userDao)
	svc := service.NewUserService(repo)

	// 引入CORS的相关中间件解决跨域问题
	server.Use(cors.New(cors.Config{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://192.168.183.132") {
				return true
			}
			return strings.Contains(origin, "company.com")
		},
		MaxAge: 12 * time.Hour,
	}))

	u := web.NewUserHandler(svc)
	u.RegisterUserRoutes(server)
	if err := server.Run(":8081"); err != nil {
		return
	}
}
