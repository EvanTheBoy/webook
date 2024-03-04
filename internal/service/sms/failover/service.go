package failover

import (
	"context"
	"webook/internal/service/sms"
)

type FailoverSMSService struct {
	services []sms.Service
	idx      uint64
}

func NewFailoverSMSService() sms.Service {

}

func (f *FailoverSMSService) Send(ctx context.Context, tplId, signName, code string, numbers []string) error {

}
