package main

import (
	"fmt"
	"go-tec-backend/config"
	"go-tec-backend/middlewares"
	"go-tec-backend/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// load environnment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	

	useDebugMode, exists := os.LookupEnv("DEBUG_MODE")

	if !exists {
		log.Fatal("DEBUG_MODE environment variable is not set")
		return
	}

	if useDebugMode == "false" {
		gin.SetMode(gin.ReleaseMode)
	}

	frontend, exists := os.LookupEnv("FRONTEND_URL")

	if !exists {
		log.Fatal("FRONTEND_URL environment variable is not set")
		return
	}

	// start backend
	config.ConnectDB()

	r := gin.Default()

	fmt.Println("current gin mode:" + gin.Mode())

	// Set up CORS middleware
	r.Use(middlewares.HandleCORS(frontend))

	routes.SetupRoutes(r)

	r.Run(":8080")
}
