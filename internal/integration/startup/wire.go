//go:build wireinject

package startup

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
		dao.NewArticleDao,
		cache.NewUserCache,
		cache.NewCodeCache,
		repository.NewCodeRepository,
		repository.NewUserRepository,
		repository.NewArticleRepository,
		service.NewCodeService,
		service.NewUserService,
		service.NewArticleService,
		ioc.InitSMSService,
		web.NewUserHandler,
		web.NewArticleHandler,
		ioc.InitGin,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
