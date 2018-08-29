package main

import (
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"gitlab.com/luxordynamics/player-resolver/mojang"
)

var api = mojang.NewApi()

func main() {
	log.SetPrefix("[PlayerResolver] ")
	log.Print("Starting player resolver...")

	router := mux.NewRouter()
	router.HandleFunc("/uuid/{name}", HandleUuidRequest).Methods("GET", "PUT")
	router.HandleFunc("/name/{uuid}", HandleNameRequest).Methods("GET", "PUT")
	http.ListenAndServe(":8080", router)
}

func HandleUuidRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Print("GET /uuid/{name}")
	} else {
		log.Print("POST /uuid/{name}")
	}
}

func HandleNameRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		log.Print("GET /name/{uuid}")
	} else {
		log.Print("PUT /name/{uuid}")
	}
}
