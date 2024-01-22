package dao

import (
	"context"
	"gorm.io/gorm"
)

// UserDAO 数据库, 存储意义上的用户
type UserDAO struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {

}

type User struct {
	Id       uint
	Email    string
	Password string

	CreatedTime int
	UpdatedTime int
}
