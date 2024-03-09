package logger

import (
	"github.com/gin-gonic/gin"
	"time"
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

	}
}

type AccessLog struct {
	// Http请求的方法
	Method string

	// 请求的Url
	Url      string
	ReqBody  string
	RespBody string
	Duration time.Duration
}
