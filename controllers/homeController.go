package controllers

import (
	"go-tec-backend/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDashboardAdminData(c *gin.Context) {
	/*
		Gets the amounts of questions made, unpublished exam scores, and upcoming exams.
	*/

	var questionsMade, unpublishedExamScores, upcomingExams int

	query := `
		SELECT
			(SELECT COUNT(*) FROM soal) questionsMade,
			(SELECT COUNT(*) FROM ujian WHERE jadwalSelesai < CURRENT_TIMESTAMP() AND nilaiDiumumkan = 0) unpublishedExamScores,
			(SELECT COUNT(*) FROM ujian WHERE jadwalMulai > CURRENT_TIMESTAMP()) upcomingExams
	`

	row := config.DB.QueryRow(query)

	if err := row.Scan(&questionsMade, &unpublishedExamScores, &upcomingExams); err != nil {
		log.Printf("Get exam error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"questions_made":     questionsMade,
		"unpublished_scores": unpublishedExamScores,
		"upcoming_exams":     upcomingExams,
	})
}
