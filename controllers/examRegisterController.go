package controllers

import (
	"database/sql"
	"net/http"
	"path/filepath"
	"time"
	"github.com/gin-gonic/gin"
	"go-tec-backend/config"
)

func RegisterExam(c *gin.Context) {
	db := config.DB // koneksi *sql.DB

	nim := c.PostForm("nim")
	examType := c.PostForm("exam_type")
	exam_id := c.PostForm("exam_id")

	if nim == "" || examType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Data tidak lengkap"})
		return
	}

	// Kalau tipe ujian offline, pastikan exam_id dikirim
	if examType == "offline" && exam_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Exam ID wajib diisi untuk ujian offline"})
		return
	}

	file, err := c.FormFile("payment_proof")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "File bukti pembayaran wajib diunggah"})
		return
	}

	// Simpan file bukti pembayaran
	filename := time.Now().Format("20060102150405") + "_" + filepath.Base(file.Filename)
	savePath := filepath.Join("uploads", filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menyimpan file"})
		return
	}

	// ✅ Cegah mahasiswa daftar dua kali di exam offline yang sama
	if examType == "offline" {
		var existingCount int
		err = db.QueryRow(`
			SELECT COUNT(*) 
			FROM pendaftaran_ujian 
			WHERE nim = ? AND exam_type = 'offline' AND exam_id = ? AND status IN ('pending', 'approved')
		`, nim, exam_id).Scan(&existingCount)

		if err != nil && err != sql.ErrNoRows {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memeriksa data duplikat"})
			return
		}
		if existingCount > 0 {
			c.JSON(http.StatusConflict, gin.H{"message": "Anda sudah terdaftar di ujian offline ini"})
			return
		}
	}

	// ✅ Simpan data pendaftaran (status = pending)
	query := `
		INSERT INTO pendaftaran_ujian (nim, exam_type, exam_id, payment_proof, status, created_at)
		VALUES (?, ?, ?, ?, ?, NOW())
	`
	_, err = db.Exec(query, nim, examType, exam_id, filename, "pending")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal menyimpan ke database",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Pendaftaran ujian berhasil dikirim dan menunggu approval admin.",
	})
}