package main

import (
	"go-tec-backend/config"
	"go-tec-backend/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	// load environnment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	frontend, exists := os.LookupEnv("FRONTEND_URL")

	if !exists {
		log.Fatal("FRONTEND_URL environment variable is not set")
		return
	}

	// start backend
	config.ConnectDB()

	r := gin.Default()

	// Set up CORS middleware
	r.Use(func(c *gin.Context) {
		if gin.Mode() == gin.ReleaseMode {
			c.Writer.Header().Set("Access-Control-Allow-Origin", frontend)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins in development mode
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // No Content
			return
		}
		c.Next()
	})

	routes.SetupRoutes(r)

	r.Run(":8080")
}
