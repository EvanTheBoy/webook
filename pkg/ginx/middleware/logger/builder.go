package logger

import (
	"github.com/gin-gonic/gin"
)

type Builder struct {
	allowReq  bool
	allowResp bool
	length    int
}

func NewBuilder() *Builder {
	return &Builder{
		length: 1024,
	}
}

func (b *Builder) AllowReq() *Builder {
	b.allowReq = true
	return b
}

func (b *Builder) AllowResp() *Builder {
	b.allowResp = true
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		url := ctx.Request.URL.String()
		if len(url) > b.length {
			url = url[:b.length]
		}

		al := &AccessLog{
			Method: ctx.Request.Method,
			Url:    url,
		}
	}
}

type AccessLog struct {
	// Http请求的方法
	Method string

	// 请求的Url
	Url      string
	ReqBody  string
	RespBody string
}
