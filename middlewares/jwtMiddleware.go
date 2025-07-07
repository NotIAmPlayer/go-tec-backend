package middlewares

import (
	"go-tec-backend/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := config.ValidateJWT(c)

		if err != nil {
			c.String(http.StatusUnauthorized, "401 - Unauthorized")
			c.Abort()
			return
		}

		c.Next()
	}
}
