package alibaba

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/ecodeclub/ekit"
	"math/rand"
	"strings"
	"time"
)

type Service struct {
	client *dysmsapi20170525.Client
}

func NewService(client *dysmsapi20170525.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) Send(ctx context.Context, tplId, signName string, numbers []string) error {
	if len(numbers) == 0 {
		return errors.New("电话号码为空")
	}
	for i := 0; i < len(numbers); i++ {
		phone := numbers[i]
		code := GenValidateCode(6)
		bcode, _ := json.Marshal(map[string]interface{}{
			"code": code,
		})
		req := &dysmsapi20170525.SendSmsRequest{
			PhoneNumbers:  ekit.ToPtr[string](phone),
			TemplateCode:  ekit.ToPtr[string](tplId),
			TemplateParam: ekit.ToPtr[string](string(bcode)),
			SignName:      ekit.ToPtr[string](signName),
		}
		runtime := &util.RuntimeOptions{}
		tryErr := func() (e error) {
			defer func() {
				if r := tea.Recover(recover()); r != nil {
					e = r
				}
			}()
			resp, err := s.client.SendSmsWithOptions(req, runtime)
			if *resp.Body.Code == "OK" {
				fmt.Println(phone, string(bcode))
			}
			if err != nil {
				fmt.Println(errors.New(*resp.Body.Message))
				return err
			}
			return nil
		}()
		if tryErr != nil {
			var e = &tea.SDKError{}
			var t *tea.SDKError
			if errors.As(tryErr, &t) {
				e = t
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
	}
	return nil
}

func GenValidateCode(width int) string {
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
