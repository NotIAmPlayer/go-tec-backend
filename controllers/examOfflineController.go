package controllers

import (
	"fmt"
	"go-tec-backend/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CreateOfflineExamInput struct {
	ExamTitle     string   `json:"exam_title" binding:"required"`
	StartDatetime string   `json:"start_datetime" binding:"required"`
	EndDatetime   string   `json:"end_datetime" binding:"required"`
	RoomName      string   `json:"room_name" binding:"required"`
	Quota         int      `json:"quota" binding:"required"`
	Students      []string `json:"students"`
}

func CreateOfflineExam(c *gin.Context) {
	var input CreateOfflineExamInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")

	startTime, err := time.ParseInLocation("2006-01-02T15:04", input.StartDatetime, loc)
	if err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_datetime format"})
	return
	}

	endTime, err := time.ParseInLocation("2006-01-02T15:04", input.EndDatetime, loc)
	if err != nil {
	c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_datetime format"})
	return
	}

	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// 🔹 Generate custom exam_id (misal: OFF20251020_120000)
	examID := fmt.Sprintf("OFF%s", time.Now().Format("20060102_150405"))

	// 1️⃣ Insert ke `exam_offline`
	_, err = tx.Exec(`
		INSERT INTO exam_offline (id, exam_title, start_datetime, end_datetime, room_name, created_at)
		VALUES (?, ?, ?, ?, ?, NOW())
	`, examID, input.ExamTitle, startTime, endTime, input.RoomName)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert exam_offline"})
		return
	}

	// 2️⃣ Hitung sisa kuota
	studentCount := len(input.Students)
	available := input.Quota - studentCount
	if available < 0 {
		available = 0
	}

	// 3️⃣ Insert ke `kuota_ujian`
	_, err = tx.Exec(`
		INSERT INTO kuota_ujian (id, idUjian, total, available)
		VALUES (?, ?, ?, ?)
	`, fmt.Sprintf("KQ_%s", examID), examID, input.Quota, available)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert quota"})
		return
	}

	// 4️⃣ Insert ke `exam_offline_students`
	for _, nim := range input.Students {
		_, err := tx.Exec(`
			INSERT INTO exam_offline_students (id, exam_id, student_nim)
			VALUES (?, ?, ?)
		`, fmt.Sprintf("ST_%s_%s", examID, nim), examID, nim)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert student relation"})
			return
		}
	}

	// 5️⃣ Commit
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Offline exam created successfully",
		"exam_id":          examID,
		"quota_total":      input.Quota,
		"quota_available":  available,
		"student_inserted": studentCount,
	})
}
