package ratelimit

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var slidingWindowLuaScript string

type RedisSlidingWindow struct {
	cmd      redis.Cmdable
	interval time.Duration
	// 阈值
	rate int
}

func NewRedisSlidingWindow(cmd redis.Cmdable, interval time.Duration, rate int) Limiter {
	return &RedisSlidingWindow{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}

func (r *RedisSlidingWindow) Limit(ctx context.Context, key string) (bool, error) {
	return r.cmd.Eval(ctx, slidingWindowLuaScript, []string{key},
		r.interval.Milliseconds(), r.rate, time.Now().UnixMilli()).Bool()
}
