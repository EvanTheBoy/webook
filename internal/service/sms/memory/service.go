package memory

import (
	"context"
	"fmt"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Send(ctx context.Context, tplId, signName, code string, numbers []string) error {
	fmt.Println(code)
	return nil
}
