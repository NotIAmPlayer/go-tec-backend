package middlewares

import "github.com/gin-gonic/gin"

func HandleCORS(frontend string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if gin.Mode() == gin.ReleaseMode {
			c.Writer.Header().Set("Access-Control-Allow-Origin", frontend)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins in development mode
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // No Content
			return
		}
		c.Next()
	}
}
