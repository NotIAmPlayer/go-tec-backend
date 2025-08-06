package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetAudioFile(c *gin.Context) {
	filename := c.Param("filename")
	fullPath := "uploads/" + filename

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 - File not found"})
		return
	}

	c.File(fullPath)
}
