package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func DumpContext(c *gin.Context) {
	fmt.Printf("Method: %s\n", c.Request.Method)
	fmt.Printf("Request Headers: %v\n", c.Request.Header)
	fmt.Printf("URL: %s\n", c.Request.URL.Path)
	fmt.Printf("User-Agent Header: %s\n", c.GetHeader("User-Agent"))
	fmt.Printf("Client IP: %s\n", c.ClientIP())
}
