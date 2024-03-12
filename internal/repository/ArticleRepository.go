package repository

import (
	"context"
	"webook/internal/repository/dao"
)

type ArticleRepository interface {
	Edit(ctx context.Context, article dao.Article)
}

type ArticleRepositoryImpl struct {
}

func (a *ArticleRepositoryImpl) Edit(ctx context.Context, article dao.Article) {

}

func NewArticleRepository() ArticleRepository {
	return &ArticleRepositoryImpl{}
}
