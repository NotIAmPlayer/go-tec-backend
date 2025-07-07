package controllers

import (
	"database/sql"
	"go-tec-backend/config"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserData struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func comparePasswordHash(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "400 - Invalid request format"})
		return
	}

	token, err := verifyLogin(request.Username, request.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "400 - Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func verifyLogin(username, password string) (string, error) {
	query := "SELECT nim AS id, email, password FROM mahasiswa WHERE nim = ? OR email = ? UNION SELECT idAdmin AS id, email, password FROM admin WHERE idAdmin = ? OR email = ?"

	rows := config.DB.QueryRow(query, username, username, username, username)

	var user UserData

	if err := rows.Scan(&user.ID, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			log.Println("User not found or password incorrect")
		} else {
			log.Printf("Error scanning user data: %v", err)
		}
		return "", err
	}

	if !comparePasswordHash(password, user.Password) {
		log.Println("Password does not match")
		return "", bcrypt.ErrMismatchedHashAndPassword
	}

	token, err := config.GenerateJWT(user.ID, user.Email)

	if err != nil {
		return "", err
	}

	return token, nil
}
