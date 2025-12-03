package controllers

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"go-tec-backend/config"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordInput struct {
	NIM         string `json:"nim" binding:"required"`
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

type PasswordResetToken struct {
	ResetID   int       `json:"reset_id"`
	NIM       string    `json:"nim"`
	TokenHash string    `json:"token_hash"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
}

func ForgotPassword(c *gin.Context) {
	var input ForgotPasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	var user UserData

	query := "SELECT nim AS id, namaMhs AS nama, email FROM mahasiswa WHERE email = ?"

	row := config.DB.QueryRow(query, input.Email)

	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "A reset link was sent."})
		return
	}

	b := make([]byte, 32)
	rand.Read(b)
	token := hex.EncodeToString(b)

	tokenHash := sha256.Sum256([]byte(token))
	stringTokenHash := hex.EncodeToString(tokenHash[:])

	query2 := "INSERT INTO password_reset_tokens (nim, token_hash, expires_at, used) VALUES (?, ?, ?, ?)"

	expiryMinutes, err := strconv.Atoi(os.Getenv("PASSWORD_RESET_EXPIRY_MINUTES"))

	if err != nil {
		log.Fatal("PASSWORD_RESET_EXPIRY_MINUTES environment variable is not set as a valid integer")
		return
	}

	_, err = config.DB.Exec(
		query2,
		user.ID,
		stringTokenHash,
		time.Now().Add(time.Duration(expiryMinutes)*time.Minute).Format(time.DateTime),
		false,
	)

	if err != nil {
		log.Println("Insert forgot password token error:", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create password reset token",
		})
		return
	}

	// Godotenv has been called at this point in time
	frontend, exists := os.LookupEnv("FRONTEND_URL")

	if !exists {
		log.Fatal("FRONTEND_URL environment variable is not set")
		return
	}

	resetURL := fmt.Sprintf(
		"%s/reset-password?token=%s&nim=%s",
		frontend, token, user.ID,
	)

	err = config.SendResetEmail(user.Email, resetURL, expiryMinutes)

	if err != nil {
		log.Println("Send email failed:", err)

		c.JSON(http.StatusOK, gin.H{"message": "A reset link was sent."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "A reset link was sent."})
}

func ResetPassword(c *gin.Context) {
	var input ResetPasswordInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	tokenHash := sha256.Sum256([]byte(input.Token))

	var reset PasswordResetToken

	// make sure used is false and check expires_at via SQL for easier checking
	query := `
		SELECT reset_id, nim, token_hash, expires_at, used
		FROM password_reset_tokens
		WHERE nim = ? AND token_hash = ? AND used = false AND expires_at > NOW()
	`

	row := config.DB.QueryRow(query, input.NIM, hex.EncodeToString(tokenHash[:]))

	if err := row.Scan(&reset.ResetID, &reset.NIM, &reset.TokenHash, &reset.ExpiresAt, &reset.Used); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "400 - Invalid or expired reset token",
			})
		} else {
			log.Printf("Find reset token error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal server error",
			})
		}

		return
	}

	if reset.Used || reset.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid or expired reset token",
		})
	}

	var u User

	query2 := "SELECT nim, namaMhs, email FROM mahasiswa WHERE nim = ?"

	row2 := config.DB.QueryRow(query2, input.NIM)

	if err := row2.Scan(&u.Nim, &u.Nama, &u.Email); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "404 - User not found",
			})
		} else {
			log.Printf("Find reset token error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Internal server error",
			})
		}
	}

	hash, err := hashPassword(input.NewPassword)

	if err != nil {
		log.Printf("Hashing password error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	query3 := "UPDATE mahasiswa SET password = ? WHERE nim = ?"

	_, err = config.DB.Exec(query3, hash, input.NIM)

	if err != nil {
		log.Printf("Reset password error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	query4 := "UPDATE password_reset_tokens SET used = ? WHERE reset_id = ?"

	_, err = config.DB.Exec(query4, true, reset.ResetID)

	if err != nil {
		log.Printf("Update reset password token error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - User password updated successfully",
	})
}
