package main

import (
	"github.com/redis/go-redis/v9"
	"webook/config"
)

func initRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return redisClient
}
