package controllers

import (
	"database/sql"
	"fmt"
	"go-tec-backend/config"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type StudentExams struct {
	ExamID        int    `json:"exam_id"`
	ExamTitle     string `json:"exam_title"`
	StartDatetime string `json:"start_datetime"`
	EndDatetime   string `json:"end_datetime"`
}

type AnswerData struct {
	Nim        string `json:"nim"`
	ExamID     int    `json:"exam_id"`
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer"`
}

func GetUpcomingExams(c *gin.Context) {
	/*
		Gets the upcoming exam data for the current student.
	*/

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "401 - Unauthorized"})
		return
	}

	query := `
		SELECT u.idUjian, u.namaUjian, u.jadwalMulai, u.jadwalSelesai
		FROM ujian u JOIN ujian_ikut i ON u.idUjian = i.idUjian
		WHERE u.jadwalSelesai >= NOW() AND i.nim = ? AND (i.statusPengerjaan <> "selesai" OR ISNULL(i.statusPengerjaan))
		ORDER BY u.jadwalMulai ASC, u.jadwalSelesai ASC, u.idUjian ASC
	`

	rows, err := config.DB.Query(query, userID)

	if err != nil {
		log.Printf("Get multiple exams (student) error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal Server Error",
		})
		return
	}

	upcomingExams := []StudentExams{}

	defer rows.Close()

	for rows.Next() {
		var e StudentExams

		if err := rows.Scan(&e.ExamID, &e.ExamTitle, &e.StartDatetime, &e.EndDatetime); err != nil {
			log.Printf("Get multiple exams (student) error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal Server Error",
			})
			return
		}

		upcomingExams = append(upcomingExams, e)
	}

	if len(upcomingExams) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "200 - No exams found",
		})
		return
	} else {
		c.JSON(http.StatusOK, upcomingExams)
	}
}

func GetExamQuestions(c *gin.Context) {
	/*
		Uses the given exam ID to give out all of the questions.
	*/

	/*
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "400 - Invalid exam ID",
			})
			return
		}

		questions := []Question{}

		query := `
			SELECT u.idSoal, s.tipeSoal, s.isiSoal, s.pilihanA, s.pilihanB, s.pilihanC, s.pilihanD, s.audio
			FROM soal_ujian u JOIN soal s ON u.idSoal = s.idSoal
			WHERE u.idUjian = ?
			ORDER BY s.tipeSoal, u.idSoal
		`

		rows, err := config.DB.Query(query, id)

		if err != nil {
			log.Printf("Get exam questions (student) error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal Server Error",
			})
			return
		}

		defer rows.Close()

		for rows.Next() {
			var q Question

			var audio sql.NullString
			if err := rows.Scan(&q.QuestionID, &q.QuestionType, &q.QuestionText, &q.ChoiceA, &q.ChoiceB, &q.ChoiceC, &q.ChoiceD, &audio); err != nil {
				log.Printf("Get exam questions (student) error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "500 - Internal Server Error",
				})
				return
			}
			if audio.Valid {
				q.AudioPath = audio.String
			} else {
				q.AudioPath = ""
			}

			questions = append(questions, q)
		}

		if len(questions) == 0 {
			c.JSON(http.StatusOK, gin.H{
				"message": "200 - No questions found",
			})
			return
		} else {
			c.JSON(http.StatusOK, questions)
		}
	*/

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - work in progress",
	})
}

func AnswerExamQuestions(c *gin.Context) {
	/*
		Uses the given exam ID to give out all of the questions.
	*/

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid exam ID",
		})
		return
	}

	var a AnswerData

	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid JSON data",
		})
		return
	}

	if a.ExamID != id {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Exam ID mismatched",
		})
		return
	}

	var tempAnswer AnswerData

	// check if exists
	query := "SELECT nim, idSoal, idUjian, jawaban FROM soal_jawaban WHERE nim = ? AND idSoal = ? AND idUjian = ?"
	row := config.DB.QueryRow(query, a.Nim, a.QuestionID, a.ExamID)

	var query2 string

	if err := row.Scan(&tempAnswer.Nim, &tempAnswer.QuestionID, &tempAnswer.ExamID, &tempAnswer.Answer); err != nil {
		if err == sql.ErrNoRows {
			// uses insert into instead of update
			query2 = "INSERT INTO soal_jawaban (jawaban, nim, idSoal, idUjian) VALUES (?, ?, ?, ?)"
		} else {
			// something went wrong
			log.Printf("Get student answer error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal Server Error",
			})
			return
		}
	}

	// check if query is at default value, instead of filled like in earlier if there are no rows
	if query2 == "" {
		if a.Answer == "" {
			// uses delete from instead of update
			query2 = "DELETE FROM soal_jawaban WHERE nim = ? AND idSoal = ? AND idUjian = ?"
		} else {
			query2 = "UPDATE soal_jawaban SET jawaban = ? WHERE nim = ? AND idSoal = ? AND idUjian = ?"
		}
	}

	if strings.HasPrefix(query2, "DELETE") {
		_, err = config.DB.Exec(query2, a.Nim, a.QuestionID, a.ExamID)
	} else {
		_, err = config.DB.Exec(query2, a.Answer, a.Nim, a.QuestionID, a.ExamID)
	}

	if err != nil {
		log.Printf("Answer question error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - Exam updated successfully",
	})
}

func GetExamAnswers(c *gin.Context) {
	examID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid exam ID",
		})
		return
	}

	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "401 - Unauthorized"})
		return
	}

	answers := []AnswerData{}

	query := "SELECT idSoal, jawaban FROM soal_jawaban WHERE nim = ? AND idUjian = ?"
	rows, err := config.DB.Query(query, userID, examID)

	if err != nil {
		log.Printf("Get multiple answers error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal Server Error",
		})
		return
	}

	for rows.Next() {
		var a AnswerData

		if err := rows.Scan(&a.QuestionID, &a.Answer); err != nil {
			log.Printf("Get multiple answers error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal Server Error",
			})
			return
		}

		answers = append(answers, a)
	}

	if len(answers) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "200 - No answers found",
		})
		return
	} else {
		c.JSON(http.StatusOK, answers)
	}
}

func StartExamStudent(c *gin.Context) {
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "401 - Unauthorized"})
		return
	}

	// reuse, only doesn't need answer and question id
	var a AnswerData

	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid JSON data",
		})
		return
	}

	if a.Nim != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "401 - Unauthorized"})
		return
	}

	var e Exam

	query := "SELECT idUjian, namaUjian, jadwalMulai, jadwalSelesai FROM ujian WHERE idUjian = ?"
	row := config.DB.QueryRow(query, a.ExamID)

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

	target, err := time.Parse(time.DateTime, e.EndDatetime)
	if err != nil {
		log.Printf("Parse start datetime error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal Server Error",
		})
		return
	}

	currentDatetime := time.Now()

	duration := target.Sub(currentDatetime)

	fmt.Print(duration.Seconds(), " ", e.StartDatetime, " ", currentDatetime.Format(time.DateTime))

	if duration.Seconds() < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Exam has already ended",
		})
		return
	}

	query2 := "UPDATE ujian_ikut SET waktuMulai = ?, statusPengerjaan = ? WHERE nim = ? AND idUjian = ?"
	_, err = config.DB.Exec(query2, currentDatetime.Format(time.DateTime), "mengikuti", a.Nim, a.ExamID)

	if err != nil {
		log.Printf("Starting exam error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - Exam updated successfully",
	})
}

func EndExamStudent(c *gin.Context) {
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "401 - Unauthorized"})
		return
	}

	// reuse, only doesn't need answer
	var a AnswerData

	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid JSON data",
		})
		return
	}

	if a.Nim != userID {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "401 - Unauthorized"})
		return
	}

	var e Exam

	query := "SELECT idUjian, namaUjian, jadwalMulai, jadwalSelesai FROM ujian WHERE idUjian = ?"
	row := config.DB.QueryRow(query, a.ExamID)

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

	target, err := time.Parse(time.DateTime, e.StartDatetime)
	if err != nil {
		log.Printf("Parse start datetime error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal Server Error",
		})
		return
	}

	currentDatetime := time.Now()

	duration := target.Sub(currentDatetime)

	fmt.Print(duration.Seconds(), " ", e.StartDatetime, " ", currentDatetime.Format(time.DateTime))

	if duration.Seconds() < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Exam has not started yet",
		})
		return
	}

	query2 := "UPDATE ujian_ikut SET waktuSelesai = ?, statusPengerjaan = ? WHERE nim = ? AND idUjian = ?"
	_, err = config.DB.Exec(query2, currentDatetime.Format(time.DateTime), "selesai", a.Nim, a.ExamID)

	if err != nil {
		log.Printf("Starting exam error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - Exam updated successfully",
	})
}
