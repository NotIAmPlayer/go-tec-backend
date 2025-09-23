package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go-tec-backend/config"
)

// Mapping jumlah benar ke skor konversi (dari tabel client)
var scoreConversion = map[int]int{
	50: 68, 49: 67, 48: 66, 47: 65, 46: 63,
	45: 62, 44: 61, 43: 60, 42: 59, 41: 58,
	40: 57, 39: 57, 38: 56, 37: 55, 36: 54,
	35: 54, 34: 53, 33: 52, 32: 52, 31: 51,
	30: 51, 29: 50, 28: 49, 27: 49, 26: 48,
	25: 48, 24: 47, 23: 47, 22: 46, 21: 45,
	20: 45, 19: 44, 18: 43, 17: 42, 16: 41,
	15: 41, 14: 39, 13: 38, 12: 37, 11: 35,
	10: 33, 9: 32, 8: 32, 7: 31, 6: 30,
	5: 29, 4: 28, 3: 27, 2: 26, 1: 25,
	0: 24,
}

// Fungsi konversi jumlah benar ke skor
func ConvertScore(correct int) int {
	if val, ok := scoreConversion[correct]; ok {
		return val
	}
	return 0
}

// Controller untuk submit ujian dan hitung skor otomatis
func SubmitExam(c *gin.Context) {
	nim := c.Param("nim")
	idUjian := c.Param("idUjian")

	idUjianInt, err := strconv.Atoi(idUjian)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "idUjian harus berupa angka"})
		return
	}

	sections := []int{1, 2, 3}
	convertedScores := []int{}
	totalConverted := 0

	for _, batch := range sections {
		query := `
			SELECT COUNT(*)
			FROM soal_jawaban sj
			JOIN soal s ON sj.idSoal = s.idSoal
			WHERE sj.nim = ? AND sj.idUjian = ? AND s.idBatch = ? AND sj.jawaban = s.jawaban
		`
		var correct int
		err := config.DB.QueryRow(query, nim, idUjianInt, batch).Scan(&correct)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		converted := ConvertScore(correct)
		convertedScores = append(convertedScores, converted)
		totalConverted += converted
	}

	finalScore := float64(totalConverted) / 3.0 * 10.0

	// Simpan atau update skor
	_, err = config.DB.Exec(`
		INSERT INTO hasil_ujian (nim, idUjian, skor)
		VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE skor = VALUES(skor)
	`, nim, idUjianInt, finalScore)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "gagal menyimpan hasil ujian",
			"detail": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nim":           nim,
		"idUjian":       idUjianInt,
		"sectionScores": convertedScores,
		"finalScore":    finalScore,
		"message":       "nilai berhasil dihitung dan disimpan",
	})
}
