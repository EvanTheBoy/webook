package ioc

import (
	"webook/internal/service/sms"
	"webook/internal/service/sms/failover"
	"webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	service := memory.NewService()
	return failover.NewTimeoutFailoverSMSService([]sms.Service{service}, 12)
}
