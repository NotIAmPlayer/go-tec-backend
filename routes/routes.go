package routes

import (
	"github.com/gin-gonic/gin"
	"go-tec-backend/controllers"
	"go-tec-backend/middlewares"
	"fmt"
)

func SetupRoutes(r *gin.Engine) {
	// Public routes - no authentication required
	r.POST("/login", controllers.Login)
	r.POST("/register", controllers.RegisterUser)

	r.GET("/audio/:filename", controllers.GetAudioFile)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/ws", controllers.ExamWebsocket)

	// Main bulk of API routes - requires JWT authentication
	api := r.Group("/api")
	api.Use(middlewares.JWTAuthMiddleware()) // Apply JWT middleware to all routes in this group

	api.GET("/me", controllers.GetMe)

	api.GET("/self/password", controllers.UpdateSelfPassword)

	student := api.Group("/student")
	student.Use(middlewares.StudentMiddleware())

	student.GET("/home", controllers.GetUpcomingExams)
	student.PUT("/exam/start", controllers.StartExamStudent)
	student.PUT("/exam/finish", controllers.EndExamStudent)
	student.GET("/exam/:id", controllers.GetExamQuestions)
	student.POST("/exam/:id", controllers.AnswerExamQuestions)
	student.GET("/answers/:id", controllers.GetExamAnswers)
	student.POST("/logs", controllers.LogActivity)
	student.POST("/exam/submit/:nim/:idUjian", controllers.SubmitExam)

	admin := api.Group("/admin")
	admin.Use(middlewares.AdminMiddleware())

	admin.GET("/home", controllers.GetDashboardAdminData)
	admin.GET("/password/:id/", controllers.UpdateAdminPassword)
	admin.GET("/home/ongoing", controllers.GetOngoingExams)

	admin.GET("/users", controllers.GetAllUsers)
	admin.GET("/users/:nim", controllers.GetUser)
	admin.POST("/users", controllers.CreateUser)
	admin.PUT("/users/:nim", controllers.UpdateUser)
	admin.PUT("/users/:nim/password", controllers.UpdateUserPassword)
	admin.DELETE("/users/:nim", controllers.DeleteUser)

	admin.GET("/questions", controllers.GetAllQuestionBatches)
	admin.GET("/questions/:id", controllers.GetQuestion)
	admin.POST("/questions", controllers.CreateQuestionBatch)
	admin.PUT("/questions/:id", controllers.UpdateQuestionBatch)
	admin.DELETE("/questions/:id", controllers.DeleteQuestionBatch)

	admin.GET("/exams", controllers.GetAllExams)
	admin.GET("/exams/:id", controllers.GetExam)
	admin.POST("/exams", controllers.CreateExam)
	admin.PUT("/exams/:id", controllers.UpdateExam)
	admin.DELETE("/exams/:id", controllers.DeleteExam)

	admin.GET("/scores/:examID", controllers.GetScoresByExam)
	admin.GET("/logs/:examID", controllers.GetLogsByExam)
	admin.GET("/logs/:examID/:nim", controllers.GetLogsByStudent)
	admin.DELETE("/logs/:examID/:nim", controllers.DeleteLogsByStudent)
	admin.GET("/scores/:examID/", controllers.GetScoresByExam)

	
	for _, ri := range r.Routes() {
		fmt.Println("Registered route:", ri.Method, ri.Path)
	}
	
	
}
