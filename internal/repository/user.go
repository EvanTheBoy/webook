package repository

import (
	"context"
	"github.com/gin-gonic/gin"
	"webook/internal/domain"
	"webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

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

func (repo *UserRepository) FindByEmail(ctx context.Context, u domain.User) (domain.User, error) {
	user, err := repo.dao.SelectEmail(ctx, u)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}, nil
}

func (repo *UserRepository) UpdateUserInfo(ctx *gin.Context, u domain.User) error {
	return repo.dao.UpdateById(ctx, dao.User{
		Id:         u.Id,
		Nickname:   u.Nickname,
		Birthday:   u.Birthday,
		Address:    u.Address,
		BriefIntro: u.BriefIntro,
	})
}

func (repo *UserRepository) FindById(ctx *gin.Context, u domain.User) (domain.User, error) {
	user, err := repo.dao.SelectUserById(ctx, domain.User{
		Id: u.Id,
	})
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{
		Email:      user.Email,
		Nickname:   user.Nickname,
		Birthday:   user.Birthday,
		Address:    user.Address,
		BriefIntro: user.BriefIntro,
	}, nil
}
