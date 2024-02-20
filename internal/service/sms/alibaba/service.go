package alibaba

import (
	"context"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/ecodeclub/ekit"
)

type Service struct {
	signName *string
	client   *dysmsapi20170525.Client
}

func NewService(client *dysmsapi20170525.Client, signName string) *Service {
	return &Service{
		signName: ekit.ToPtr[string](signName),
	}
}

func (s *Service) Send(ctx context.Context, tpl string, args []string, numbers ...string) {
	req := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers: tea.String("your_value"),
		SignName:     s.signName,
	}
}
