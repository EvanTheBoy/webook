package ratelimit

import (
	"context"
	"fmt"
	"webook/internal/service/sms"
	"webook/pkg/ratelimit"
)

type RateLimitSMSService struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewRateLimitSMSService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &RateLimitSMSService{
		svc:     svc,
		limiter: limiter,
	}
}

func (s *RateLimitSMSService) Send(ctx context.Context, tplId, signName, code string, numbers []string) error {
	limited, err := s.limiter.Limit(ctx, "alibaba:send")
	if err != nil {
		return fmt.Errorf("系统错误")
	}
	if limited {
		return fmt.Errorf("触发了限流")
	}
	return s.svc.Send(ctx, tplId, signName, code, numbers)
}
