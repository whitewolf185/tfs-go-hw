package handlers

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func (obj Storage) GetMessageHandler(w http.ResponseWriter, r *http.Request) {
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

	chat := obj.service.GetMessage(AuthorisedUsername)

	_, err := w.Write([]byte(chat))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("bad write.\n Error", err)
	}
}

func (obj Storage) GetChatMessagesHandler(w http.ResponseWriter, r *http.Request) {
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

	chat := obj.service.GetMessage("main")

	_, err := w.Write([]byte(chat))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("bad write.\n Error", err)
	}
}
