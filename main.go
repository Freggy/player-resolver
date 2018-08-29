package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	log.SetPrefix("[PlayerResolver]")
	log.Print("Starting player resolver...")
	router := mux.NewRouter()
	router.HandleFunc("/uuid", GetUuidFromName).Methods("GET")
	http.ListenAndServe(":8080", router)
}


func GetUuidFromName(w http.ResponseWriter, r *http.Request) {

}