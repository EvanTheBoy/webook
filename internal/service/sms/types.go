package sms

import "context"

type Service interface {
	Send(ctx context.Context, tplId, signName, code string, numbers []string) error
}
