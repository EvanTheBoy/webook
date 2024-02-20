package alibaba

import (
	"context"
	"encoding/json"
	"fmt"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/ecodeclub/ekit"
	"strings"
)

type Service struct {
	signName *string
	client   *dysmsapi20170525.Client
}

func NewService(client *dysmsapi20170525.Client, signName string) *Service {
	return &Service{
		client:   client,
		signName: ekit.ToPtr[string](signName),
	}
}

func (s *Service) Send(ctx context.Context, tpl string, args string, numbers string) error {
	req := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  ekit.ToPtr[string](numbers),
		TemplateCode:  ekit.ToPtr[string](tpl),
		TemplateParam: ekit.ToPtr[string](args),
		SignName:      s.signName,
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				e = r
			}
		}()
		_, err := s.client.SendSmsWithOptions(req, runtime)
		if err != nil {
			return err
		}
		return nil
	}()
	if tryErr != nil {
		var e = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			e = _t
		} else {
			e.Message = tea.String(tryErr.Error())
		}
		// 错误 message
		fmt.Println(tea.StringValue(e.Message))
		// 诊断地址
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(e.Data)))
		err := d.Decode(&data)
		if err != nil {
			return err
		}
		if m, ok := data.(map[string]interface{}); ok {
			recommend, _ := m["Recommend"]
			fmt.Println(recommend)
		}
		_, err = util.AssertAsString(e.Message)
		if err != nil {
			return err
		}
	}
	return nil
}
