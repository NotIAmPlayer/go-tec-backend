package controllers

import (
	"database/sql"
	"go-tec-backend/config"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Exams struct {
	ExamID        int      `json:"exam_id"`
	ExamTitle     string   `json:"exam_title"`
	StartDatetime string   `json:"start_datetime"`
	EndDatetime   string   `json:"end_datetime"`
	Questions     []int    `json:"questions"` // List of question IDs
	QuestionCount int      `json:"question_count"`
	Students      []string `json:"students"` // List of student NIMs
	StudentCount  int      `json:"student_count"`
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

		// amount of questions - no need to get details of every questions
		query2 := "SELECT COUNT(*) AS question_count FROM soal_ujian WHERE idUjian = ?"
		row2 := config.DB.QueryRow(query2, e.ExamID)

		if err := row2.Scan(&e.QuestionCount); err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "404 - Exam questions not found",
				})
			} else {
				log.Printf("Get exam question error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "500 - Internal Server Error",
				})
			}
			return
		}

		// amount of students
		query3 := "SELECT COUNT(*) AS student_count FROM ujian_ikut WHERE idUjian = ?"
		row3 := config.DB.QueryRow(query3, e.ExamID)

		if err := row3.Scan(&e.StudentCount); err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "404 - Exam questions not found",
				})
			} else {
				log.Printf("Get exam question error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "500 - Internal Server Error",
				})
			}
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
	/*
		Get a exam by ID from the database as JSON.
	*/
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

	// amount of questions - can be changed to list of questions and details later
	query2 := "SELECT COUNT(*) AS question_count FROM soal_ujian WHERE idUjian = ?"
	row2 := config.DB.QueryRow(query2, id)

	if err := row2.Scan(&e.QuestionCount); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "404 - Exam questions not found",
			})
		} else {
			log.Printf("Get exam question error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal Server Error",
			})
		}
		return
	}

	// amount of students - can be changed to list of students and details later
	query3 := "SELECT COUNT(*) AS student_count FROM ujian_ikut WHERE idUjian = ?"
	row3 := config.DB.QueryRow(query3, e.ExamID)

	if err := row3.Scan(&e.StudentCount); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "404 - Exam questions not found",
			})
		} else {
			log.Printf("Get exam question error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal Server Error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, e)
}

func CreateExam(c *gin.Context) {
	/*
		Create a new question in the database from JSON data.
	*/
	var e Exams

	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid JSON data",
		})
		return
	}

	createdTime := time.Now().Format("2025-05-26 08:47:59")

	query := "INSERT INTO ujian (idUjian, namaUjian, jadwalMulai, jadwalSelesai, dateCreated, dateUpdated) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := config.DB.Exec(query, e.ExamID, e.ExamTitle, e.StartDatetime, e.EndDatetime, createdTime, createdTime)

	if err != nil {
		log.Printf("Create question error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	// fetch the newly made exam ID
	query2 := "SELECT idUjian FROM ujian WHERE dateCreated = ?"
	row := config.DB.QueryRow(query2, createdTime)

	if err := row.Scan(&e.ExamID); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "404 - Exam not found",
			})
			return
		}
		log.Printf("Get question error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal Server Error",
		})
		return
	}

	// add valid question IDs to soal_ujian and student NIMs to ujian_ikut if questions and students arrays aren't empty.
	if len(e.Questions) > 0 {
		for _, qid := range e.Questions {
			// check if question with this id exists
			var q Questions
			query3 := "SELECT isiSoal FROM soal WHERE idSoal = ?"

			row := config.DB.QueryRow(query3, qid)

			if err := row.Scan(&q.QuestionText); err != nil {
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

			// if exists, place the question id for that exam
			query4 := "INSERT INTO soal_ujian (idUjian, idSoal) VALUES (?, ?)"
			_, err := config.DB.Exec(query4, e.ExamID, q)

			if err != nil {
				log.Printf("Create question error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "500 - Internal server error",
				})
				return
			}
		}
	}

	if len(e.Students) > 0 {
		for _, sid := range e.Students {
			// check if student with this npm exists
			var s Users
			query3 := "SELECT namaMhs FROM mahasiswa WHERE nim = ?"

			row := config.DB.QueryRow(query3, sid)

			if err := row.Scan(&s.Nama); err != nil {
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

			// if exists, place the student npm for that exam
			query4 := "INSERT INTO ujian_ikut (nim, idUjian) VALUES (?, ?)"
			_, err := config.DB.Exec(query4, sid, e.ExamID)

			if err != nil {
				log.Printf("Create question error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "500 - Internal server error",
				})
				return
			}
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "201 - Exam created successfully",
	})
}
