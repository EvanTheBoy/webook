package repository

import (
	"context"
	"webook/internal/repository/cache"
)

type CacheRepository struct {
	cache *cache.CodeCache
}

func NewCacheRepository(c *cache.CodeCache) *CacheRepository {
	return &CacheRepository{
		cache: c,
	}
}

func (cr *CacheRepository) Store(ctx context.Context, biz, phone, code string) error {
	return cr.cache.Set(ctx, biz, code, phone)
}
