package service

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository"
)

type ArticleService interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
}

type ArticleServiceImpl struct {
	repo repository.ArticleRepository
}

func NewArticleService(r repository.ArticleRepository) ArticleService {
	return &ArticleServiceImpl{
		repo: r,
	}
}

func (a *ArticleServiceImpl) Create(ctx context.Context, article domain.Article) (int64, error) {
	if article.Id > 0 {
		err := a.repo.Update(ctx, article)
		return article.Id, err
	}
	return a.repo.Create(ctx, article)
}
