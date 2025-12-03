package controllers

import (
	"database/sql"
	"fmt"
	"go-tec-backend/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	//"go-tec-backend/models"
)

// Struktur request dari client
type LogRequest struct {
	Nim           string `json:"nim"`
	IdUjian       int    `json:"idUjian"`
	IdSoal        int    `json:"idSoal,omitempty"`
	TipeAktivitas string `json:"tipeAktivitas"`
	Aktivitas     string `json:"aktivitas"`
}

// Handler untuk menyimpan log aktivitas
func LogActivity(c *gin.Context) {
	var req LogRequest

	// Binding JSON request ke struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Query insert ke tabel log_aktivitas
	query := `
		INSERT INTO log_aktivitas (nim, idUjian, idSoal, tipeAktivitas, aktivitas, waktu)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	var soalID sql.NullInt64
	if req.IdSoal == 0 {
		soalID = sql.NullInt64{Int64: 0, Valid: false} // jadi NULL
	} else {
		soalID = sql.NullInt64{Int64: int64(req.IdSoal), Valid: true}
	}

	loc, _ := time.LoadLocation("Asia/Jakarta")
	now := time.Now().In(loc)

	_, err := config.DB.Exec(query,
		req.Nim,
		req.IdUjian,
		soalID,
		req.TipeAktivitas,
		req.Aktivitas,
		now.In(loc),
	)

	if err != nil {
		// kirimkan error asli supaya bisa debug
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "log saved"})
}

func GetScoresByExam(c *gin.Context) {
	examID := c.Param("examID")
	fmt.Println("Handler GetScoresByExam called with examID =", examID)

	query := `
        SELECT m.nim, m.namaMhs AS name, es.listening, es.grammar, es.reading, es.skor
        FROM hasil_ujian es
        JOIN mahasiswa m ON es.nim = m.nim
        WHERE es.idUjian = ?
    `
	rows, err := config.DB.Query(query, examID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "query failed",
			//"detail": err.Error(),
		})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var nim, name string
		var listening, grammar, reading int
		var score float64 // skor di DB decimal(5,2), jadi float64
		if err := rows.Scan(&nim, &name, &listening, &grammar, &reading, &score); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, map[string]interface{}{
			"nim":       nim,
			"name":      name,
			"listening": listening,
			"grammar":   grammar,
			"reading":   reading,
			"score":     score,
		})
	}

	c.JSON(http.StatusOK, results)
}

func GetLogsByExam(c *gin.Context) {
	examID := c.Param("examID")
	fmt.Println("Handler GetScoresByExam called with examID =", c.Param("examID"))
	query := `
        SELECT waktu, nim, tipeAktivitas, aktivitas
        FROM log_aktivitas
        WHERE idUjian = ?
        ORDER BY waktu DESC
    `
	rows, err := config.DB.Query(query, examID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var waktu string
		var nim, tipe, aktivitas string
		rows.Scan(&waktu, &nim, &tipe, &aktivitas)
		results = append(results, map[string]interface{}{
			"waktu":         waktu,
			"nim":           nim,
			"tipeAktivitas": tipe,
			"aktivitas":     aktivitas,
		})
	}
	c.JSON(http.StatusOK, results)
}

func GetLogsByStudent(c *gin.Context) {
	examID := c.Param("examID")
	nim := c.Param("nim")

	query := `
        SELECT waktu, nim, tipeAktivitas, aktivitas
        FROM log_aktivitas
        WHERE idUjian = ? AND nim = ?
        ORDER BY waktu DESC
    `
	rows, err := config.DB.Query(query, examID, nim)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var waktu, nimStr, tipe, aktivitas string
		if err := rows.Scan(&waktu, &nimStr, &tipe, &aktivitas); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, gin.H{
			"waktu":         waktu,
			"nim":           nimStr,
			"tipeAktivitas": tipe,
			"aktivitas":     aktivitas,
		})
	}
	c.JSON(http.StatusOK, results)
}

func DeleteLogsByStudent(c *gin.Context) {
	examID := c.Param("examID")
	nim := c.Param("nim")

	query := "DELETE FROM log_aktivitas WHERE idUjian = ? AND nim = ?"
	result, err := config.DB.Exec(query, examID, nim)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No logs found for this student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logs deleted successfully"})
}
