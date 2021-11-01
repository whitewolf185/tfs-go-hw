package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

// LoginHandler нужно отправлять пост запросом структуру вида username="some user"
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var username User
	if err := json.NewDecoder(r.Body).Decode(&username); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("bad decode")
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Fatal("bad body close")
		}
	}()
	fmt.Println(username.Username)

	token, err := GenerateToken(username.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal(err)
	}

	cookie := &http.Cookie{
		Name:  cookieAuth,
		Value: token,
		Path:  "/",
	}

	log.Println(token)

	obj.users[token] = username.Username
	http.SetCookie(w, cookie)
}

// SendMessageHandler тут нужно отсылать json вида "send_to="some user" & message="some message""
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var bodyMessage SendMessage
	if err := json.NewDecoder(r.Body).Decode(&bodyMessage); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("bad decode")
	}
	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Fatal("bad body close")
		}
	}()

	username, ok := r.Context().Value(userID).(cookieVal)
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	AuthorisedUsername := chi.URLParam(r, "username")
	if string(username) != AuthorisedUsername {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	_ = WriteToFile(bodyMessage)

	w.WriteHeader(http.StatusOK)
}
