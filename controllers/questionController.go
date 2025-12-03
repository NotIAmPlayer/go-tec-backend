package controllers

import (
	"database/sql"
	"encoding/json"
	"go-tec-backend/config"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Question struct {
	QuestionID   int    `json:"question_id"`
	QuestionText string `json:"question_text"`
	ChoiceA      string `json:"choice_a"`
	ChoiceB      string `json:"choice_b"`
	ChoiceC      string `json:"choice_c"`
	ChoiceD      string `json:"choice_d"`
	Answer       string `json:"answer"`
}

type QuestionBatch struct {
	BatchID      int        `json:"batch_id" form:"batch_id"`
	BatchType    string     `json:"batch_type" form:"batch_type"`
	BatchText    string     `json:"batch_text" form:"batch_text"`
	QuestionJSON string     `json:"question_json" form:"questions"`
	Questions    []Question `json:"questions"`
	AudioPath    string     `json:"audio_path" form:"audio_path"`
}

func GetAllQuestions(c *gin.Context) {
	/*
		Get all questions from the database as JSON.
	*/

	/*
		questions := []Question{}

		query := "SELECT idSoal, tipeSoal, isiSoal, pilihanA, pilihanB, pilihanC, pilihanD, jawaban, audio FROM soal ORDER BY idSoal ASC"

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
			var q Question

			var audio sql.NullString
			if err := rows.Scan(&q.QuestionID, &q.QuestionType, &q.QuestionText, &q.ChoiceA, &q.ChoiceB, &q.ChoiceC, &q.ChoiceD, &q.Answer, &audio); err != nil {
				log.Printf("Get multiple questions error: %v", err)
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

func GetAllQuestionBatches(c *gin.Context) {
	/*
		Get all question batches in the database as JSON.
	*/

	questionBatches := []QuestionBatch{}

	query := "SELECT idBatch, tipeBatch, textBatch, audio FROM batch_soal ORDER BY idBatch ASC"
	rows, err := config.DB.Query(query)

	if err != nil {
		log.Printf("Get multiple question batches error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal Server Error",
		})
		return
	}

	defer rows.Close()

	batchFailed := false

	for rows.Next() {
		var b QuestionBatch
		var audio sql.NullString

		// question batch
		if err := rows.Scan(&b.BatchID, &b.BatchType, &b.BatchText, &audio); err != nil {
			log.Printf("Get multiple question batches error: %v", err)
			batchFailed = true
			break
		}

		if audio.Valid {
			b.AudioPath = audio.String
		} else {
			b.AudioPath = ""
		}

		// all questions in a batch
		query2 := "SELECT idSoal, isiSoal, pilihanA, pilihanB, pilihanC, pilihanD, jawaban FROM soal WHERE idBatch = ? ORDER BY idSoal ASC"
		rows2, err := config.DB.Query(query2, b.BatchID)

		if err != nil {
			log.Printf("Get multiple question batches error: %v", err)
			batchFailed = true
			break
		}

		questionFailed := false

		for rows2.Next() {
			var q Question

			if err := rows2.Scan(&q.QuestionID, &q.QuestionText, &q.ChoiceA, &q.ChoiceB, &q.ChoiceC, &q.ChoiceD, &q.Answer); err != nil {
				log.Printf("Get multiple question batches error: %v", err)
				questionFailed = true
				break
			}

			b.Questions = append(b.Questions, q)
		}

		if questionFailed {
			batchFailed = true
			rows2.Close()
			break
		}

		questionBatches = append(questionBatches, b)
		rows2.Close()
	}

	if batchFailed {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal Server Error",
		})
		return
	}

	if len(questionBatches) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "200 - No questions found",
		})
		return
	} else {
		c.JSON(http.StatusOK, questionBatches)
	}
}

func GetQuestion(c *gin.Context) {
	/*
		Get a question by ID from the database as JSON.
	*/

	/*
		id, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "400 - Invalid question ID",
			})
			return
		}

		var q Question

		query := "SELECT idSoal, tipeSoal, isiSoal, pilihanA, pilihanB, pilihanC, pilihanD, jawaban, audio FROM soal WHERE idSoal = ?"

		row := config.DB.QueryRow(query, id)

		var audio sql.NullString
		if err := row.Scan(&q.QuestionID, &q.QuestionType, &q.QuestionText, &q.ChoiceA, &q.ChoiceB, &q.ChoiceC, &q.ChoiceD, &q.Answer, &audio); err != nil {
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
		if audio.Valid {
			q.AudioPath = audio.String
		} else {
			q.AudioPath = ""
		}

		c.JSON(http.StatusOK, q)
	*/

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - work in progress",
	})
}

func CreateQuestionBatch(c *gin.Context) {
	/*
		Create question batches in the database from a form data.
	*/
	var b QuestionBatch

	if err := c.ShouldBind(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid form data",
		})
		return
	}

	if b.QuestionJSON == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Questions cannot be empty",
		})
		return
	}

	if b.BatchType != "listening" && b.BatchType != "reading" && b.BatchType != "grammar" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid question type",
		})
		return
	}

	if b.BatchType == "reading" && b.BatchText == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Reading question batches must have a reading text.",
		})
		return
	}

	// turn the json string for array of questions back to array of questions
	err := json.Unmarshal([]byte(b.QuestionJSON), &b.Questions)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid questions JSON string",
		})
		return
	}

	// debug prints for each questions
	/*
		fmt.Println(b.BatchID, b.BatchType, b.BatchText)

		for _, q := range b.Questions {
			fmt.Println(q)
		}
	*/

	var filename string

	if b.BatchType == "listening" {
		file, err := c.FormFile("file")

		if err != nil && err != http.ErrMissingFile {
			log.Printf("File upload error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "400 - Invalid file upload",
			})
			return
		}

		if file == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "400 - File is required",
			})
			return
		}

		// If uploads directory does not exist, create it
		if _, err := os.Stat("uploads"); os.IsNotExist(err) {
			if err := os.Mkdir("uploads", 0755); err != nil {
				log.Printf("Directory creation error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "500 - Internal server error",
				})
				return
			}
			log.Println("Uploads directory created")
		}

		// Save file to the uploads directory
		dest := "uploads/" + file.Filename

		// Limit file size to 50MB
		if file.Size > 50*1024*1024 {
			log.Printf("File size error: %s exceeds 50MB limit", file.Filename)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "400 - File size exceeds 50MB limit",
			})
			return
		}

		if err := c.SaveUploadedFile(file, dest); err != nil {
			log.Printf("File save error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal server error",
			})
			return
		}

		log.Printf("File uploaded successfully: %s", dest)
		filename = file.Filename
	}

	createdTime := time.Now().Format(time.DateTime)

	// use a transaction for multi-table creation
	tx, err := config.DB.Begin()

	if err != nil {
		log.Printf("Begin transaction error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	commited := false

	defer func() {
		if !commited {
			tx.Rollback()
		}
	}()

	query := "INSERT INTO batch_soal (tipeBatch, textBatch, audio, dateCreated, dateUpdated) VALUES (?, ?, ?, ?, ?)"
	_, err = tx.Exec(query, b.BatchType, b.BatchText, filename, createdTime, createdTime)

	if err != nil {
		log.Printf("Create question batch error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	// fetch the newly created question batch's id
	query2 := "SELECT idBatch FROM batch_soal WHERE dateCreated = ?"
	row := tx.QueryRow(query2, createdTime)

	if err := row.Scan(&b.BatchID); err != nil {
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

	questionsFailed := false

	// time to actually upload each questions
	for _, q := range b.Questions {
		query3 := "INSERT INTO soal (idBatch, isiSoal, pilihanA, pilihanB, pilihanC, pilihanD, jawaban) VALUES (?, ?, ?, ?, ?, ?, ?)"
		_, err := tx.Exec(query3, b.BatchID, q.QuestionText, q.ChoiceA, q.ChoiceB, q.ChoiceC, q.ChoiceD, q.Answer)

		if err != nil {
			log.Printf("Create question error: %v", err)
			questionsFailed = true
			break
		}
	}

	if questionsFailed {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Commit transaction error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	commited = true
	c.JSON(http.StatusCreated, gin.H{
		"message": "201 - Question batches created successfully",
	})
}

func CreateQuestion(c *gin.Context) {
	/*
		Create a new question in the database from JSON data.
	*/

	/*
		var q Question

		if err := c.ShouldBind(&q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "400 - Invalid form data",
			})
			return
		}

		if q.QuestionType != "listening" && q.QuestionType != "reading" && q.QuestionType != "grammar" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "400 - Invalid question type",
			})
			return
		}

		if q.QuestionText == "" || q.ChoiceA == "" || q.ChoiceB == "" || q.ChoiceC == "" || q.ChoiceD == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "400 - Question text and answer texts are required",
			})
			return
		}

		if q.Answer != "a" && q.Answer != "b" && q.Answer != "c" && q.Answer != "d" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "400 - Invalid answer letter",
			})
			return
		}

		var filename string

		if q.QuestionType == "listening" {
			file, err := c.FormFile("file")

			if err != nil && err != http.ErrMissingFile {
				log.Printf("File upload error: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "400 - Invalid file upload",
				})
				return
			}

			if file == nil {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "400 - File is required",
				})
				return
			}

			// If uploads directory does not exist, create it
			if _, err := os.Stat("uploads"); os.IsNotExist(err) {
				if err := os.Mkdir("uploads", 0755); err != nil {
					log.Printf("Directory creation error: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "500 - Internal server error",
					})
					return
				}
				log.Println("Uploads directory created")
			}

			// Save file to the uploads directory
			dest := "uploads/" + file.Filename

			// Limit file size to 50MB
			if file.Size > 50*1024*1024 {
				log.Printf("File size error: %s exceeds 50MB limit", file.Filename)
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "400 - File size exceeds 50MB limit",
				})
				return
			}

			if err := c.SaveUploadedFile(file, dest); err != nil {
				log.Printf("File save error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "500 - Internal server error",
				})
				return
			}

			log.Printf("File uploaded successfully: %s", dest)
			filename = file.Filename
		}

		query := "INSERT INTO soal (tipeSoal, isiSoal, pilihanA, pilihanB, pilihanC, pilihanD, jawaban, audio) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
		_, err := config.DB.Exec(query, q.QuestionType, q.QuestionText, q.ChoiceA, q.ChoiceB, q.ChoiceC, q.ChoiceD, q.Answer, filename)

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
	*/

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - work in progress",
	})
}

func UpdateQuestionBatch(c *gin.Context) {
	/*
		Update question batches in the database from a form data.
	*/
	id := c.Param("id")

	var b QuestionBatch

	if err := c.ShouldBind(&b); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid form data",
		})
		return
	}

	if b.QuestionJSON == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Questions cannot be empty",
		})
		return
	}

	if b.BatchType != "listening" && b.BatchType != "reading" && b.BatchType != "grammar" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid question type",
		})
		return
	}

	if b.BatchType == "reading" && b.BatchText == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Reading question batches must have a reading text.",
		})
		return
	}

	// turn the json string for array of questions back to array of questions
	err := json.Unmarshal([]byte(b.QuestionJSON), &b.Questions)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid questions JSON string",
		})
		return
	}

	file, err := c.FormFile("file")

	if err != nil && err != http.ErrMissingFile {
		log.Printf("File upload error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid file upload",
		})
		return
	}

	var filename string

	if file != nil {
		// If uploads directory does not exist, create it
		if _, err := os.Stat("uploads"); os.IsNotExist(err) {
			if err := os.Mkdir("uploads", 0755); err != nil {
				log.Printf("Directory creation error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "500 - Internal server error",
				})
				return
			}
			log.Println("Uploads directory created")
		}

		// Save file to the uploads directory
		dest := "uploads/" + file.Filename

		if file.Size > 50*1024*1024 {
			log.Printf("File size error: %s exceeds 50MB limit", file.Filename)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "400 - File size exceeds 50MB limit",
			})
			return
		}

		if err := c.SaveUploadedFile(file, dest); err != nil {
			log.Printf("File save error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal server error",
			})
			return
		}

		log.Printf("File uploaded successfully: %s", dest)
		filename = file.Filename
	}

	// use a transaction for multi-table creation
	tx, err := config.DB.Begin()

	if err != nil {
		log.Printf("Begin transaction error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	commited := false

	defer func() {
		if !commited {
			tx.Rollback()
		}
	}()

	updates := []string{}
	args := []interface{}{}

	if file != nil {
		updates = append(updates, "audio = ?")
		args = append(args, filename)
	}

	if b.BatchType != "" && (b.BatchType == "listening" || b.BatchType == "reading" || b.BatchType == "grammar") {
		updates = append(updates, "tipeBatch = ?")
		args = append(args, b.BatchType)
	}

	if b.BatchText != "" {
		updates = append(updates, "textBatch = ?")
		args = append(args, b.BatchText)
	}

	updateQuestions := false
	questionsFailed := false

	if len(b.Questions) > 0 {
		query := "SELECT idSoal FROM soal WHERE idBatch = ?"
		rows, err := tx.Query(query, id)

		if err != nil {
			log.Printf("Get batch questions error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal server error",
			})
			return
		}

		defer rows.Close()

		oldQuestionIDs := make(map[int]bool)

		for rows.Next() {
			var qid int

			if err := rows.Scan(&qid); err != nil {
				log.Printf("Get batch questions error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "500 - Internal server error",
				})
				return
			}

			oldQuestionIDs[qid] = true
		}

		newQuestionIDs := make(map[int]bool)

		for _, q := range b.Questions {
			if q.QuestionID != 0 {
				newQuestionIDs[q.QuestionID] = true
			}
		}

		var toAdd []Question
		var toUpdate []Question
		var toDelete []int

		for _, q := range b.Questions {
			if q.QuestionID == 0 {
				toAdd = append(toAdd, q)
			}

			if oldQuestionIDs[q.QuestionID] {
				toUpdate = append(toUpdate, q)
			}
		}

		for qid := range oldQuestionIDs {
			if !newQuestionIDs[qid] {
				toDelete = append(toDelete, qid)
			}
		}

		if len(toAdd) > 0 {
			for _, q := range toAdd {
				query2 := "INSERT INTO soal (idBatch, isiSoal, pilihanA, pilihanB, pilihanC, pilihanD, jawaban) VALUES (?, ?, ?, ?, ?, ?, ?)"
				_, err := tx.Exec(query2, id, q.QuestionText, q.ChoiceA, q.ChoiceB, q.ChoiceC, q.ChoiceD, q.Answer)

				if err != nil {
					log.Printf("Create question (update) error: %v", err)
					questionsFailed = true
					break
				}
			}
		}

		if len(toUpdate) > 0 && !questionsFailed {
			for _, q := range toUpdate {
				query2 := "UPDATE soal SET isiSoal = ?, pilihanA = ?, pilihanB = ?, pilihanC = ?, pilihanD = ?, jawaban = ? WHERE idSoal = ? AND idBatch = ?"
				_, err := tx.Exec(query2, q.QuestionText, q.ChoiceA, q.ChoiceB, q.ChoiceC, q.ChoiceD, q.Answer, q.QuestionID, id)

				if err != nil {
					log.Printf("Update question (update) error: %v", err)
					questionsFailed = true
					break
				}
			}
		}

		if len(toDelete) > 0 && !questionsFailed {
			for _, qid := range toDelete {
				query2 := "DELETE FROM soal WHERE idSoal = ? AND idBatch = ?"
				_, err := tx.Exec(query2, qid, id)

				if err != nil {
					log.Printf("Delete question (update) error: %v", err)
					questionsFailed = true
					break
				}
			}
		}

		if questionsFailed {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal server error",
			})
			return
		}

		updateQuestions = true
	}

	if len(updates) == 0 && !updateQuestions {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - No fields to update",
		})
		return
	}

	updatedTime := time.Now().Format(time.DateTime)

	updates = append(updates, "dateUpdated = ?")
	args = append(args, updatedTime)

	args = append(args, id)

	query := "UPDATE batch_soal SET " + strings.Join(updates, ", ") + " WHERE idBatch = ?"
	_, err = tx.Exec(query, args...)

	if err != nil {
		log.Printf("Update question batch error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Commit transaction error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	commited = true

	c.JSON(http.StatusOK, gin.H{
		"message": "Question batch updated successfully.",
	})
}

func UpdateQuestion(c *gin.Context) {
	/*
		Update a question in the database from JSON data.
	*/

	/*
		id := c.Param("id")

		file, err := c.FormFile("file")

		if err != nil && err != http.ErrMissingFile {
			log.Printf("File upload error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "400 - Invalid file upload",
			})
			return
		}

		var filename string

		if file != nil {
			// If uploads directory does not exist, create it
			if _, err := os.Stat("uploads"); os.IsNotExist(err) {
				if err := os.Mkdir("uploads", 0755); err != nil {
					log.Printf("Directory creation error: %v", err)
					c.JSON(http.StatusInternalServerError, gin.H{
						"message": "500 - Internal server error",
					})
					return
				}
				log.Println("Uploads directory created")
			}

			// Save file to the uploads directory
			dest := "uploads/" + file.Filename

			if file.Size > 50*1024*1024 {
				log.Printf("File size error: %s exceeds 50MB limit", file.Filename)
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "400 - File size exceeds 50MB limit",
				})
				return
			}

			if err := c.SaveUploadedFile(file, dest); err != nil {
				log.Printf("File save error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "500 - Internal server error",
				})
				return
			}

			log.Printf("File uploaded successfully: %s", dest)
			filename = file.Filename
		}

		var q Question

		if err := c.ShouldBind(&q); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "400 - Invalid form data",
			})
			return
		}

		updates := []string{}
		args := []interface{}{}

		if file != nil {
			updates = append(updates, "audio = ?")
			args = append(args, filename)
		}

		if q.QuestionType != "" && (q.QuestionType == "Listening" || q.QuestionType == "Reading" || q.QuestionType == "Grammar") {
			updates = append(updates, "tipeSoal = ?")
			args = append(args, q.QuestionType)
		}

		if q.QuestionText != "" {
			updates = append(updates, "isiSoal = ?")
			args = append(args, q.QuestionText)
		}

		if q.ChoiceA != "" {
			updates = append(updates, "pilihanA = ?")
			args = append(args, q.ChoiceA)
		}

		if q.ChoiceB != "" {
			updates = append(updates, "pilihanB = ?")
			args = append(args, q.ChoiceB)
		}

		if q.ChoiceC != "" {
			updates = append(updates, "pilihanC = ?")
			args = append(args, q.ChoiceC)
		}

		if q.ChoiceD != "" {
			updates = append(updates, "pilihanD = ?")
			args = append(args, q.ChoiceD)
		}

		if q.Answer != "" {
			updates = append(updates, "jawaban = ?")
			args = append(args, q.Answer)
		}

		if len(updates) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "400 - No fields to update",
			})
			return
		}

		args = append(args, id)
		query := "UPDATE soal SET " + strings.Join(updates, ", ") + " WHERE idSoal = ?"

		_, err = config.DB.Exec(query, args...)

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

	*/

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - work in progress",
	})
}

func DeleteQuestionBatch(c *gin.Context) {
	id := c.Param("id")

	query := "SELECT COUNT(*) batchUsedCount FROM batch_ujian WHERE idBatch = ?"
	row := config.DB.QueryRow(query, id)

	var batchUsedCount int

	if err := row.Scan(&batchUsedCount); err != nil {
		log.Printf("Get batch details error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal Server Error",
		})
		return
	}

	if batchUsedCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - This batch is used in at least one exam.",
		})
		return
	}

	query2 := "DELETE FROM soal WHERE idBatch = ?"
	_, err := config.DB.Exec(query2, id)

	if err != nil {
		log.Printf("Delete exam error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	query3 := "DELETE FROM batch_soal WHERE idBatch = ?"
	_, err = config.DB.Exec(query3, id)

	if err != nil {
		log.Printf("Delete exam error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - Question batch deleted successfully",
	})
}

func DeleteQuestion(c *gin.Context) {
	/*
		Delete a question from the database.
	*/

	/*
		id := c.Param("id")

		query := "DELETE FROM soal WHERE idSoal = ?"
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

	*/

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - work in progress",
	})
}
