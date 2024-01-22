package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
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
	// 针对存储意义上的User进行存储, 在这里主要存创建时间
	// 与更新时间, 然后调用context让它能一直在链路上传递下去
	// 这里我们只需要存时间的毫秒数, 用来抵消时区的影响, 只用
	// 在前端展示的时候再做一下处理就好了
	now := time.Now().UnixMilli()
	u.CreatedTime = now
	u.UpdatedTime = now
	return dao.db.WithContext(ctx).Create(&u).Error
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string

	CreatedTime int64
	UpdatedTime int64
}
