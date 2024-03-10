package logger

import (
	"bytes"
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"time"
)

type Builder struct {
	allowReqBody  bool
	allowRespBody bool
	length        int
	loggerFunc    func(ctx context.Context, al *AccessLog)
}

func NewBuilder(fn func(ctx context.Context, al *AccessLog)) *Builder {
	return &Builder{
		length:     1024,
		loggerFunc: fn,
	}
}

func (b *Builder) AllowReq() *Builder {
	b.allowReqBody = true
	return b
}

func (b *Builder) AllowResp() *Builder {
	b.allowRespBody = true
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		url := ctx.Request.URL.String()
		if len(url) > b.length {
			url = url[:b.length]
		}
		al := &AccessLog{
			Method: ctx.Request.Method,
			Url:    url,
		}
		if b.allowReqBody && ctx.Request.Body != nil {
			body, _ := ctx.GetRawData()
			ctx.Request.Body = io.NopCloser(bytes.NewReader(body))
			if len(body) > b.length {
				body = body[:b.length]
			}
			al.ReqBody = string(body)
		}
		defer func() {
			al.Duration = time.Since(start)
			if b.allowRespBody && ctx.Request.Body != nil {
				// 打日志
				b.loggerFunc(ctx, al)

			}
		}()
		// 执行业务逻辑
		ctx.Next()
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
