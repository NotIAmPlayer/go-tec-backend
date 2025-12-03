package routes

import (
	"go-tec-backend/controllers"
	"go-tec-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Public routes - no authentication required
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.Login)
	r.POST("/refresh", controllers.Refresh)
	r.POST("/logout", controllers.Logout)

	r.GET("/audio/:filename", controllers.GetAudioFile)

	r.POST("/forgot", controllers.ForgotPassword)
	r.POST("/forgot/reset", controllers.ResetPassword)

	r.GET("/ws", controllers.ExamWebsocket)
	r.GET("/api/exam/quota", controllers.GetExamQuota)
	r.POST("/api/exam/register", controllers.RegisterExam)

	r.Static("/uploads", "./uploads")

	// Simulate an exam answer by a test student to get all correct answers
	// Do NOT leave this on production (please)
	r.POST("/debug/:idUjian", controllers.DebugExamAnswers)

	// Main bulk of API routes - requires JWT authentication
	api := r.Group("/api")
	api.Use(middlewares.JWTAuthMiddleware()) // Apply JWT middleware to all routes in this group

	api.GET("/me", controllers.GetMe)

	api.PUT("/self/password", controllers.UpdateSelfPassword)

	student := api.Group("/student")
	student.Use(middlewares.StudentMiddleware())

	student.GET("/home", controllers.GetUpcomingExams)

	student.PUT("/exam/start", controllers.StartExamStudent)
	student.PUT("/exam/finish", controllers.EndExamStudent)
	student.GET("/exam/:id", controllers.GetExamQuestions)
	student.POST("/exam/:id", controllers.AnswerExamQuestions)
	student.POST("/exam/submit/:nim/:idUjian", controllers.SubmitExam)

	student.GET("/answers/:id", controllers.GetExamAnswers)

	student.POST("/logs", controllers.LogActivity)

	student.GET("/offline", controllers.GetOfflineExamsForStudent) // ✅ baru
	student.GET("/offline/available", controllers.GetAvailableOfflineExams)

	student.GET("/online/available", controllers.GetAvailableOnlineExams)

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

	admin.GET("/registrations", controllers.GetAllRegistrations)
	admin.POST("/registrations/verify", controllers.VerifyRegistration)

	admin.POST("/exams/offline", controllers.CreateOfflineExam)
	admin.GET("/exams/offline", controllers.GetOfflineExams)
	admin.GET("/exams/offline/:id", controllers.GetOfflineExamByID)
	admin.PUT("/exams/offline/:id", controllers.UpdateOfflineExam)
	admin.DELETE("/exams/offline/:id", controllers.DeleteOfflineExam)
	admin.GET("/exams/:id/students", controllers.GetExamStudents)

	//admin.GET("/scores/:examID", controllers.GetScoresByExam)
	admin.GET("/logs/:examID", controllers.GetLogsByExam)
	admin.GET("/logs/:examID/:nim", controllers.GetLogsByStudent)
	admin.DELETE("/logs/:examID/:nim", controllers.DeleteLogsByStudent)
	admin.GET("/scores/:examID", controllers.GetScoresByExam)
}
