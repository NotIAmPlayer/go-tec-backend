package controllers

import (
	"database/sql"
	"go-tec-backend/config"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Exams struct {
	ExamID        int      `json:"exam_id"`
	ExamTitle     string   `json:"exam_title"`
	StartDatetime string   `json:"start_datetime"`
	EndDatetime   string   `json:"end_datetime"`
	Questions     []int    `json:"questions"` // List of question IDs
	Students      []string `json:"students"`  // List of student NIMs
}

func GetExams(c *gin.Context) {
	/*
		Get exams on a specific page from the database as JSON.
	*/

	page, err := strconv.Atoi(c.Param("page"))

	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid page number",
		})
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid page size number",
		})
		return
	}

	var offset = (page - 1) * limit

	exams := []Exams{}

	query := "SELECT idUjian, namaUjian, jadwalMulai, jadwalSelesai FROM ujian LIMIT ? OFFSET ?"
	rows, err := config.DB.Query(query, limit, offset)

	if err != nil {
		log.Printf("Get multiple questions error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal Server Error",
		})
		return
	}

	for rows.Next() {
		var e Exams

		if err := rows.Scan(&e.ExamID, &e.ExamTitle, &e.StartDatetime, &e.EndDatetime); err != nil {
			log.Printf("Get multiple exams error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal Server Error",
			})
			return
		}

		exams = append(exams, e)
	}

	if len(exams) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "404 - No exams found",
		})
		return
	} else {
		c.JSON(http.StatusOK, exams)
	}
}

func GetExam(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid exam ID",
		})
		return
	}

	var e Exams

	query := "SELECT idUjian, namaUjian, jadwalMulai, jadwalSelesai FROM ujian WHERE idUjian = ?"
	row := config.DB.QueryRow(query, id)

	if err := row.Scan(&e.ExamID, &e.ExamTitle, &e.StartDatetime, &e.EndDatetime); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "404 - Exam not found",
			})
		} else {
			log.Printf("Get exam error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal Server Error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, e)
}
