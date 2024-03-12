package web

import "github.com/gin-gonic/gin"

type ArticleHandler struct {
}

func (a *ArticleHandler) RegisterRoutes(server *gin.Engine) {
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
}
