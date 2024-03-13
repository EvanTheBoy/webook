package dao

import (
	"context"
	"gorm.io/gorm"
)

type ArticleDao interface {
	Insert(ctx context.Context, d Article) (int64, error)
}

type ArticleDaoImpl struct {
	db *gorm.DB
}

func NewArticleDao(db *gorm.DB) ArticleDao {
	return &ArticleDaoImpl{
		db: db,
	}
}

func (a *ArticleDaoImpl) Insert(ctx context.Context, d Article) (int64, error) {

}

type Article struct {
	Id          int64  `gorm:"primaryKey;autoIncrement"`
	Title       string `gorm:"type=varchar(1024)"`
	Content     string `gorm:"BLOB"`
	AuthorId    int64  `gorm:"index"`
	CreatedTime int64
	UpdatedTime int64
}
