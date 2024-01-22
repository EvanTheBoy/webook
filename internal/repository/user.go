package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

// var ErrUserDuplicateEmail = fmt.Errorf("%w 邮箱冲突", dao.ErrUserDuplicateEmail)
var ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}
