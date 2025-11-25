package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")

		if !exists || (role != "admin") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "401 - Unauthorized"})
			return
		}

		c.Next()
	}
}
