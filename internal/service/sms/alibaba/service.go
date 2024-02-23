package alibaba

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	"github.com/ecodeclub/ekit"
)

type Service struct {
	client *dysmsapi20170525.Client
}

func NewService(client *dysmsapi20170525.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) Send(ctx context.Context, tplId, signName, code string, numbers []string) error {
	if len(numbers) == 0 {
		return errors.New("电话号码为空")
	}
	for i := 0; i < len(numbers); i++ {
		phone := numbers[i]
		bcode, _ := json.Marshal(map[string]interface{}{
			"code": code,
		})
		req := &dysmsapi20170525.SendSmsRequest{
			PhoneNumbers:  ekit.ToPtr[string](phone),
			TemplateCode:  ekit.ToPtr[string](tplId),
			TemplateParam: ekit.ToPtr[string](string(bcode)),
			SignName:      ekit.ToPtr[string](signName),
		}
		_, err := s.client.SendSms(req)
		fmt.Println(phone, string(bcode))
		if err != nil {
			//fmt.Println(errors.New(*resp.Body.Message))
			return err
		}
	}
	return nil
}
