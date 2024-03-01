//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webook/internal/repository"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, ioc.InitRedis,
		dao.NewUserDao,
		cache.NewUserCache,
		cache.NewCodeCache,
		repository.NewCodeRepository,
		repository.NewUserRepository,
		service.NewCodeService,
		service.NewUserService,
		ioc.InitSMSService,
		web.NewUserHandler,
		ioc.InitGin,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
	//cmdable := ioc.InitRedis()
	//v := ioc.InitMiddlewares(cmdable)
	//db := ioc.InitDB()
	//userDao := dao.NewUserDao(db)
	//userCache := cache.NewUserCache(cmdable)
	//userRepository := repository.NewUserRepository(userDao, userCache)
	//userService := service.NewUserService(userRepository)
	//codeCache := cache.NewCodeCache(cmdable)
	//codeRepository := repository.NewCodeRepository(codeCache)
	//smsService := memory.NewService()
	//codeService := service.NewCodeService(codeRepository, smsService)
	//userHandler := web.NewUserHandler(userService, codeService)
	//engine := ioc.InitGin(v, userHandler)
	//return engine
}
