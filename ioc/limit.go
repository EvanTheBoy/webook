package ioc

import (
	"github.com/redis/go-redis/v9"
	"time"
	ratelimit2 "webook/pkg/ratelimit"
)

func InitLimiter(redisClient redis.Cmdable) ratelimit2.Limiter {
	return ratelimit2.NewRedisSlidingWindow(redisClient, time.Second, 100)
}
