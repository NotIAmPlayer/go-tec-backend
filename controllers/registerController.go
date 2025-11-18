package controllers

import (
	"database/sql"
	"fmt"
	"go-tec-backend/config"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var u User

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(400, gin.H{
			"message": "400 - Bad request",
		})
		return
	}

	DumpContext(c)
	fmt.Println("Registering user:", u)

	if u.Nim == "" || u.Nama == "" || u.Password == "" || u.Email == "" {
		c.JSON(400, gin.H{
			"message": "400 - Missing required fields",
		})
		return
	}

	// validate email format
	// validate email format
	emailSubstring := strings.Split(u.Email, "@")

	if len(emailSubstring) != 2 || emailSubstring[0] == "" || emailSubstring[1] == "" {
		c.JSON(400, gin.H{
			"message": "400 - Invalid email format",
		})
		return
	}

	// valid domains
	domain := emailSubstring[1]
	if domain != "ukdc.ac.id" && domain != "student.ukdc.ac.id" {
		c.JSON(400, gin.H{
			"message": "400 - Email must be from ukdc.ac.id or student.ukdc.ac.id domain",
		})
		return
	}


	// check if user already exists (via NIM or email)
	query := "SELECT nim, namaMhs, email, password FROM mahasiswa WHERE nim = ? OR email = ?"

	var existingUser User

	row := config.DB.QueryRow(query, u.Nim, u.Email)

	if err := row.Scan(&existingUser.Nim, &existingUser.Nama, &existingUser.Email, &existingUser.Password); err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Get existing users error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal server error",
			})
			return
		}
	}

	if existingUser.Nim == u.Nim {
		c.JSON(http.StatusConflict, gin.H{
			"message": "409 - User with this NIM already exists",
		})
		return
	}

	if existingUser.Email == u.Email {
		c.JSON(http.StatusConflict, gin.H{
			"message": "409 - User with this email already exists",
		})
		return
	}

	hash, err := hashPassword(u.Password)

	if err != nil {
		log.Printf("Hashing new user password error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	query2 := "INSERT INTO mahasiswa (nim, namaMhs, email, password) VALUES (?, ?, ?, ?)"

	_, err = config.DB.Exec(query2, u.Nim, u.Nama, u.Email, hash)

	if err != nil {
		log.Printf("Create user error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "201 - User created successfully",
	})
}
