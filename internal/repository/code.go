package repository

import (
	"context"
	"webook/internal/repository/cache"
)

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

type CodeRepositoryImpl struct {
	cache cache.CodeCache
}

func NewCodeRepository(c cache.CodeCache) CodeRepository {
	return &CodeRepositoryImpl{
		cache: c,
	}
}

func (cr *CodeRepositoryImpl) Store(ctx context.Context, biz, phone, code string) error {
	return cr.cache.Set(ctx, biz, code, phone)
}

func (cr *CodeRepositoryImpl) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return cr.cache.Verify(ctx, biz, code, phone)
}
