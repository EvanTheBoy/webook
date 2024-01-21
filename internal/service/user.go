package service

import (
	"context"
	"webook/internal/domain"
	"webook/internal/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	// 存储
	err := svc.repo.Create(ctx, u)
	if err != nil {
		return err
	}
	return nil
}
