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

	//api.GET("/users", controllers.GetAllUsers) // In case you want to get all users
	api.GET("/users/page/:page", controllers.GetUsers)
	api.GET("/users/:nim", controllers.GetUser)
	api.POST("/users", controllers.CreateUser)
}
