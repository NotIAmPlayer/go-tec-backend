package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		if gin.Mode() == gin.ReleaseMode {
			origin := r.Header.Get("Origin")

			frontend, exists := os.LookupEnv("FRONTEND_URL")

			if !exists {
				fmt.Println("FRONTEND_URL environment variable is not set")
				return false
			}

			if origin == frontend {
				return true
			}

			return false
		} else {
			return true
		}
	},
}

type WebsocketUser struct {
	UserID string
	Role   string
}

func ExamWebsocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	var students []WebsocketUser
	var admins []WebsocketUser

	if err != nil {
		fmt.Println("Error opening websocket:", err)
		return
	}

	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			fmt.Println("Websocket - Reading error:", err)
		}
		fmt.Println("Received:", message)

		var msg map[string]string

		err = json.Unmarshal(message, &msg)

		if err != nil {
			fmt.Println("Websocket - Error decoding JSON string:", err)
		}

		if msg["type"] == "auth" {
			token, err := jwt.Parse(msg["token"], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(os.Getenv("TOKEN_SECRET")), nil
			})

			if err != nil {
				fmt.Println("Websocket - Missing or invalid token")
				continue
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				var u WebsocketUser

				if role, ok := claims["role"].(string); ok {
					u.Role = role
				} else {
					u.Role = ""
				}
				if userID, ok := claims["user_id"].(string); ok {
					u.UserID = userID
				} else {
					u.UserID = ""
				}

				if claims["role"] == "admin" {
					admins = append(admins, u)
				} else {
					students = append(students, u)
				}
			}
		}
	}
}
