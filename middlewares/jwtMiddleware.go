package middlewares

import (
	"go-tec-backend/config"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cookieToken, err := c.Cookie("jwt-access")

		if err != nil || cookieToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "401 - Unauthorized")
			return
		}

		tokenClaims, err := config.ParseToken(cookieToken, []byte(os.Getenv("ACCESS_TOKEN_SECRET")))

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "401 - Missing or invalid access token")
			return
		}

		c.Set("user_id", (*tokenClaims)["id"])
		c.Set("email", (*tokenClaims)["email"])
		c.Set("role", (*tokenClaims)["role"])
		c.Next()
	}
}
