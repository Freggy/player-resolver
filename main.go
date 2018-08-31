package main

import (
	"log"
	"gitlab.com/luxordynamics/player-resolver/internal/mojang"
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
	"encoding/json"
)

var api = mojang.NewApi()

type NameResolveRequest struct {
	name string
}

func main() {
	log.SetPrefix("[PlayerResolver] ")
	log.Print("Starting player resolver...")

	router := fasthttprouter.New()
	router.GET("/uuid/:name", HandleUuidRequest)
	//router.HandleFunc("/name/{uuid}", HandleNameRequest).Methods("GET", "PUT")
	fasthttp.ListenAndServe(":8080", router.Handler)
}

// Handles requests for UUID resolving
func HandleUuidRequest(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)

	name := ctx.UserValue("name").(string)

	if !mojang.ValidUserNameRegex.MatchString(name) {
		log.Println("Given name is not valid. (" + name + ")")
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(`{"error": "NameNotValidException"}`)
		return
	}

	mapping, err := api.UuidFromName(name)

	if err != nil {
		log.Fatal(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(`{"error": "MojangRequestException"}`)
		return
	}

	resp, err := json.Marshal(mapping)

	if err != nil {
		log.Fatal(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(`{"error": "ProcessFailedException"}`)
		return
	}

	ctx.SetBody(resp)
}
