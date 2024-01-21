package repository

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDAO
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	err := repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		return err
	}
	return nil
}
