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
		ioc.InitLogger,
		ioc.InitDB, ioc.InitRedis, ioc.InitLimiter,
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
}
