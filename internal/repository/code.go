package repository

import (
	"context"
	"webook/internal/repository/cache"
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCacheRepository(c *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: c,
	}
}

func (cr *CodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return cr.cache.Set(ctx, biz, code, phone)
}

func (cr *CodeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return cr.cache.Verify(ctx, biz, code, phone)
}
