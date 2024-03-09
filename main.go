package main

import "go.uber.org/zap"

func InitLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)
	zap.L().Info("搞好了")
}

func main() {
	InitLogger()
	server := InitWebServer()
	if err := server.Run(":8081"); err != nil {
		return
	}
}
