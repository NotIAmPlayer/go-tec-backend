package controllers

import (
	"database/sql"
	"go-tec-backend/config"
	"log"
	"net/http"
	"strconv"
	"strings"

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

/*
func comparePasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
*/

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

	query := "SELECT nim, nama_mhs, email FROM mahasiswa LIMIT ? OFFSET ?"

	rows, err := config.DB.Query(query, limit, offset)

	if err != nil {
		log.Printf("Get multiple users error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

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

	u.Password = hash

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
	})
}

func UpdateUser(c *gin.Context) {
	/*
		Update a user by NIM in the database from JSON data.
		Does not update password.
	*/
	nim := c.Param("nim")

	var u Users

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid JSON data",
		})
		return
	}

	updates := []string{}
	args := []interface{}{}

	if u.Nama != "" {
		updates = append(updates, "nama_mhs = ?")
		args = append(args, u.Nama)
	}

	if u.Email != "" {
		updates = append(updates, "email = ?")
		args = append(args, u.Email)
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - No fields to update",
		})
		return
	}

	args = append(args, nim)
	query := "UPDATE mahasiswa SET " + strings.Join(updates, ", ") + " WHERE nim = ?"

	_, err := config.DB.Exec(query, args...)

	if err != nil {
		log.Printf("Update user error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - User updated successfully",
	})
}

func UpdateUserPassword(c *gin.Context) {
	/*
		Update a user's password by NIM in the database from JSON data.
	*/
	nim := c.Param("nim")

	var u Users

	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid JSON data",
		})
		return
	}

	if u.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - No password to update",
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

	u.Password = hash

	query := "UPDATE mahasiswa SET password = ? WHERE nim = ?"
	_, err = config.DB.Exec(query, u.Password, nim)

	if err != nil {
		log.Printf("Update user password error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - User password updated successfully",
	})
}

func DeleteUser(c *gin.Context) {
	/*
		Delete a user by NIM from the database.
	*/
	nim := c.Param("nim")

	query := "DELETE FROM mahasiswa WHERE nim = ?"
	_, err := config.DB.Exec(query, nim)

	if err != nil {
		log.Printf("Delete user error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - User deleted successfully",
	})
}
