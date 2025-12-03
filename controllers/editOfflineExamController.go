package controllers

import (
	"go-tec-backend/config"
	"net/http"
	"time"
	"fmt"
	"github.com/gin-gonic/gin"
)

type UpdateOfflineExamInput struct {
	ExamTitle     string   `json:"exam_title" binding:"required"`
	StartDatetime string   `json:"start_datetime" binding:"required"`
	EndDatetime   string   `json:"end_datetime" binding:"required"`
	RoomName      string   `json:"room_name" binding:"required"`
	Quota         int      `json:"quota"`
	Students      []string `json:"students"`
}

// 🔧 Fungsi bantu: parsing waktu lokal tanpa geser ke UTC
func parseFlexibleTime(input string) (time.Time, error) {
	loc, _ := time.LoadLocation("Asia/Jakarta")
	layouts := []string{
		"2006-01-02T15:04",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
	}

	var t time.Time
	var err error
	for _, layout := range layouts {
		t, err = time.ParseInLocation(layout, input, loc)
		if err == nil {
			return t, nil
		}
	}
	if err != nil {
    	fmt.Println("parseFlexibleTime failed to parse:", input, "err:", err)
	}
	return time.Time{}, err

}



func UpdateOfflineExam(c *gin.Context) {
	id := c.Param("id")
	var input UpdateOfflineExamInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Parsing waktu dari FE (lokal Jakarta)
	startTime, err := parseFlexibleTime(input.StartDatetime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_datetime format"})
		return
	}
	endTime, err := parseFlexibleTime(input.EndDatetime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_datetime format"})
		return
	}

	

	// Pastikan tetap di zona Asia/Jakarta (tidak bergeser)
	loc, _ := time.LoadLocation("Asia/Jakarta")
	// ubah ke lokal WIB tanpa offset timezone (agar tidak bergeser di MySQL)
	startTimeLocal := time.Date(
		startTime.Year(), startTime.Month(), startTime.Day(),
		startTime.Hour(), startTime.Minute(), startTime.Second(), 0,
		loc,
	)

	fmt.Println("Received start:", input.StartDatetime)
	fmt.Println("Parsed as:", startTime)
	fmt.Println("Stored as:", startTimeLocal)
	endTimeLocal := time.Date(
		endTime.Year(), endTime.Month(), endTime.Day(),
		endTime.Hour(), endTime.Minute(), endTime.Second(), 0,
		loc,
	)

	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// 1️⃣ Update exam_offline
	_, err = tx.Exec(`
		UPDATE exam_offline
		SET exam_title = ?, start_datetime = ?, end_datetime = ?, room_name = ?
		WHERE id = ?
	`, input.ExamTitle, startTimeLocal, endTimeLocal, input.RoomName, id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exam_offline"})
		return
	}

	// 2️⃣ Upsert kuota
	studentCount := len(input.Students)
	totalQuota := input.Quota
	if totalQuota < studentCount {
		totalQuota = studentCount
	}
	available := totalQuota - studentCount

	_, err = tx.Exec(`
		INSERT INTO kuota_ujian (idUjian, total, available)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE
			total = VALUES(total),
			available = VALUES(available)
	`, id, totalQuota, available)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upsert exam quota"})
		return
	}

	// 3️⃣ Reset students
	_, err = tx.Exec(`DELETE FROM exam_offline_students WHERE exam_id = ?`, id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear old student relations"})
		return
	}

	for _, nim := range input.Students {
		_, err := tx.Exec(`
			INSERT INTO exam_offline_students (exam_id, student_nim)
			VALUES (?, ?)
		`, id, nim)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert new student relation"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Offline exam updated successfully",
		"exam_id":         id,
		"quota_total":     totalQuota,
		"quota_available": available,
		"student_count":   studentCount,
	})
}
