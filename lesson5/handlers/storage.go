package handlers

import (
	"context"
	"net/http"
	"sync"
)

type Storage struct {
	mutex    *sync.Mutex
	users    map[string]string
	messages map[string]string
}

func NewStorage() Storage {
	users := make(map[string]string)
	messages := make(map[string]string)
	return Storage{users: users, messages: messages}
}

func (obj *Storage) SetTokenToUser(token, user string) {

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
