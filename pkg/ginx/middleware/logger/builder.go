package logger

type Builder struct {
	allowReq  bool
	allowResp bool
}

func (b *Builder) AllowReq() *Builder {
	b.allowReq = true
	return b
}

func (b *Builder) AllowResp() *Builder {
	b.allowResp = true
	return b
}

type AccessLog struct {
	// Http请求的方法
	Method string

	// 请求的Url
	Url      string
	ReqBody  string
	RespBody string
}
