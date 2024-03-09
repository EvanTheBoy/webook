package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"webook/internal/domain"
	"webook/internal/repository"
	log2 "webook/pkg/logs"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrUserNotFound          = repository.ErrUserNotFound
	ErrInvalidUserOrPassword = errors.New("账号或密码错误")
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, u domain.User) (domain.User, error)
	UpdateUserInfo(ctx *gin.Context, u domain.User) error
	SearchById(ctx *gin.Context, u domain.User) (domain.User, error)
	FindOrCreate(ctx *gin.Context, phone string) (domain.User, error)
}

type UserServiceImpl struct {
	repo   repository.UserRepository
	logger log2.Logger
}

func NewUserService(repo repository.UserRepository, l log2.Logger) UserService {
	return &UserServiceImpl{
		repo:   repo,
		logger: l,
	}
}

func (svc *UserServiceImpl) SignUp(ctx context.Context, u domain.User) error {
	// 对密码进行加密后存储
	password, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(password)
	// 考虑数据库存储的操作
	return svc.repo.Create(ctx, u)
}

func (svc *UserServiceImpl) Login(ctx context.Context, u domain.User) (domain.User, error) {
	// 先查找用户
	user, err := svc.repo.FindByEmail(ctx, u.Email)
	if errors.Is(err, repository.ErrUserNotFound) {
		svc.logger.Error("账号或密码错误")
		// 笼统化, 不能告诉用户具体是账号有问题还是密码有问题
		return domain.User{}, ErrInvalidUserOrPassword
	} else if err != nil {
		return domain.User{}, err
	}
	// 比较密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if err != nil {
		// 笼统化, 不能告诉用户具体是账号有问题还是密码有问题
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return user, nil
}

func (svc *UserServiceImpl) UpdateUserInfo(ctx *gin.Context, u domain.User) error {
	return svc.repo.UpdateUserInfo(ctx, u)
}

func (svc *UserServiceImpl) SearchById(ctx *gin.Context, u domain.User) (domain.User, error) {
	user, err := svc.repo.FindById(ctx, u)
	return user, err
}

func (svc *UserServiceImpl) FindOrCreate(ctx *gin.Context, phone string) (domain.User, error) {
	user, err := svc.repo.FindByPhone(ctx, phone)
	// 若存在用户, 把用户和错误一并返回
	if !errors.Is(err, ErrUserNotFound) {
		svc.logger.Error("用户不存在")
		return user, err
	}
	// 若不存在用户, 当场创建
	user = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, user)
	if err != nil {
		return user, err
	}
	// 会遇到主从延迟的问题
	return svc.repo.FindByPhone(ctx, phone)
}
