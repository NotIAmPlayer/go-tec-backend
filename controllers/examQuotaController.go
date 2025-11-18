package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-tec-backend/config" // ⬅️ ganti sesuai module kamu
)

// Struktur respons JSON
type QuotaResponse struct {
	Total     int `json:"total"`
	Available int `json:"available"`
}

// GetExamQuota mengambil data kuota dari database
func GetExamQuota(c *gin.Context) {
	db := config.DB // pastikan DB kamu bertipe *sql.DB

	var quota QuotaResponse

	// Contoh ambil kuota ujian pertama (atau kuota global)
	query := "SELECT total, available FROM kuota_ujian LIMIT 1"
	err := db.QueryRow(query).Scan(&quota.Total, &quota.Available)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"message": "Data kuota belum tersedia"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil data", "error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, quota)
}
