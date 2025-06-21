package controllers

import (
	"os"

	"github.com/gin-gonic/gin"
)

func GetAudioFile(c *gin.Context) {
	filename := c.Param("filename")
	fullPath := "uploads/" + filename

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}

	c.File(fullPath)
}
