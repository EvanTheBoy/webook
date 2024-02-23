package service

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
	"webook/internal/repository"
)

type CodeService struct {
	repo *repository.CacheRepository
}

func NewCodeService(r *repository.CacheRepository) *CodeService {
	return &CodeService{
		repo: r,
	}
}

func (svc *CodeService) Send() {
	code := svc.GenValidateCode(6)

}

func (svc *CodeService) GenValidateCode(width int) string {
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
