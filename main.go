package main

import (
	"fmt"
	"go-tec-backend/config"
	"go-tec-backend/routes"
	"log"
	"os"
	"strings"
	"time"

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
	var origin []string
	if gin.Mode() == gin.ReleaseMode {
		if frontend == "" {
			log.Fatal("FRONTEND_URL environment variable cannot be an empty string")
		}

		origin = strings.Split(frontend, ",")

		for i := range origin {
			origin[i] = strings.TrimSpace(origin[i])
		}
	} else {
		origin = []string{"http://localhost:5173"}
	}

	originStr := strings.Join(origin, ", ")
	fmt.Println("origin:" + originStr)

	r.Use(func(c *gin.Context) {
		originHeader := c.GetHeader("Origin")
		log.Printf("Received Origin: %s", originHeader)
		log.Printf("Allowed Origin: %s", origin)
		c.Next()
	})

	r.Use(cors.New(cors.Config{
		AllowOrigins:     origin,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.SetupRoutes(r)

	r.Run(":8080")
}
