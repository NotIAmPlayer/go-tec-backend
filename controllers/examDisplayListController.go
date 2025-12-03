package controllers

import (
	"database/sql"
	"go-tec-backend/config"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type OfflineExamResponse struct {
	ExamID         string  `json:"exam_id"` // ubah dari int → string
	ExamTitle      string  `json:"exam_title"`
	StartDatetime  *string `json:"start_datetime"`
	EndDatetime    *string `json:"end_datetime"`
	RoomName       string  `json:"room_name"`
	StudentCount   int     `json:"student_count"`
	TotalQuota     int     `json:"total_quota"`
	AvailableQuota int     `json:"available_quota"`
}

// Helper untuk konversi waktu
func convertToWIB(t sql.NullTime) *string {
	if !t.Valid {
		return nil
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		log.Println("Error loading location Asia/Jakarta:", err)
		// Fallback ke string UTC jika gagal load
		utcTime := t.Time.Format("2006-01-02 15:04:05")
		return &utcTime
	}

	wibTime := t.Time.In(loc).Format("2006-01-02 15:04:05")
	return &wibTime
}

func GetOfflineExams(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT 
			e.id AS exam_id,
			e.exam_title,
			e.start_datetime,
			e.end_datetime,
			e.room_name,
			COALESCE(COUNT(s.student_nim), 0) AS student_count,
			COALESCE(k.total, 0) AS total_quota,
			COALESCE(k.available, 0) AS available_quota
		FROM exam_offline e
		LEFT JOIN exam_offline_students s ON e.id = s.exam_id
		LEFT JOIN kuota_ujian k ON e.id = k.idUjian
		GROUP BY e.id, e.exam_title, e.start_datetime, e.end_datetime, e.room_name, k.total, k.available
		ORDER BY e.created_at DESC;
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch offline exams"})
		return
	}
	defer rows.Close()

	var exams []OfflineExamResponse

	for rows.Next() {
		var exam OfflineExamResponse
		// Variabel Scan sudah benar, tidak perlu diubah
		err := rows.Scan(
			&exam.ExamID,
			&exam.ExamTitle,
			&exam.StartDatetime, // Ini akan menerima string waktu yg sudah dikonversi
			&exam.EndDatetime,   // Ini akan menerima string waktu yg sudah dikonversi
			&exam.RoomName,
			&exam.StudentCount,
			&exam.TotalQuota,
			&exam.AvailableQuota,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error scanning exam data"})
			return
		}
		exams = append(exams, exam)
	}

	if len(exams) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "No offline exams found"})
		return
	}

	c.JSON(http.StatusOK, exams)
}

func DeleteOfflineExam(c *gin.Context) {
	id := c.Param("id")

	tx, err := config.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// 1️⃣ Hapus relasi siswa
	_, err = tx.Exec(`DELETE FROM exam_offline_students WHERE exam_id = ?`, id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related students"})
		return
	}

	// 2️⃣ Hapus kuota ujian
	_, err = tx.Exec(`DELETE FROM kuota_ujian WHERE idUjian = ?`, id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete related quota"})
		return
	}

	// 3️⃣ Hapus data ujian offline utama
	res, err := tx.Exec(`DELETE FROM exam_offline WHERE id = ?`, id)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exam_offline"})
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Exam not found"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Offline exam deleted successfully"})
}

func GetOfflineExamByID(c *gin.Context) {
	id := c.Param("id")
	var exam OfflineExamResponse

	// 🔹 Ambil data utama ujian offline
	err := config.DB.QueryRow(`
		SELECT 
			e.id AS exam_id,
			e.exam_title,

			-- UBAH DI SINI: Konversi waktu ke Asia/Jakarta
			e.start_datetime,
			e.end_datetime,

			e.room_name,
			COALESCE(k.total, 0) AS total_quota,
			COALESCE(k.available, 0) AS available_quota
		FROM exam_offline e
		LEFT JOIN kuota_ujian k ON e.id = k.idUjian
		WHERE e.id = ?
	`, id).Scan(
		&exam.ExamID,
		&exam.ExamTitle,
		&exam.StartDatetime, // Ini akan menerima string waktu yg sudah dikonversi
		&exam.EndDatetime,   // Ini akan menerima string waktu yg sudah dikonversi
		&exam.RoomName,
		&exam.TotalQuota,
		&exam.AvailableQuota,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Offline exam not found"})
		return
	}

	// 🔹 Ambil daftar mahasiswa yang ikut ujian ini
	rows, err := config.DB.Query(`
	SELECT student_nim
	FROM exam_offline_students
	WHERE exam_id = ?
	`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students"})
		return
	}
	defer rows.Close()

	students := []string{}
	for rows.Next() {
		var nim string
		if err := rows.Scan(&nim); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan student_nim"})
			return
		}
		students = append(students, nim)
	}

	// 🔹 Gabungkan hasil semuanya
	c.JSON(http.StatusOK, gin.H{
		"exam_id":         exam.ExamID,
		"exam_title":      exam.ExamTitle,
		"start_datetime":  exam.StartDatetime, // Ini adalah string waktu WIB
		"end_datetime":    exam.EndDatetime,   // Ini adalah string waktu WIB
		"room_name":       exam.RoomName,
		"total_quota":     exam.TotalQuota,
		"available_quota": exam.AvailableQuota,
		"students":        students,
	})
}
