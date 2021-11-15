package handlers

import (
	"context"
	"net/http"
	"sync"
)

type Users struct {
	mutex *sync.Mutex
	users map[string]string
}

func NewUsers() *Users {
	var user Users
	user.users = make(map[string]string)

	return &user
}

type Messages struct {
	messages map[string]string
	mutex    *sync.Mutex
}

func NewMessages() *Messages {
	var msg Messages
	msg.messages = make(map[string]string)

	return &msg
}

type Storage struct {
	users    Users
	messages Messages
}

func NewStorage() *Storage {
	users := NewUsers()
	messages := NewMessages()
	return &Storage{users: *users, messages: *messages}
}

func (obj *Users) SetTokenToUser(token, user string) {
	obj.mutex.Lock()
	obj.users[token] = user
	obj.mutex.Unlock()
}

func (obj *Messages) GetMessage(user string) string {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	return obj.messages[user]
}

func (obj Storage) Auth(handler http.Handler) http.Handler {
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
		obj.users.mutex.Lock()
		username, ok := obj.users.users[c.Value]
		obj.users.mutex.Unlock()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		idCtx := context.WithValue(r.Context(), userID, cookieVal(username))

		handler.ServeHTTP(w, r.WithContext(idCtx))
	}

	return http.HandlerFunc(fn)
}

func (obj *Messages) WriteMessage(bodyMessage SendMessage) error {

	obj.mutex.Lock()
	obj.messages[bodyMessage.SendTo] = bodyMessage.Message
	obj.mutex.Unlock()

	return nil
}
