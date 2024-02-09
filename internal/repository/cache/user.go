package cache

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"time"
	"webook/internal/domain"
)

type UserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

func (cache *UserCache) Set(ctx *gin.Context, u domain.User) error {
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return cache.client.Set(ctx, cache.key(u.Id), data, cache.expiration).Err()
}

func (cache *UserCache) Get(ctx *gin.Context, id int64) (domain.User, error) {
	data, err := cache.client.Get(ctx, cache.key(id)).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(data, &u)
	return u, err
}

func (cache *UserCache) key(id int64) string {
	return fmt.Sprintf("cache:info:%d", id)
}
