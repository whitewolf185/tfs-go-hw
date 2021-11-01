package handlers

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

const (
	cookieAuth              = "token"
	userID        cookieVal = "username"
	PORT          string    = ":8080"
	workingFolder string    = "D:/Documents/tfs-go-hw/lesson5/"
	tokenTLL                = 24 * time.Hour
	signingKey              = "f&vafqJ=~seYCHcU;Fg?"
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
