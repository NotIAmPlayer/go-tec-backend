package main

import (
	"go-tec-backend/config"
	"go-tec-backend/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	r := gin.Default()

	routes.SetupRoutes(r)

	r.Run(":8080")
}
