package ioc

import (
	"go.uber.org/zap"
	log2 "webook/pkg/logs"
)

func InitLogger() log2.Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return log2.NewZapLogger(l)
}
