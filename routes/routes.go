package routes

import (
	"github.com/gin-gonic/gin"

	"go-tec-backend/controllers"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api.GET("/users", controllers.GetUsers)
}
