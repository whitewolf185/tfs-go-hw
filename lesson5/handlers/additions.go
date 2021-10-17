package handlers

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
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

var users = make(map[string]string)

type tokenClaims struct {
	jwt.StandardClaims
	UserName string
}

func Auth(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(cookieAuth)
		switch err {
		case nil:
		case http.ErrNoCookie:
			w.WriteHeader(http.StatusUnauthorized)
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if c.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		username, ok := users[c.Value]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		idCtx := context.WithValue(r.Context(), userID, cookieVal(username))

		handler.ServeHTTP(w, r.WithContext(idCtx))
	}

	return http.HandlerFunc(fn)
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

func GetFilepath(filename string) string {
	return fmt.Sprintf("%s%s/%s_chat.txt", workingFolder, "chats", filename)
}