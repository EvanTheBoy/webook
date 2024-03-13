package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

type ArticleRepository interface {
	Create(ctx context.Context, article domain.Article) (int64, error)
}

type ArticleRepositoryImpl struct {
	dao dao.ArticleDao
}

func (a *ArticleRepositoryImpl) Create(ctx context.Context, article domain.Article) (int64, error) {
	return a.dao.Insert(ctx, dao.Article{
		Title:    article.Title,
		Content:  article.Content,
		AuthorId: article.Author.Id,
	})
}

func NewArticleRepository(d dao.ArticleDao) ArticleRepository {
	return &ArticleRepositoryImpl{
		dao: d,
	}
}
