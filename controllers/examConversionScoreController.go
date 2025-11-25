package controllers

import (
	"math"
	"net/http"
	"strconv"

	"go-tec-backend/config"

	"github.com/gin-gonic/gin"
)

// Mapping jumlah benar ke skor konversi (dari tabel client)
// scoreConversion[section][jumlahBenar] = skor konversi
var scoreConversion = map[int]map[int]int{
	1: { // Section 1
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
	},
	2: { // Section 2
		40: 68, 39: 67, 38: 65, 37: 63, 36: 61,
		35: 60, 34: 58, 33: 57, 32: 56, 31: 55,
		30: 54, 29: 53, 28: 52, 27: 51, 26: 50,
		25: 49, 24: 48, 23: 47, 22: 46, 21: 45,
		20: 44, 19: 43, 18: 42, 17: 41, 16: 40,
		15: 40, 14: 38, 13: 37, 12: 36, 11: 36,
		10: 33, 9: 31, 8: 29, 7: 27, 6: 26,
		5: 25, 4: 24, 3: 23, 2: 22, 1: 21,
		0: 20,
	},
	3: { // Section 3
		50: 67, 49: 66, 48: 65, 47: 63, 46: 61,
		45: 60, 44: 59, 43: 58, 42: 57, 41: 56,
		40: 55, 39: 54, 38: 54, 37: 53, 36: 52,
		35: 52, 34: 51, 33: 50, 32: 49, 31: 48,
		30: 48, 29: 47, 28: 46, 27: 46, 26: 45,
		25: 44, 24: 43, 23: 43, 22: 42, 21: 41,
		20: 40, 19: 39, 18: 38, 17: 37, 16: 36,
		15: 35, 14: 34, 13: 32, 12: 31, 11: 30,
		10: 29, 9: 28, 8: 28, 7: 27, 6: 26,
		5: 25, 4: 24, 3: 24, 2: 23, 1: 22,
		0: 21,
	},
}

// Fungsi konversi jumlah benar ke skor
func ConvertScore(section int, correct int) int {
	if val, ok := scoreConversion[section][correct]; ok {
		return val
	}
	return 0
}

func round(x float64) float64 {
	mod := math.Mod(x, 1)

	if mod >= 0.5 {
		return math.Ceil(x)
	} else {
		return math.Floor(x)
	}
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

	sections := []int{1, 2, 3} // tipeBatch

	convertedScores := []int{}
	totalConverted := 0

	for _, section := range sections {
		query := `
			SELECT COUNT(*)
			FROM soal_jawaban sj
			JOIN soal s ON sj.idSoal = s.idSoal
			WHERE sj.nim = ? AND sj.idUjian = ? 
				AND sj.tipeBatch = ? 
				AND sj.jawaban = s.jawaban
		`
		var correct int
		err := config.DB.QueryRow(query, nim, idUjianInt, section).Scan(&correct)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		converted := ConvertScore(section, correct)
		convertedScores = append(convertedScores, converted)
		totalConverted += converted
	}

	finalScore := float64(totalConverted) / 3.0 * 10.
	finalScore = round(finalScore)

	// Simpan atau update skor
	_, err = config.DB.Exec(`
		INSERT INTO hasil_ujian (nim, idUjian, listening, grammar, reading, skor)
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE skor = VALUES(skor)
	`, nim, idUjianInt, convertedScores[0], convertedScores[1], convertedScores[2], finalScore)
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
