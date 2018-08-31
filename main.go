package main

import (
	"log"
	"gitlab.com/luxordynamics/player-resolver/internal/mojang"
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
)

var api = mojang.NewApi()

func main() {
	log.SetPrefix("[PlayerResolver] ")
	log.Print("Starting player resolver...")

	router := fasthttprouter.New()
	router.GET("/uuid/:name", HandleUuidRequest)
	router.PUT("/uuid/:name", HandleUuidRequest)
	//router.HandleFunc("/name/{uuid}", HandleNameRequest).Methods("GET", "PUT")
	fasthttp.ListenAndServe(":8080", router.Handler)
}

func HandleUuidRequest(ctx *fasthttp.RequestCtx) {
	go func() {
		if ctx.IsPut() {

		} else {

		}
	}()
}
