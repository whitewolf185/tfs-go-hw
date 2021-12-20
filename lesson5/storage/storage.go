package storage

import (
	"sync"
)

type Users struct {
	mutex sync.Mutex
	users map[string]string
}

func NewUsers() *Users {
	var user Users
	user.users = make(map[string]string)

	return &user
}

type Messages struct {
	messages map[string]string
	mutex    sync.Mutex
}

func NewMessages() *Messages {
	var msg Messages
	msg.messages = make(map[string]string)

	return &msg
}

type Storage struct {
	Users    *Users
	Messages *Messages
}

func NewStorage() *Storage {
	users := NewUsers()
	messages := NewMessages()
	return &Storage{Users: users, Messages: messages}
}

func (obj *Users) SetTokenToUser(token, user string) {
	obj.mutex.Lock()
	obj.users[token] = user
	obj.mutex.Unlock()
}

func (obj *Users) GetTokenUser(Value string) (string, bool) {
	obj.mutex.Lock()
	username, ok := obj.users[Value]
	obj.mutex.Unlock()
	return username, ok
}

func (obj *Messages) GetMessage(user string) string {
	obj.mutex.Lock()
	defer obj.mutex.Unlock()
	return obj.messages[user]
}

func (obj *Messages) WriteMessage(SendTo string, Message string) {
	obj.mutex.Lock()
	obj.messages[SendTo] += Message + "\n"
	obj.mutex.Unlock()
}
