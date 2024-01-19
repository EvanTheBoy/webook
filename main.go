package main

import (
	"github.com/gin-gonic/gin"
	"webook/internal/web"
)

func main() {
	server := gin.Default()
	u := web.NewUserHandler()
	u.RegisterUserRoutes(server)
	if err := server.Run(":8080"); err != nil {
		return
	}
}
