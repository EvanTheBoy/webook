package failover

import (
	"context"
	"errors"
	"log"
	"sync/atomic"
	"webook/internal/service/sms"
)

type FailoverSMSService struct {
	services []sms.Service
	idx      uint64
}

func NewFailoverSMSService(services []sms.Service, idx uint64) sms.Service {
	return &FailoverSMSService{
		services: services,
		idx:      idx,
	}
}

func (f *FailoverSMSService) Send(ctx context.Context, tplId, signName, code string, numbers []string) error {
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.services))
	for i := idx; i < idx+length; i++ {
		svc := f.services[int(i%length)]
		err := svc.Send(ctx, tplId, signName, code, numbers)
		switch {
		case err == nil:
			return nil
		case errors.Is(err, context.DeadlineExceeded), errors.Is(err, context.Canceled):
			return err
		default:
			log.Default()
		}
	}
	return errors.New("全部服务商都失败了")
}
