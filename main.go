package main

import "webook/internal/integration/startup"

func main() {
	server := startup.InitWebServer()
	if err := server.Run(":8081"); err != nil {
		return
	}
}
