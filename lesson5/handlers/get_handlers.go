package handlers

import (
	"context"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func (obj *ChatHandlers) GetMessageHandler(w http.ResponseWriter, r *http.Request) {
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

	chat := obj.storage.Messages.GetMessage(AuthorisedUsername)

	_, err := w.Write([]byte(chat))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("bad write.\n Error", err)
	}
}

func (obj *ChatHandlers) GetChatMessagesHandler(w http.ResponseWriter, r *http.Request) {
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

	chat := obj.storage.Messages.GetMessage("main")

	_, err := w.Write([]byte(chat))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Fatal("bad write.\n Error", err)
	}
}

func (obj *ChatHandlers) Auth(handler http.Handler) http.Handler {
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

		username, ok := obj.storage.Users.GetTokenUser(c.Value)

		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		idCtx := context.WithValue(r.Context(), userID, cookieVal(username))

		handler.ServeHTTP(w, r.WithContext(idCtx))
	}

	return http.HandlerFunc(fn)
}
