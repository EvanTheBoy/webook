package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrSendCodeTooManyTimes = errors.New("验证码发送太频繁")
	ErrVerifyTooManyTimes   = errors.New("验证次数太多")
	ErrSystemAnomaly        = errors.New("系统错误")
	ErrCodeNotCorrect       = errors.New("验证码错误")
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache interface {
	Set(ctx context.Context, biz, code, phone string) error
	Verify(ctx context.Context, biz, code, phone string) error
}

type CodeCacheImpl struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeCache {
	return &CodeCacheImpl{
		client: client,
	}
}

func (c *CodeCacheImpl) Set(ctx context.Context, biz, code, phone string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return ErrSendCodeTooManyTimes
	default:
		return ErrSystemAnomaly
	}
}

func (c *CodeCacheImpl) Verify(ctx context.Context, biz, code, phone string) error {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return ErrVerifyTooManyTimes
	case -2:
		return ErrCodeNotCorrect
	default:
		return ErrSystemAnomaly
	}
}

func (c *CodeCacheImpl) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
