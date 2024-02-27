//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"
	"webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB, initRedis,
		dao.NewUserDao,
		repository.NewUserRepository,
		service.NewUserService,
		web.NewUserHandler,
	)
	return new(gin.Engine)
}
