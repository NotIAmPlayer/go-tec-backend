package controllers

import (
	"go-tec-backend/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Users struct {
	Nim      string `json:"nim"`
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func GetUsers(c *gin.Context) {
	users := []Users{}

	var query = "SELECT nim, nama_mhs, email, password FROM mahasiswa"

	rows, err := config.DB.Query(query)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "500 - Error fetching users",
		})
		return
	}

	defer rows.Close()

	for rows.Next() {
		var u Users

		if err := rows.Scan(&u.Nim, &u.Nama, &u.Email, &u.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "500 - Error scanning user",
			})
			return
		}

		users = append(users, u)
	}

	if len(users) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "404 - No users found",
		})
		return
	} else {
		c.JSON(http.StatusOK, users)
	}
}
