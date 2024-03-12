package service

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository"
)

type ArticleService interface {
	Edit(ctx context.Context, article domain.Article)
}

type ArticleServiceImpl struct {
	repo *repository.ArticleRepository
}

func NewArticleService(r *repository.ArticleRepository) ArticleService {
	return &ArticleServiceImpl{
		repo: r,
	}
}

func (a *ArticleServiceImpl) Edit(ctx context.Context, article domain.Article) {

}
