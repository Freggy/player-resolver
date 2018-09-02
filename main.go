package player_resolver

import (
	"encoding/json"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"gitlab.com/luxordynamics/player-resolver/internal/cassandra"
	"gitlab.com/luxordynamics/player-resolver/internal/mojang"
	"log"
	"strings"
)

var api = mojang.NewApi()
var session cassandra.CassandraSession

func main() {
	log.SetPrefix("[PlayerResolver] ")
	log.Print("Starting player resolver...")

	session, err := cassandra.New()

	if err != nil {
		log.Fatal(err)
	}

	defer session.Close()

	// 20109332-3b1b-4dcb-9fd1-1b3468f05572
	// 92de217b-8b2b-403b-86a5-fe26fa3a9b5f
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

// Handles requests for resolving names to UUIDs
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

	// TODO: check if uuid is already in database

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

// Handles requests for resolving UUIDs to names
func HandleNameRequest(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	uuid := ctx.UserValue("uuid").(string)

	if mojang.ValidLongRegex.MatchString(uuid) {
		uuid = strings.Replace(uuid, "-", "", -1)
	} else if !mojang.ValidShortUuidRegex.MatchString(uuid) {
		handleError(ctx, `{"error": "MalformedUuidException"}`)
		return
	}

	exists, err := session.UuidEntryExists(uuid)

	if err != nil {
		handleError(ctx, `{"error": "InternalServiceException"}`)
		return
	}

	var data *mojang.PlayerNameMapping

	if exists {
		data, err = retrieveByUuid(uuid)
	} else {
		data, err = api.NameFromUuid(uuid)
	}

	if err != nil {
		handleError(ctx, `{"error": "InternalServiceException"}`)
		return
	}

	resp, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
		handleError(ctx, `{"error": "InternalServiceException"}`)
		return
	}

	ctx.SetBody(resp)
	// TODO: check if uuid is already in database
}


func retrieveByUuid(uuid string) (mapping *mojang.PlayerNameMapping, err error) {
	entry, err := session.EntryByUuid(uuid)

	if err != nil {
		return nil, err
	}

	// TODO: check if last update was x days ago

	return &entry.Mapping, nil
}

func handleError(ctx *fasthttp.RequestCtx, body string) {
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	ctx.SetBodyString(body)
}


