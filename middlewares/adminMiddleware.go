package middlewares

import (
	"net/http"
	"fmt"
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
		fmt.Println("JWTAuthMiddleware running, token role =", role)
	}
}
