package routes

import (
	"github.com/gin-gonic/gin"

	"go-tec-backend/controllers"
	"go-tec-backend/middlewares"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/login", controllers.Login)

	r.GET("/audio/:filename", controllers.GetAudioFile)

	api := r.Group("/api")
	api.Use(middlewares.JWTAuthMiddleware()) // Apply JWT middleware to all routes in this group

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api.POST("/admin/:id/password", controllers.UpdateAdminPassword)

	api.GET("/users", controllers.GetAllUsers)
	api.GET("/users/count", controllers.GetUserCount)
	//api.GET("/users/page/:page", controllers.GetUsers) // unused, not needed for frontend
	api.GET("/users/:nim", controllers.GetUser)
	api.POST("/users", controllers.CreateUser)
	api.PUT("/users/:nim", controllers.UpdateUser)
	api.PUT("/users/:nim/password", controllers.UpdateUserPassword)
	api.DELETE("/users/:nim", controllers.DeleteUser)

	api.GET("/questions", controllers.GetAllQuestions)
	api.GET("/questions/count", controllers.GetQuestionCount)
	//api.GET("/questions/page/:page", controllers.GetQuestions) // unused, not needed for frontend
	api.GET("/questions/:id", controllers.GetQuestion)
	api.POST("/questions", controllers.CreateQuestion)
	api.PUT("/questions/:id", controllers.UpdateQuestion)
	api.DELETE("/questions/:id", controllers.DeleteQuestion)

	api.GET("/exams", controllers.GetAllExams)
	api.GET("/exams/count", controllers.GetExamCount)
	//api.GET("/exams/page/:page", controllers.GetExams) // unused, not needed for frontend
	api.GET("/exams/:id", controllers.GetExam)
	api.POST("/exams", controllers.CreateExam)
	api.PUT("/exams/:id", controllers.UpdateExam)
	api.DELETE("/exams/:id", controllers.DeleteExam)
}
