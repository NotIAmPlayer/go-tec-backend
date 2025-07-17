package middlewares

import (
	"go-tec-backend/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := config.ValidateJWT(c)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "401 - Missing or invalid token")
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("user_id", claims["id"])
			c.Set("email", claims["email"])
			c.Set("role", claims["role"])
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "401 - Invalid token")
			return
		}
	}
}
