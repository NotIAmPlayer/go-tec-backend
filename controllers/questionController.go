package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	QuestionJSON string     `form:"questions"`
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
	var b QuestionBatch

	DumpContext(c)

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

	/*
		if b.BatchType == "listening" {

		}
	*/

	err := json.Unmarshal([]byte(b.QuestionJSON), &b.Questions)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid questions JSON string",
		})
		return
	}

	fmt.Println(b.BatchID, b.BatchType, b.BatchText)

	for _, q := range b.Questions {
		fmt.Println(q)
	}
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

			// Limit file size to 10MB
			if file.Size > 10*1024*1024 {
				log.Printf("File size error: %s exceeds 10MB limit", file.Filename)
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "400 - File size exceeds 10MB limit",
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

			if file.Size > 10*1024*1024 {
				log.Printf("File size error: %s exceeds 10MB limit", file.Filename)
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "400 - File size exceeds 10MB limit",
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
