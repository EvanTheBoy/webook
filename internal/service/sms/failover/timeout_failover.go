package failover

import (
	"context"
	"sync/atomic"
	"webook/internal/service/sms"
)

type TimeoutFailoverSMSService struct {
	// 服务商
	services []sms.Service

	// 服务商索引
	idx int32

	// 重试次数
	cnt int32

	// 阈值(最大重试次数)
	threshold int32
}

func NewTimeoutFailoverSMSService(svcs []sms.Service, th int32) sms.Service {
	return &TimeoutFailoverSMSService{
		services:  svcs,
		threshold: th,
	}
}

func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tplId, signName, code string, numbers []string) error {
	cnt := atomic.LoadInt32(&t.cnt)
	idx := atomic.LoadInt32(&t.idx)

	// 超过阈值, 重试次数太多了, 要换服务商
	if cnt > t.threshold {
		newIdx := (idx + 1) % int32(len(t.services))

	}
}
