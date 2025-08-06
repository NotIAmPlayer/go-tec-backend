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

func GetOngoingExams(c *gin.Context) {
	/*
		Gets the upcoming exam data for the current student.
	*/

	query := `
		SELECT idUjian, namaUjian, jadwalMulai, jadwalSelesai FROM ujian
		WHERE jadwalMulai <= NOW() AND jadwalSelesai >= NOW()
		ORDER BY jadwalMulai ASC, jadwalSelesai ASC, idUjian ASC 
	`

	rows, err := config.DB.Query(query)

	if err != nil {
		log.Printf("Get multiple exams (student) error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal Server Error",
		})
		return
	}

	upcomingExams := []StudentExams{}

	defer rows.Close()

	for rows.Next() {
		var e StudentExams

		if err := rows.Scan(&e.ExamID, &e.ExamTitle, &e.StartDatetime, &e.EndDatetime); err != nil {
			log.Printf("Get multiple exams (student) error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal Server Error",
			})
			return
		}

		upcomingExams = append(upcomingExams, e)
	}

	if len(upcomingExams) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "200 - No exams found",
		})
		return
	} else {
		c.JSON(http.StatusOK, upcomingExams)
	}
}
