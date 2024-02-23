package repository

import "webook/internal/repository/cache"

type CacheRepository struct {
	cache *cache.CodeCache
}

func NewCacheRepository(c *cache.CodeCache) *CacheRepository {
	return &CacheRepository{
		cache: c,
	}
}
