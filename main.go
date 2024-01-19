package main

import "github.com/gin-gonic/gin"

func main() {
	server := gin.Default()

	if err := server.Run(":8080"); err != nil {
		return
	}
}
