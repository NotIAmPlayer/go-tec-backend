package controllers

import (
	"net/http"
	"database/sql"
	"log"
	"time"
	"go-tec-backend/config" // ubah sesuai nama module kamu
	"github.com/gin-gonic/gin"
)

type Registration struct {
	ID           int    `json:"id"`
	Nim          string `json:"nim"`
	ExamType     string `json:"exam_type"`
	PaymentProof string `json:"payment_proof"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
}

func GetAllRegistrations(c *gin.Context) {
	rows, err := config.DB.Query(`
		SELECT id, nim, exam_type, payment_proof, status, created_at
		FROM pendaftaran_ujian
		ORDER BY created_at DESC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal mengambil data dari database",
			"error":   err.Error(),
		})
		return
	}
	defer rows.Close()

	var regs []Registration
	for rows.Next() {
		var r Registration
		if err := rows.Scan(&r.ID, &r.Nim, &r.ExamType, &r.PaymentProof, &r.Status, &r.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Gagal membaca data",
				"error":   err.Error(),
			})
			return
		}
		regs = append(regs, r)
	}

	c.JSON(http.StatusOK, regs)
}

func VerifyRegistration(c *gin.Context) {
	db := config.DB
	var input struct {
		ID     int    `json:"id"`
		Status string `json:"status"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Data tidak valid"})
		return
	}

	if input.Status != "approved" && input.Status != "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Status tidak valid"})
		return
	}

	// Ambil data pendaftaran
	var nim, examType, examID, prevStatus string
	err := db.QueryRow(`
		SELECT nim, exam_type, exam_id, status
		FROM pendaftaran_ujian
		WHERE id = ?
	`, input.ID).Scan(&nim, &examType, &examID, &prevStatus)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"message": "Data pendaftaran tidak ditemukan"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil data pendaftaran"})
		return
	}

	// 🟡 DEBUG: Tampilkan data awal
	log.Printf("DEBUG VerifyRegistration: ID=%d | NIM=%s | examType=%s | examID=%s | prevStatus=%s",
		input.ID, nim, examType, examID, prevStatus)

	// Update status di tabel pendaftaran_ujian
	_, err = db.Exec(`
		UPDATE pendaftaran_ujian
		SET status = ?
		WHERE id = ?
	`, input.Status, input.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memperbarui status pendaftaran"})
		return
	}

	// =====================================================
	// === OFFLINE SECTION =================================
	// =====================================================
	if input.Status == "approved" && examType == "offline" {
		// Kurangi kuota
		res, err := db.Exec(`
			UPDATE kuota_ujian
			SET available = available - 1
			WHERE idUjian = ? AND available > 0
		`, examID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memperbarui kuota ujian"})
			return
		}

		rows, _ := res.RowsAffected()
		if rows == 0 {
			log.Println("DEBUG Kuota penuh, mencoba mengalihkan ke ujian online...")

			now := time.Now()
			log.Println("DEBUG backend time:", now.Format("2006-01-02 15:04:05 MST"))

			var dbTime string
			errTime := db.QueryRow("SELECT NOW()").Scan(&dbTime)
			if errTime != nil {
				log.Println("DEBUG gagal ambil waktu database:", errTime)
			} else {
				log.Println("DEBUG database time (NOW()):", dbTime)
			}

			// Cari ujian online alternatif di tabel ujian
			var onlineExamID string
			err := db.QueryRow(`
				SELECT idUjian
				FROM ujian
				WHERE jadwalSelesai > NOW()
				ORDER BY jadwalMulai ASC
				LIMIT 1
			`).Scan(&onlineExamID)


			if err != nil {
				log.Printf("DEBUG Tidak menemukan ujian online: %v\n", err)
				c.JSON(http.StatusConflict, gin.H{"message": "Kuota offline penuh dan tidak ada ujian online tersedia"})
				return
			}

			// Update data pendaftaran agar diarahkan ke ujian online
			_, err = db.Exec(`
				UPDATE pendaftaran_ujian
				SET exam_type = 'online', exam_id = ?
				WHERE id = ?
			`, onlineExamID, input.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memperbarui data pendaftaran ke ujian online"})
				return
			}

			// Tambahkan ke ujian_ikut
			_, err = db.Exec(`
				INSERT INTO ujian_ikut (nim, idUjian, statusPengerjaan)
				VALUES (?, ?, 'belum_mulai')
			`, nim, onlineExamID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menambahkan mahasiswa ke ujian online"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Kuota ujian offline penuh, diarahkan ke ujian online"})
			return
		}

		// Tambahkan ke exam_offline_students
		var exists int
		_ = db.QueryRow(`
			SELECT COUNT(*) FROM exam_offline_students
			WHERE exam_id = ? AND student_nim = ?
		`, examID, nim).Scan(&exists)

		if exists == 0 {
			_, err = db.Exec(`
				INSERT INTO exam_offline_students (exam_id, student_nim)
				VALUES (?, ?)
			`, examID, nim)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menambahkan mahasiswa ke ujian offline"})
				return
			}
		}
	}

	// =====================================================
	// === ONLINE SECTION ==================================
	// =====================================================
	if input.Status == "approved" && examType == "online" {
		var exists int
		_ = db.QueryRow(`
			SELECT COUNT(*) FROM ujian_ikut
			WHERE nim = ? AND idUjian = ?
		`, nim, examID).Scan(&exists)
		log.Printf("DEBUG ONLINE exists=%d", exists)

		if exists == 0 {
			_, err = db.Exec(`
				INSERT INTO ujian_ikut (nim, idUjian, statusPengerjaan)
				VALUES (?, ?, 'belum_mulai')
			`, nim, examID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menambahkan mahasiswa ke ujian online"})
				return
			}
			log.Println("DEBUG Mahasiswa berhasil ditambahkan ke ujian online.")
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Status pendaftaran berhasil diperbarui"})
}

