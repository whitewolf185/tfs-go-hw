package handlers

import (
	"context"
	"net/http"
	"sync"
)

type StorageService interface {
	SetTokenToUser(token, user string)
	WriteMessage(bodyMessage SendMessage) error
	GetMessage(user string) string
}

type Storage struct {
	mutex    *sync.Mutex
	users    map[string]string
	messages map[string]string
	service  StorageService
}

func NewStorage() *Storage {
	users := make(map[string]string)
	messages := make(map[string]string)
	return &Storage{users: users, messages: messages}
}

func (obj *Storage) SetTokenToUser(token, user string) {
	obj.mutex.Lock()
	obj.users[token] = user
	obj.mutex.Unlock()
}

func (obj *Storage) GetMessage(user string) string {
	return obj.users[user]
}

func (obj *Storage) Auth(handler http.Handler) http.Handler {
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
		username, ok := obj.users[c.Value]
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		idCtx := context.WithValue(r.Context(), userID, cookieVal(username))

		handler.ServeHTTP(w, r.WithContext(idCtx))
	}

	return http.HandlerFunc(fn)
}

func (obj *Storage) WriteMessage(bodyMessage SendMessage) error {

	obj.mutex.Lock()
	obj.messages[bodyMessage.SendTo] = bodyMessage.Message
	obj.mutex.Unlock()

	return nil
}
