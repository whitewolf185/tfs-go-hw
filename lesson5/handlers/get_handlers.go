package handlers

import (
	"github.com/go-chi/chi/v5"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func GetMessageHandler(w http.ResponseWriter, r *http.Request) {
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
	chatPath := GetFilepath(AuthorisedUsername)

	chat, err := ioutil.ReadFile(chatPath)
	if os.IsNotExist(err) {
		_, err := os.Create(chatPath)
		if err != nil {
			log.Fatal("bad crete file")
		}
	}

	_, err = w.Write(chat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("bad write.\n Error", err)
	}
}

func GetChatMessagesHandler(w http.ResponseWriter, r *http.Request) {
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

	chatPath := GetFilepath("main")

	chat, err := ioutil.ReadFile(chatPath)
	if os.IsNotExist(err) {
		_, err := os.Create(chatPath)
		if err != nil {
			log.Fatal("bad crete file")
		}
	}

	_, err = w.Write(chat)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("bad write.\n Error", err)
	}
}
