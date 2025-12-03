package controllers

import (
	"net/http"

	"go-tec-backend/config"

	"github.com/gin-gonic/gin"
)

// Struct untuk response ke frontend
type ScoreResponse struct {
	Nim   string  `json:"nim"`
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

func DebugExamAnswers(c *gin.Context) {
	// only use this on localhost
	/*
		idUjian := c.Param("idUjian")

		query := `
			SELECT s.idSoal, u.idUjian, b.tipeBatch, s.jawaban
			FROM batch_soal b
				JOIN batch_ujian u ON b.idBatch = u.idBatch
				JOIN soal s ON b.idBatch = s.idBatch
			WHERE u.idUjian = ?
		`

		rows, err := config.DB.Query(query, idUjian)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var questionID, examID int
			var batchType, answer string

			if err := rows.Scan(&questionID, &examID, &batchType, &answer); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			query2 := `
				INSERT INTO soal_jawaban (nim, idSoal, idUjian, tipeBatch, jawaban)
				VALUES (?, ?, ?, ?, ?)
			`

			var batchNumber int

			if batchType == "listening" {
				batchNumber = 1
			} else if batchType == "grammar" {
				batchNumber = 2
			} else {
				batchNumber = 3
			}

			_, err := config.DB.Exec(query2, "223400005", questionID, examID, batchNumber, answer)

			if err != nil {
				log.Printf("Create exam error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "500 - Internal server error",
				})
				return
			}
		}
	*/

	c.JSON(http.StatusOK, gin.H{"message": "done!"})
}

// Endpoint: GET /api/admin/scores/:idUjian
func GetExamScores(c *gin.Context) {
	idUjian := c.Param("idUjian")

	query := `
		SELECT h.nim, m.namaMhs, h.skor
		FROM hasil_ujian h
		JOIN mahasiswa m ON h.nim = m.nim
		WHERE h.idUjian = ?
		ORDER BY h.skor DESC
	`

	rows, err := config.DB.Query(query, idUjian)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var results []ScoreResponse
	for rows.Next() {
		var r ScoreResponse
		if err := rows.Scan(&r.Nim, &r.Name, &r.Score); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		results = append(results, r)
	}

	c.JSON(http.StatusOK, results)
}
