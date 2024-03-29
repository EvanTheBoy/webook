package repository

import (
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"time"
	"webook/internal/domain"
	"webook/internal/repository/cache"
	"webook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFound       = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx *gin.Context, phone string) (domain.User, error)
	UpdateUserInfo(ctx *gin.Context, u domain.User) error
	FindById(ctx *gin.Context, u domain.User) (domain.User, error)
}

type GORMUserRepository struct {
	dao   dao.UserDao
	cache cache.UserCache
}

func NewUserRepository(d dao.UserDao, c cache.UserCache) UserRepository {
	return &GORMUserRepository{
		dao:   d,
		cache: c,
	}
}

func (repo *GORMUserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, repo.domainToEntity(u))
}

func (repo *GORMUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := repo.dao.SelectEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(user), nil
}

func (repo *GORMUserRepository) FindByPhone(ctx *gin.Context, phone string) (domain.User, error) {
	user, err := repo.dao.SelectPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.entityToDomain(user), nil
}

func (repo *GORMUserRepository) UpdateUserInfo(ctx *gin.Context, u domain.User) error {
	err := repo.dao.UpdateById(ctx, repo.domainToEntity(u))
	if err != nil {
		return err
	}
	err = repo.cache.Set(ctx, u)
	if err != nil {
		return err
	}
	return nil
}

func (repo *GORMUserRepository) FindById(ctx *gin.Context, u domain.User) (domain.User, error) {
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
	var ur = repo.entityToDomain(user)
	go func() {
		err = repo.cache.Set(ctx, ur)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return ur, err
}

func (repo *GORMUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:          u.Id,
		Email:       u.Email.String,
		Phone:       u.Phone.String,
		Password:    u.Password,
		Address:     u.Address,
		BriefIntro:  u.BriefIntro,
		Birthday:    u.Birthday,
		Nickname:    u.Nickname,
		CreatedTime: time.UnixMilli(u.CreatedTime),
		UpdatedTime: time.UnixMilli(u.UpdatedTime),
	}
}

func (repo *GORMUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password:    u.Password,
		Address:     u.Address,
		BriefIntro:  u.BriefIntro,
		Birthday:    u.Birthday,
		Nickname:    u.Nickname,
		CreatedTime: u.CreatedTime.UnixMilli(),
		UpdatedTime: u.UpdatedTime.UnixMilli(),
	}
}
