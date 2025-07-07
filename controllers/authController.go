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
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
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

	token, user, err := verifyLogin(request.Username, request.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "400 - Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name, // For frontend display
			"email": user.Email,
			"role":  user.Role, // 'mahasiswa' or 'admin'
		},
	})
}

func verifyLogin(username, password string) (string, UserData, error) {
	query := `
		SELECT nim AS id, namaMhs AS nama, email, password, 'mahasiswa' AS role FROM mahasiswa WHERE nim = ? OR email = ?
		UNION
		SELECT idAdmin AS id, namaAdmin AS nama, email, password, 'admin' AS role FROM admin WHERE idAdmin = ? OR email = ?
	`

	rows := config.DB.QueryRow(query, username, username, username, username)

	var user UserData

	if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			log.Println("User not found or password incorrect")
		} else {
			log.Printf("Error scanning user data: %v", err)
		}
		return "", user, err
	}

	if !comparePasswordHash(password, user.Password) {
		log.Println("Password does not match")
		return "", user, bcrypt.ErrMismatchedHashAndPassword
	}

	token, err := config.GenerateJWT(user.ID, user.Email, user.Role)

	if err != nil {
		return "", user, err
	}

	return token, user, nil
}
