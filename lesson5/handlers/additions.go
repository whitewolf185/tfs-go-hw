package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"hw-chat/storage"
	"time"
)

const (
	cookieAuth           = "token"
	userID     cookieVal = "username"
	PORT       string    = ":8080"
	tokenTLL             = 24 * time.Hour
	signingKey           = "f&vafqJ=~seYCHcU;Fg?"
)

type cookieVal string

type User struct {
	Username string `json:"username"`
}

type SendMessage struct {
	SendTo  string `json:"send_to"`
	Message string `json:"message"`
}

type tokenClaims struct {
	jwt.StandardClaims
	UserName string
}

type ChatHandlers struct {
	storage *storage.Storage
}

func NewChatHandlers() *ChatHandlers {
	var ChatHand ChatHandlers
	ChatHand.storage = storage.NewStorage()

	return &ChatHand
}

func GenerateToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTLL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		username,
	})

	return token.SignedString([]byte(signingKey))
}
