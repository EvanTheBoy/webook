package repository

import (
	"context"
	"webook/internal/domain"
)

type ArticleRepository interface {
	Edit(ctx context.Context, article domain.Article)
}
