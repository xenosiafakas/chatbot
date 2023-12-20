package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xenosiafakas/review-chatbot/pkg/routes"
)


func main() {
	
	r := mux.NewRouter()

	routes.ChatBotRoutes(r)
	http.Handle("/",r)

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", http.StatusNotFound)
	})
	
	if err := http.ListenAndServe(":5000", r); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}

