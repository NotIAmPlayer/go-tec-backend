package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go-tec-backend/config"
)

// Struct untuk response ke frontend
type ScoreResponse struct {
	Nim   string  `json:"nim"`
	Name  string  `json:"name"`
	Score float64 `json:"score"`
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
