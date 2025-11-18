package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"strings"
	"go-tec-backend/config"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Student struct {
	NIM   string `json:"nim"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func GetExamStudents(c *gin.Context) {
	examID := c.Param("id")
	db := config.DB

	log.Println("🔍 GetExamStudents called with examID =", examID)

	var rows *sql.Rows
	var err error

	if strings.HasPrefix(examID, "OFF") {
		query := `
			SELECT m.nim, m.namaMhs AS name, m.email
			FROM exam_offline_students eos
			JOIN mahasiswa m ON m.nim = eos.student_nim
			WHERE eos.exam_id = ?`
		rows, err = db.Query(query, examID)
	} else {
		query := `
			SELECT m.nim, m.namaMhs AS name, m.email
			FROM ujian_ikut ui
			JOIN mahasiswa m ON m.nim = ui.nim
			WHERE ui.idUjian = ?`
		rows, err = db.Query(query, examID)
	}

	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}); return }
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.NIM, &s.Name, &s.Email); err != nil { continue }
		students = append(students, s)
	}
	c.JSON(http.StatusOK, students)
}
