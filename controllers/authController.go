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

	token, err := verifyLogin(request.Username, request.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "400 - Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func verifyLogin(username, password string) (string, error) {
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
		return "", err
	}

	if !comparePasswordHash(password, user.Password) {
		log.Println("Password does not match")
		return "", bcrypt.ErrMismatchedHashAndPassword
	}

	token, err := config.GenerateJWT(user.ID, user.Email, user.Role)

	if err != nil {
		return "", err
	}

	return token, nil
}

func GetMe(c *gin.Context) {
	/*
		Get logged in user's credentials (ID, name, email, role) for frontend purposes
	*/
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "401 - Unauthorized"})
		return
	}

	email, exists := c.Get("email")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "401 - Unauthorized"})
		return
	}

	role, exists := c.Get("role")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "401 - Unauthorized"})
		return
	}

	var query string

	if role == "mahasiswa" {
		query = "SELECT nim AS id, namaMhs AS nama, email, 'mahasiswa' AS role FROM mahasiswa WHERE nim = ? OR email = ?"
	} else if role == "admin" {
		query = "SELECT idAdmin AS id, namaAdmin AS nama, email, 'admin' AS role FROM admin WHERE idAdmin = ? OR email = ?"
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "403 - Invalid role"})
		return
	}

	row := config.DB.QueryRow(query, userID, email)

	var u UserData

	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Role); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "404 - User data not found",
			})
		} else {
			log.Printf("Get one user data error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    u.ID,
		"name":  u.Name,
		"email": u.Email,
		"role":  u.Role,
	})
}
