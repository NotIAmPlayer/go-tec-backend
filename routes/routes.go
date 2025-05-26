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
	api.PUT("/users/:nim", controllers.UpdateUser)
	api.PUT("/users/:nim/password", controllers.UpdateUserPassword)
	api.DELETE("/users/:nim", controllers.DeleteUser)

	//api.GET("/questions", controllers.GetAllQuestions) // In case you want to get all questions
	api.GET("/questions/page/:page", controllers.GetQuestions)
	api.GET("/questions/:id", controllers.GetQuestion)
	api.POST("/questions", controllers.CreateQuestion)
	api.PUT("/questions/:id", controllers.UpdateQuestion)
	api.DELETE("/questions/:id", controllers.DeleteQuestion)

	api.GET("/exams/page/:page", controllers.GetExams)
	api.GET("/exams/:id", controllers.GetExam)
}
