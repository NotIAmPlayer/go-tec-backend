package controllers

import (
	"database/sql"
	"go-tec-backend/config"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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
	accessTokenLifespan, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_HOUR_LIFESPAN"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	refreshTokenLifespan, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DAYS_LIFESPAN"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	domain := os.Getenv("FRONTEND_URL_JWT")

	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid request format",
		})
		return
	}

	accessToken, refreshToken, err := verifyLogin(request.Username, request.Password)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "400 - Invalid username or password",
		})
		return
	}

	var secure bool
	if gin.Mode() == gin.ReleaseMode {
		secure = true
	} else {
		secure = false
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"jwt-access",
		accessToken,
		accessTokenLifespan*int(time.Hour),
		"/",
		domain,
		secure,
		true,
	)

	c.SetCookie(
		"jwt-refresh",
		refreshToken,
		refreshTokenLifespan*int(time.Hour)*24,
		"/",
		domain,
		secure,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - Login successful",
		"user":    request.Username,
	})
}

func Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("jwt-refresh")

	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "401 - No refresh token",
		})
		return
	}

	tokenClaims, err := config.ParseToken(refreshToken, []byte(os.Getenv("REFRESH_TOKEN_SECRET")))

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "401 - Missing or invalid refresh token")
		return
	}

	userId := (*tokenClaims)["id"].(string)
	userEmail := (*tokenClaims)["email"].(string)
	userRole := (*tokenClaims)["role"].(string)

	newAccessToken, err := config.GenerateJWTAccessToken(userId, userEmail, userRole)

	// prepare assigning a new access token cookie
	accessTokenLifespan, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_HOUR_LIFESPAN"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Internal server error",
		})
		return
	}

	domain := os.Getenv("FRONTEND_URL_JWT")

	var secure bool
	if gin.Mode() == gin.ReleaseMode {
		secure = true
	} else {
		secure = false
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"jwt-access",
		newAccessToken,
		accessTokenLifespan*int(time.Hour),
		"/",
		domain,
		secure,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Access token refreshed",
	})
}

func Logout(c *gin.Context) {
	var secure bool
	if gin.Mode() == gin.ReleaseMode {
		secure = true
	} else {
		secure = false
	}

	c.SetCookie("jwt-access", "", -1, "/", "", secure, true)
	c.SetCookie("jwt-refresh", "", -1, "/", "", secure, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "200 - Logout successful",
	})
}

func verifyLogin(username, password string) (string, string, error) {
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
		return "", "", err
	}

	if !comparePasswordHash(password, user.Password) {
		log.Println("Password does not match")
		return "", "", bcrypt.ErrMismatchedHashAndPassword
	}

	accessToken, err := config.GenerateJWTAccessToken(user.ID, user.Email, user.Role)

	if err != nil {
		return "", "", err
	}

	refreshToken, err := config.GenerateJWTRefreshToken(user.ID, user.Email, user.Role)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func GetMe(c *gin.Context) {
	/*
		Get logged in user's credentials (ID, name, email, role) for frontend purposes
	*/
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "401 - Unauthorized"})
		return
	}

	email, exists := c.Get("email")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "401 - Unauthorized"})
		return
	}

	role, exists := c.Get("role")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "401 - Unauthorized"})
		return
	}

	var query string

	if role == "mahasiswa" {
		query = "SELECT nim AS id, namaMhs AS nama, email, 'mahasiswa' AS role FROM mahasiswa WHERE nim = ? OR email = ?"
	} else if role == "admin" {
		query = "SELECT idAdmin AS id, namaAdmin AS nama, email, 'admin' AS role FROM admin WHERE idAdmin = ? OR email = ?"
	} else {
		c.JSON(http.StatusForbidden, gin.H{"message": "403 - Invalid role"})
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
