package failover

import (
	"context"
	"errors"
	"log"
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
		// 成功则表示idx更新, 就是换了服务商了
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			atomic.StoreInt32(&t.cnt, 0)
		}
		// 因为也有可能if语句进不去, 那样的话就是并发操作导致的
		// 不管怎样，这里都要更新一下 idx
		idx = atomic.LoadInt32(&t.idx)
	}

	svc := t.services[idx]
	err := svc.Send(ctx, tplId, signName, code, numbers)
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		atomic.AddInt32(&t.cnt, 1)
	case err == nil:
		// 连续重试失败的节奏被打断了
		atomic.StoreInt32(&t.idx, 0)
	default:
		log.Default()
	}
	return err
}
