package service

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"webook/internal/repository"
	"webook/internal/service/sms"
)

var (
	tplId    = "SMS_464325417"
	signName = "hmdpLogin"
)

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, code, phone string) (bool, error)
}

type CodeServiceImpl struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
}

func NewCodeService(r repository.CodeRepository, svc sms.Service) CodeService {
	return &CodeServiceImpl{
		repo:   r,
		smsSvc: svc,
	}
}

func (svc *CodeServiceImpl) Send(ctx context.Context, biz, phone string) error {
	code := svc.genValidateCode(6)
	// 先存入 redis
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	// 再发送
	err = svc.smsSvc.Send(ctx, tplId, signName, code, []string{phone})
	if err != nil {
		return err
	}
	return nil
}

func (svc *CodeServiceImpl) Verify(ctx context.Context, biz, code, phone string) (bool, error) {
	return svc.repo.Verify(ctx, biz, code, phone)
}

func (svc *CodeServiceImpl) genValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.NewSource(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < width; i++ {
		_, err := fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
		if err != nil {
			return ""
		}
	}
	return sb.String()
}
