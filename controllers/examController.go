package controllers

import (
	"database/sql"
	"go-tec-backend/config"
	"log"
	"net/http"
	"strconv"
	"strings"
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

func GetExamCount(c *gin.Context) {
	/*
		Get the total number of exams in the database.
		Used by the frontend for pagination.
	*/
	var count int

	query := "SELECT COUNT(*) FROM ujian"

	row := config.DB.QueryRow(query)

	if err := row.Scan(&count); err != nil {
		log.Printf("Get exam count error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": count,
	})
}

func GetAllExams(c *gin.Context) {
	/*
		Get all exams from the database as JSON.
	*/
	exams := []Exams{}

	query := "SELECT idUjian, namaUjian, jadwalMulai, jadwalSelesai FROM ujian ORDER BY idUjian ASC"
	rows, err := config.DB.Query(query)

	if err != nil {
		log.Printf("Get multiple questions error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal Server Error",
		})
		return
	}

	defer rows.Close()

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
		c.JSON(http.StatusOK, gin.H{
			"message": "200 - No exams found",
		})
		return
	} else {
		c.JSON(http.StatusOK, exams)
	}
}

/* Unused
func GetExams(c *gin.Context) {
	//	Get exams on a specific page from the database as JSON.

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

	query := "SELECT idUjian, namaUjian, jadwalMulai, jadwalSelesai FROM ujian ORDER BY idSoal ASC LIMIT ? OFFSET ?"
	rows, err := config.DB.Query(query, limit, offset)

	if err != nil {
		log.Printf("Get multiple questions error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal Server Error",
		})
		return
	}

	defer rows.Close()

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
*/

func GetExam(c *gin.Context) {
	/*
		Get an exam by ID from the database as JSON.
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
				"message": "404 - Exam students not found",
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
		Create a new exam in the database from JSON data.
	*/
	var e Exams

	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid JSON data",
		})
		return
	}

	createdTime := time.Now().Format(time.DateTime)

	query := "INSERT INTO ujian (idUjian, namaUjian, jadwalMulai, jadwalSelesai, dateCreated, dateUpdated) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := config.DB.Exec(query, e.ExamID, e.ExamTitle, e.StartDatetime, e.EndDatetime, createdTime, createdTime)

	if err != nil {
		log.Printf("Create exam error: %v", err)
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
		log.Printf("Get exam error: %v", err)
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
			_, err := config.DB.Exec(query4, e.ExamID, qid)

			if err != nil {
				log.Printf("Create exam question error: %v", err)
				// log.Println(e.ExamID, qid)
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
						"message": "404 - Student not found",
					})
					return
				}
				log.Printf("Get student error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "500 - Internal Server Error",
				})
				return
			}

			// if exists, place the student npm for that exam
			query4 := "INSERT INTO ujian_ikut (nim, idUjian) VALUES (?, ?)"
			_, err := config.DB.Exec(query4, sid, e.ExamID)

			if err != nil {
				log.Printf("Create exam student error: %v", err)
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

func UpdateExam(c *gin.Context) {
	/*
		Update an exam in the database from JSON data.
	*/
	id := c.Param("id")

	var e Exams

	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid JSON data",
		})
		return
	}

	updates := []string{}
	args := []interface{}{}

	if e.ExamTitle != "" {
		updates = append(updates, "namaUjian = ?")
		args = append(args, e.ExamTitle)
	}

	if e.StartDatetime != "" {
		// determine if it's a valid datetime
		if _, err := time.Parse(time.DateTime, e.StartDatetime); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "400 - Invalid starting date and time",
			})
			return
		}

		updates = append(updates, "jadwalMulai = ?")
		args = append(args, e.StartDatetime)
	}

	if e.EndDatetime != "" {
		if _, err := time.Parse(time.DateTime, e.EndDatetime); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "400 - Invalid ending date and time",
			})
			return
		}

		updates = append(updates, "jadwalSelesai = ?")
		args = append(args, e.EndDatetime)
	}

	updateQuestions := false
	updateStudents := false

	// TODO: add questions and students updates
	if len(e.Questions) > 0 {
		// get every existing questions on this exam
		q := "SELECT idSoal FROM soal_ujian WHERE idUjian = ?"
		rows, err := config.DB.Query(q, id)

		if err != nil {
			log.Printf("Get exam questions error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal server error",
			})
			return
		}

		defer rows.Close()

		// map for old question ids
		oldQuestionIDs := make(map[int]bool)

		for rows.Next() {
			var qid int

			if err := rows.Scan(&qid); err != nil {
				log.Printf("Get exam questions error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "500 - Internal server error",
				})
				return
			}

			oldQuestionIDs[qid] = true
		}

		// map for new question ids
		newQuestionIDs := make(map[int]bool)

		for _, qid := range e.Questions {
			newQuestionIDs[qid] = true
		}

		var toAdd []int
		var toDelete []int

		for _, qid := range e.Questions {
			if !oldQuestionIDs[qid] {
				toAdd = append(toAdd, qid)
			}
		}

		for qid := range oldQuestionIDs {
			if !newQuestionIDs[qid] {
				toDelete = append(toDelete, qid)
			}
		}

		if len(toAdd) > 0 {
			for _, qid := range toAdd {
				q2 := "INSERT INTO soal_ujian (idUjian, idSoal) VALUES (?, ?)"
				_, err := config.DB.Exec(q2, id, qid)

				if err != nil {
					log.Printf("Update exam question (create) error: %v", err)
					log.Println(id, qid)
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "500 - Internal server error",
					})
					return
				}
			}
		}

		if len(toDelete) > 0 {
			for _, qid := range toDelete {
				q3 := "DELETE FROM soal_ujian WHERE idUjian = ? AND idSoal = ?"
				_, err := config.DB.Exec(q3, id, qid)

				if err != nil {
					log.Printf("Update exam question (delete) error: %v", err)
					log.Println(id, qid)
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "500 - Internal server error",
					})
					return
				}
			}
		}

		updateQuestions = true
	}

	if len(e.Students) > 0 {
		q := "SELECT nim FROM ujian_ikut WHERE idUjian = ?"
		rows, err := config.DB.Query(q, id)

		if err != nil {
			log.Printf("Get exam students error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal server error",
			})
			return
		}

		defer rows.Close()

		// map for old student NIMs
		oldStudentNIMs := make(map[string]bool)

		for rows.Next() {
			var nim string

			if err := rows.Scan(&nim); err != nil {
				log.Printf("Get exam students error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "500 - Internal server error",
				})
				return
			}

			oldStudentNIMs[nim] = true
		}

		// map for new student NIMs
		newStudentNIMs := make(map[string]bool)

		for _, nim := range e.Students {
			newStudentNIMs[nim] = true
		}

		var toAdd []string
		var toDelete []string

		for _, nim := range e.Students {
			if !oldStudentNIMs[id] {
				toAdd = append(toAdd, nim)
			}
		}

		for nim := range oldStudentNIMs {
			if !newStudentNIMs[id] {
				toDelete = append(toDelete, nim)
			}
		}

		if len(toAdd) > 0 {
			for _, nim := range toAdd {
				q2 := "INSERT INTO ujian_ikut (nim, idUjian) VALUES (?, ?)"
				_, err := config.DB.Exec(q2, nim, id)

				if err != nil {
					log.Printf("Update exam student (create) error: %v", err)
					log.Println(id, nim)
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "500 - Internal server error",
					})
					return
				}
			}
		}

		if len(toDelete) > 0 {
			for _, nim := range toDelete {
				q3 := "DELETE FROM ujian_ikut WHERE idUjian = ? AND nim = ?"
				_, err := config.DB.Exec(q3, e.ExamID, nim)

				if err != nil {
					log.Printf("Update exam student (delete) error: %v", err)
					log.Println(e.ExamID, nim)
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "500 - Internal server error",
					})
					return
				}
			}
		}

		updateStudents = true
	}

	if len(updates) == 0 && !updateQuestions && !updateStudents {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - No fields to update",
		})
		return
	}

	updatedTime := time.Now().Format(time.DateTime)

	updates = append(updates, "dateUpdated = ?")
	args = append(args, updatedTime)

	args = append(args, id)

	query := "UPDATE ujian SET " + strings.Join(updates, ", ") + " WHERE idUjian = ?"
	_, err := config.DB.Exec(query, args...)

	if err != nil {
		log.Printf("Update exam error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - Exam updated successfully",
	})
}

func DeleteExam(c *gin.Context) {
	/*
		Delete an exam from the database.
	*/
	id := c.Param("id")

	query := "DELETE FROM ujian WHERE idUjian = ?"
	_, err := config.DB.Exec(query, id)

	if err != nil {
		log.Printf("Delete exam error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - Exam deleted successfully",
	})
}
