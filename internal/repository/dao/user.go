package dao

import "context"

// UserDAO 数据库, 存储意义上的用户
type UserDAO struct {
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
