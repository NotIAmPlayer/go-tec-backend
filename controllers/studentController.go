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

type StudentQuestion struct {
	QuestionID   int    `json:"question_id"`
	QuestionText string `json:"question_text"`
	ChoiceA      string `json:"choice_a"`
	ChoiceB      string `json:"choice_b"`
	ChoiceC      string `json:"choice_c"`
	ChoiceD      string `json:"choice_d"`
	Answer       string `json:"answer"`
	AudioPath    string `json:"audio_path"`
	BatchID      int    `json:"batch_id"`
	BatchType    string `json:"batch_type"`
	BatchText    string `json:"batch_text"`
}

type AnswerData struct {
	Nim        string `json:"nim"`
	ExamID     int    `json:"exam_id"`
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer"`
	TipeBatch  int    `json:"tipeBatch"`
}

func GetUpcomingExams(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "401 - Unauthorized"})
		return
	}

	query := `
		SELECT u.idUjian, u.namaUjian, u.jadwalMulai, u.jadwalSelesai
		FROM ujian u
		JOIN ujian_ikut i ON u.idUjian = i.idUjian
		WHERE 
			u.jadwalSelesai >= NOW()
			AND i.nim = ?
			AND (i.statusPengerjaan <> 'selesai' OR ISNULL(i.statusPengerjaan))
		ORDER BY u.jadwalMulai ASC
	`

	rows, err := config.DB.Query(query, userID)
	if err != nil {
		log.Printf("Get multiple exams (student) error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "500 - Internal Server Error"})
		return
	}
	defer rows.Close()

	upcomingExams := []StudentExams{}
	for rows.Next() {
		var e StudentExams
		if err := rows.Scan(&e.ExamID, &e.ExamTitle, &e.StartDatetime, &e.EndDatetime); err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		upcomingExams = append(upcomingExams, e)
	}

	if len(upcomingExams) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "200 - No exams found"})
		return
	}

	c.JSON(http.StatusOK, upcomingExams)
}

func GetOfflineExamsForStudent(c *gin.Context) {
	var userID string

	if val, exists := c.Get("user_id"); exists {
		userID = fmt.Sprintf("%v", val)
	} else if val, exists := c.Get("user"); exists {
		if claims, ok := val.(map[string]interface{}); ok {
			if id, ok := claims["id"].(string); ok {
				userID = id
			} else if nim, ok := claims["nim"].(string); ok {
				userID = nim
			} else if idNum, ok := claims["id"].(float64); ok {
				userID = fmt.Sprintf("%.0f", idNum)
			}
		}
	}

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	type OfflineExam struct {
		ExamID        string `json:"exam_id"`
		ExamTitle     string `json:"exam_title"`
		StartDatetime string `json:"start_datetime"`
		EndDatetime   string `json:"end_datetime"`
		RoomName      string `json:"room_name"`
	}

	query := `
		SELECT 
			eo.id AS exam_id,
			eo.exam_title,
			eo.start_datetime,
			eo.end_datetime,
			eo.room_name
		FROM exam_offline eo
		INNER JOIN exam_offline_students eos ON eos.exam_id = eo.id
		WHERE eos.student_nim = ?
		  AND eo.end_datetime >= NOW()  -- ✅ hanya tampilkan yang belum lewat
		ORDER BY eo.start_datetime ASC
	`

	rows, err := config.DB.Query(query, userID)
	if err != nil {
		log.Printf("Error fetching offline exams for student (%s): %v", userID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "500 - Internal Server Error"})
		return
	}
	defer rows.Close()

	exams := []OfflineExam{}
	for rows.Next() {
		var e OfflineExam
		if err := rows.Scan(&e.ExamID, &e.ExamTitle, &e.StartDatetime, &e.EndDatetime, &e.RoomName); err == nil {
			exams = append(exams, e)
		} else {
			log.Printf("Scan error: %v", err)
		}
	}

	if len(exams) == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "200 - No offline exams found"})
		return
	}

	c.JSON(http.StatusOK, exams)
}

func GetAvailableOnlineExams(c *gin.Context) {
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "401 - Unauthorized"})
		return
	}

	query := `
		SELECT u.idUjian, u.namaUjian, u.jadwalMulai, u.jadwalSelesai
		FROM ujian u
		WHERE
			u.jadwalSelesai >= NOW() 		-- ujian yang belum lewat
			AND NOT EXISTS (				-- ujian yang belum mendaftar
				SELECT 1
				FROM pendaftaran_ujian p
				WHERE p.exam_id = u.idUjian
					AND p.nim = ?
			)
		ORDER BY u.jadwalMulai ASC
	`

	rows, err := config.DB.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Gagal ambil data ujian online",
			"error":   err.Error(),
		})
		return
	}
	defer rows.Close()

	type Exam struct {
		IDUjian       int    `json:"idUjian"`
		NamaUjian     string `json:"namaUjian"`
		JadwalMulai   string `json:"jadwalMulai"`
		JadwalSelesai string `json:"jadwalSelesai"`
	}

	var exams []Exam

	for rows.Next() {
		var (
			idUjian       int
			namaUjian     string
			jadwalMulai   time.Time
			jadwalSelesai time.Time
		)

		if err := rows.Scan(&idUjian, &namaUjian, &jadwalMulai, &jadwalSelesai); err == nil {
			exams = append(exams, Exam{
				IDUjian:       idUjian,
				NamaUjian:     namaUjian,
				JadwalMulai:   jadwalMulai.Format(time.DateTime),
				JadwalSelesai: jadwalSelesai.Format(time.DateTime),
			})
		}
	}

	c.JSON(http.StatusOK, exams)
}

func GetAvailableOfflineExams(c *gin.Context) {
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "401 - Unauthorized"})
		return
	}

	type AvailableExam struct {
		ExamID         string `json:"exam_id"`
		ExamTitle      string `json:"exam_title"`
		StartDatetime  string `json:"start_datetime"`
		EndDatetime    string `json:"end_datetime"`
		Room           string `json:"room"`
		QuotaTotal     int    `json:"quota_total"`
		QuotaAvailable int    `json:"quota_available"`
	}

	query := `
		SELECT 
			eo.id AS exam_id,
			eo.exam_title,
			eo.start_datetime,
			eo.end_datetime,
			COALESCE(eo.room_name, '-') AS room,
			IFNULL(k.total, 0) AS quota_total,
			IFNULL(k.available, 0) AS quota_available
		FROM exam_offline eo
		LEFT JOIN kuota_ujian k ON k.idUjian = eo.id
		WHERE 
			eo.end_datetime >= NOW()        -- ✅ hanya ujian yang belum lewat
			AND IFNULL(k.available, 0) > 0
			AND NOT EXISTS (				-- ujian yang belum mendaftar
				SELECT 1
				FROM pendaftaran_ujian p
				WHERE p.exam_id = eo.id
					AND p.nim = ?
			)
		ORDER BY eo.start_datetime ASC
	`

	rows, err := config.DB.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Gagal mengambil daftar ujian offline",
			"detail": err.Error(),
		})
		return
	}
	defer rows.Close()

	var exams []AvailableExam
	for rows.Next() {
		var e AvailableExam
		if err := rows.Scan(
			&e.ExamID,
			&e.ExamTitle,
			&e.StartDatetime,
			&e.EndDatetime,
			&e.Room,
			&e.QuotaTotal,
			&e.QuotaAvailable,
		); err == nil {
			exams = append(exams, e)
		}
	}

	c.JSON(http.StatusOK, exams)
}

func GetExamQuestions(c *gin.Context) {
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

	questions := []StudentQuestion{}

	query := `
		SELECT s.idSoal, b.idBatch, b.tipeBatch, b.textBatch, s.isiSoal, s.pilihanA, s.pilihanB, s.pilihanC, s.pilihanD, b.audio
		FROM batch_ujian u JOIN batch_soal b ON b.idBatch = u.idBatch JOIN soal s ON b.idBatch = s.idBatch
		WHERE u.idUjian = ?
		ORDER BY b.tipeBatch, s.idSoal
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
		var q StudentQuestion

		var audio sql.NullString
		if err := rows.Scan(&q.QuestionID, &q.BatchID, &q.BatchType, &q.BatchText, &q.QuestionText, &q.ChoiceA, &q.ChoiceB, &q.ChoiceC, &q.ChoiceD, &audio); err != nil {
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
			query2 = "INSERT INTO soal_jawaban (jawaban, nim, idSoal, idUjian, tipeBatch) VALUES (?, ?, ?, ?, ?)"
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
			query2 = "UPDATE soal_jawaban SET jawaban = ?, tipeBatch = ? WHERE nim = ? AND idSoal = ? AND idUjian = ?"
		}
	}

	if strings.HasPrefix(query2, "DELETE") {
		_, err = config.DB.Exec(query2, a.Nim, a.QuestionID, a.ExamID)
	} else if strings.HasPrefix(query2, "INSERT") {
		_, err = config.DB.Exec(query2, a.Answer, a.Nim, a.QuestionID, a.ExamID, a.TipeBatch)
	} else { // UPDATE
		_, err = config.DB.Exec(query2, a.Answer, a.TipeBatch, a.Nim, a.QuestionID, a.ExamID)
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

	query := "SELECT idSoal, jawaban, tipeBatch FROM soal_jawaban WHERE nim = ? AND idUjian = ?"
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

		if err := rows.Scan(&a.QuestionID, &a.Answer, &a.TipeBatch); err != nil {
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

	var target time.Time
	var err error

	// Coba parse dengan format baru (ISO/RFC3339)
	target, err = time.Parse(time.RFC3339, e.EndDatetime)
	if err != nil {
		// Kalau gagal (misal format lama dari DB), coba format lama
		target, err = time.Parse("2006-01-02 15:04:05", e.EndDatetime)
		if err != nil {
			log.Printf("Parse end datetime error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal Server Error",
			})
			return
		}
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
	_, err = config.DB.Exec(query2, currentDatetime.Format("2006-01-02 15:04:05"), "selesai", a.Nim, a.ExamID)

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

	loc, _ := time.LoadLocation("Asia/Jakarta")

	target, err := time.Parse(time.RFC3339, e.StartDatetime)
	if err != nil {
		// Kalau gagal (mungkin data lama dari DB), coba format MySQL
		target, err = time.ParseInLocation("2006-01-02 15:04:05", e.StartDatetime, loc)
		if err != nil {
			log.Printf("Parse start datetime error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Error parsing start time",
			})
			return
		}
	}

	currentDatetime := time.Now().In(loc)
	duration := currentDatetime.Sub(target)

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
