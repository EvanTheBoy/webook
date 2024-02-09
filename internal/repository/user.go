package repository

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(d *dao.UserDAO, c *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   d,
		cache: c,
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
	// 先从缓存中查找
	uc, err := repo.cache.Get(ctx, u.Id)
	if err == nil {
		return uc, nil
	}
	// 缓存中没有就从数据库查
	user, err := repo.dao.SelectUserById(ctx, domain.User{
		Id: u.Id,
	})
	if err != nil {
		return domain.User{}, err
	}
	var ur = domain.User{
		Email:      user.Email,
		Nickname:   user.Nickname,
		Birthday:   user.Birthday,
		Address:    user.Address,
		BriefIntro: user.BriefIntro,
	}
	go func() {
		err = repo.cache.Set(ctx, ur)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return ur, err
}
