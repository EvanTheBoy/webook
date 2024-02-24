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
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) *CodeCache {
	return &CodeCache{
		client: client,
	}
}

func (c *CodeCache) Set(ctx context.Context, biz, code, phone string) error {
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

func (c *CodeCache) Verify(ctx context.Context, biz, code, phone string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		return false, ErrVerifyTooManyTimes
	default:
		return false, ErrSystemAnomaly
	}
}

func (c *CodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
