package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTAccessToken(id, email, role string) (string, error) {
	accessTokenLifespan, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_HOUR_LIFESPAN"))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Duration(accessTokenLifespan) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("ACCESS_TOKEN_SECRET")))
}

func GenerateJWTRefreshToken(id, email, role string) (string, error) {
	refreshTokenLifespan, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_DAYS_LIFESPAN"))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"id":    id,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(time.Duration(refreshTokenLifespan) * time.Hour * 24),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(os.Getenv("REFRESH_TOKEN_SECRET")))
}

func ParseToken(tokenString string, secret []byte) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	return &claims, nil
}

var ErrMissingAccessToken = errors.New("missing jwt access token cookie")
