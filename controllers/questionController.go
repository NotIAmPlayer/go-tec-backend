package controllers

import (
	"database/sql"
	"go-tec-backend/config"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Questions struct {
	QuestionID   int    `json:"question_id"`
	QuestionType string `json:"question_type"`
	QuestionText string `json:"question_text"`
	ChoiceA      string `json:"choice_a"`
	ChoiceB      string `json:"choice_b"`
	ChoiceC      string `json:"choice_c"`
	ChoiceD      string `json:"choice_d"`
	Answer       string `json:"answer"`
}

func GetQuestions(c *gin.Context) {
	/*
		Get questions on a specific page from the database as JSON.
	*/

	page, err := strconv.Atoi(c.Param("page"))

	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid page number",
		})
		return
	}

	var limit = 20
	var offset = (page - 1) * limit

	questions := []Questions{}

	query := "SELECT id_soal, tipe_soal, isi_soal, pilihan_1, pilihan_2, pilihan_3, pilihan_4, kunci_jawaban FROM soal LIMIT ? OFFSET ?"

	rows, err := config.DB.Query(query, limit, offset)

	if err != nil {
		log.Printf("Get multiple questions error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal Server Error",
		})
		return
	}

	for rows.Next() {
		var q Questions

		if err := rows.Scan(&q.QuestionID, &q.QuestionType, &q.QuestionText, &q.ChoiceA, &q.ChoiceB, &q.ChoiceC, &q.ChoiceD, &q.Answer); err != nil {
			log.Printf("Get multiple questions error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal Server Error",
			})
			return
		}

		questions = append(questions, q)
	}

	if len(questions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "404 - No questions found on page " + strconv.Itoa(page),
		})
		return
	} else {
		c.JSON(http.StatusOK, questions)
	}
}

func GetQuestion(c *gin.Context) {
	/*
		Get a question by ID from the database as JSON.
	*/
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid question ID",
		})
		return
	}

	var q Questions

	query := "SELECT id_soal, tipe_soal, isi_soal, pilihan_1, pilihan_2, pilihan_3, pilihan_4, kunci_jawaban FROM soal WHERE id_soal = ?"

	row := config.DB.QueryRow(query, id)

	if err := row.Scan(&q.QuestionID, &q.QuestionType, &q.QuestionText, &q.ChoiceA, &q.ChoiceB, &q.ChoiceC, &q.ChoiceD, &q.Answer); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "404 - Question not found",
			})
			return
		}
		log.Printf("Get question error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal Server Error",
		})
		return
	}

	c.JSON(http.StatusOK, q)
}

func CreateQuestion(c *gin.Context) {
	/*
		Create a new question in the database from JSON data.
	*/
	var q Questions

	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid JSON data",
		})
		return
	}

	if q.QuestionType != "Listening" && q.QuestionType != "Reading" && q.QuestionType != "Grammar" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid question type",
		})
		return
	}

	query := "INSERT INTO soal (tipe_soal, isi_soal, pilihan_1, pilihan_2, pilihan_3, pilihan_4, kunci_jawaban) VALUES (?, ?, ?, ?, ?, ?, ?)"
	_, err := config.DB.Exec(query, q.QuestionType, q.QuestionText, q.ChoiceA, q.ChoiceB, q.ChoiceC, q.ChoiceD, q.Answer)

	if err != nil {
		log.Printf("Create question error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "201 - Question created successfully",
	})
}

func UpdateQuestion(c *gin.Context) {
	/*
		Update a question in the database from JSON data.
	*/
	id := c.Param("id")

	var q Questions

	if err := c.ShouldBindJSON(&q); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid JSON data",
		})
		return
	}

	updates := []string{}
	args := []interface{}{}

	if q.QuestionType != "" && (q.QuestionType == "Listening" || q.QuestionType == "Reading" || q.QuestionType == "Grammar") {
		updates = append(updates, "tipe_soal = ?")
		args = append(args, q.QuestionType)
	}

	if q.QuestionText != "" {
		updates = append(updates, "isi_soal = ?")
		args = append(args, q.QuestionText)
	}

	if q.ChoiceA != "" {
		updates = append(updates, "pilihan_1 = ?")
		args = append(args, q.ChoiceA)
	}

	if q.ChoiceB != "" {
		updates = append(updates, "pilihan_2 = ?")
		args = append(args, q.ChoiceB)
	}

	if q.ChoiceC != "" {
		updates = append(updates, "pilihan_3 = ?")
		args = append(args, q.ChoiceC)
	}

	if q.ChoiceD != "" {
		updates = append(updates, "pilihan_4 = ?")
		args = append(args, q.ChoiceD)
	}

	if q.Answer != "" {
		updates = append(updates, "kunci_jawaban = ?")
		args = append(args, q.Answer)
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - No fields to update",
		})
		return
	}

	args = append(args, id)
	query := "UPDATE soal SET " + strings.Join(updates, ", ") + " WHERE id_soal = ?"

	_, err := config.DB.Exec(query, args...)

	if err != nil {
		log.Printf("Update question error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - Question updated successfully",
	})
}

func DeleteQuestion(c *gin.Context) {
	/*
		Delete a question from the database.
	*/
	id := c.Param("id")

	query := "DELETE FROM soal WHERE id_soal = ?"
	_, err := config.DB.Exec(query, id)

	if err != nil {
		log.Printf("Delete question error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - Question deleted successfully",
	})
}
