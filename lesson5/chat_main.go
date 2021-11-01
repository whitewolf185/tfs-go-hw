package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	hand "hw-chat/handlers"
	"log"
	"net/http"
)

func main() {
	root := chi.NewRouter()
	root.Use(middleware.Logger)
	root.Post("/login", hand.LoginHandler)
	messageHandler := hand.NewStorage()

	r := chi.NewRouter()
	r.Use(messageHandler.Auth)
	r.Get("/users/{username}/message", hand.GetMessageHandler)
	r.Get("/users/{username}/main_chat", hand.GetChatMessagesHandler)

	// чтобы отправить сообщение в вобщий чат, нужно в поле send_to написать main
	r.Post("/users/{username}/message", messageHandler.SendMessageHandler)
	root.Mount("/api", r)

	log.Fatal(http.ListenAndServe(hand.PORT, root))
}
