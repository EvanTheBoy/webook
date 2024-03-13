package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webook/internal/domain"
	"webook/internal/service"
	log2 "webook/pkg/logs"
)

type ArticleHandler struct {
	svc    service.ArticleService
	logger log2.Logger
}

func NewArticleHandler(s service.ArticleService, l log2.Logger) *ArticleHandler {
	return &ArticleHandler{
		svc:    s,
		logger: l,
	}
}

func (a *ArticleHandler) RegisterArticleRoutes(server *gin.Engine) {
	g := server.Group("/article")
	g.POST("/edit", a.Edit)
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Content string `json:"content"`
		Title   string `json:"title"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	id, err := a.svc.Create(ctx, domain.Article{
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.logger.Error("保存帖子失败", log2.Error(err))
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})
}
