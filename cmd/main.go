package main

import (
	"encoding/json"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"gitlab.com/luxordynamics/player-resolver/internal/mojang"
	"log"
)

var api = mojang.NewApi()

type NameResolveRequest struct {
	name string
}

func main() {
	log.SetPrefix("[PlayerResolver] ")
	log.Print("Starting player resolver...")

	/*
		_, b, er := fasthttp.Get(nil, "https://api.mojang.com/users/profiles/minecraft/freggyy")

		if er != nil {
			log.Println(er)
			return
		}

		log.Println(string(b)) */

	router := fasthttprouter.New()
	router.GET("/uuid/:name", HandleUuidRequest)
	router.GET("/name/:uuid", HandleNameRequest)
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

func HandleNameRequest(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	uuid := ctx.UserValue("uuid").(string)

	if mojang.ValidShortUuidRegex.MatchString(uuid) {
		uuid = mojang.ValidLongRegex.ReplaceAllString(uuid, "$1-$2-$3-$4-$5")
	} else if !mojang.ValidLongRegex.MatchString(uuid) {
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBodyString(`{"error": "MalformedUuidException"}`)
		return
	}

	// TODO: check if uuid is already in database
}
