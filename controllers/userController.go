package controllers

import (
	"database/sql"
	"fmt"
	"go-tec-backend/config"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	Nim      string `json:"nim"`
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(bytes), err
}

func comparePasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func GetUsers(c *gin.Context) {
	/*
		Get users on a specific page from the database as JSON.
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

	users := []Users{}

	var query = "SELECT nim, nama_mhs, email FROM mahasiswa LIMIT ? OFFSET ?"

	rows, err := config.DB.Query(query, limit, offset)

	if err != nil {
		log.Printf("Get multiple users error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	fmt.Println("Query:", query, "Limit:", limit, "Offset:", offset, "Row Length:")

	defer rows.Close()

	for rows.Next() {
		var u Users

		if err := rows.Scan(&u.Nim, &u.Nama, &u.Email); err != nil {
			log.Printf("Get multiple users error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal server error",
			})
			return
		}

		users = append(users, u)
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "404 - No users found on page " + strconv.Itoa(page),
		})
		return
	} else {
		c.JSON(http.StatusOK, users)
	}
}

func GetUser(c *gin.Context) {
	/*
		Get a user by NIM from the database as JSON.
	*/
	nim := c.Param("nim")

	var u Users

	query := "SELECT nim, nama_mhs, email FROM mahasiswa WHERE nim = ?"

	row := config.DB.QueryRow(query, nim)

	if err := row.Scan(&u.Nim, &u.Nama, &u.Email); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "404 - User not found",
			})
		} else {
			log.Printf("Get one user error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, u)
}

func CreateUser(c *gin.Context) {
	/*
		Create a new user in the database from JSON data.
	*/

	var u Users

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid JSON data",
		})
		return
	}

	hash, err := hashPassword(u.Password)

	if err != nil {
		log.Printf("Hashing password error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	u.Password = string(hash)

	query := "INSERT INTO mahasiswa (nim, nama_mhs, email, password) VALUES (?, ?, ?, ?)"
	_, err = config.DB.Exec(query, u.Nim, u.Nama, u.Email, u.Password)

	if err != nil {
		log.Printf("Create user error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "201 - User created successfully",
		"user": gin.H{
			"nim":   u.Nim,
			"nama":  u.Nama,
			"email": u.Email,
		},
	})
}
