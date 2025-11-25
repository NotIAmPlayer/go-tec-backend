package main

import (
	"fmt"
	"go-tec-backend/config"
	"go-tec-backend/routes"
	"log"
	"os"

	"github.com/gin-contrib/cors"
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
	var origin string
	if gin.Mode() == gin.ReleaseMode {
		origin = frontend
	} else {
		origin = "http://localhost:5173"
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{origin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	routes.SetupRoutes(r)

	r.Run(":8080")
}
