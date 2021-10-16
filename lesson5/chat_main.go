package main

import (
	"encoding/json"
	"fmt"
	chi "github.com/go-chi/chi/v5"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

const PORT string = ":8080"
const workingFolder string = "D:/GO/tfs-go-hw/lesson5/"

type PostBody struct {
	Message string `json:"message"`
}

func main() {
	r := chi.NewRouter()
	r.HandleFunc("/", mainHandler)
	r.Post("/messageSend", messageSend)
	log.Fatal(http.ListenAndServe(PORT, r))
}

func mainHandler(w http.ResponseWriter, r *http.Request) {
	file, err := ioutil.ReadFile(workingFolder + "index.html")
	if err != nil {
		w.WriteHeader(500)
		log.Fatal("Bad read file")
		return
	}

	_, err = w.Write(file)
	if err != nil {
		w.WriteHeader(500)
		log.Fatal("Bad sended response")
	}
}

func messageSend(w http.ResponseWriter, r *http.Request) {
	var message PostBody
	d, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	if err = json.Unmarshal(d, &message); err != nil {
		log.Fatal("Bad decoding")
		return
	}

	fmt.Println(message)
}
