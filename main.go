package main

import (
	"encoding/json"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"log"
	"strings"
	"gitlab.com/luxordynamics/player-resolver/internal/cassandra"
	"gitlab.com/luxordynamics/player-resolver/internal/mojang"
)

// TODO: change error responses

var api = mojang.NewApi()
var session cassandra.Session

func main() {

	session, err := cassandra.New()

	if err != nil {
		log.Fatal(err)
	}

	defer session.Close()

	router := fasthttprouter.New()
	router.GET("/uuid/:name", HandleUuidRequest)
	router.GET("/name/:uuid", HandleNameRequest)
	fasthttp.ListenAndServe(":8080", router.Handler)
}

// HandleUuidRequest handles requests for resolving names to UUIDs
func HandleUuidRequest(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)

	name := ctx.UserValue("name").(string)

	if !mojang.ValidUserNameRegex.MatchString(name) {
		log.Println("Given name is not valid. (" + name + ")")
		handleError(ctx, fasthttp.StatusBadRequest, `{"code": 400, "message": "Provided name is not valid", "type": "InvalidNameException"}`)
		return
	}

	// TODO: check if uuid is already in database

	mapping, err := api.UuidFromName(name)

	if err != nil {
		log.Fatal(err)
		handleError(ctx, fasthttp.StatusInternalServerError, `{"code": 500, "message": "Error while querying Mojang API", "type": "MojangApiException"}`)
		return
	}

	resp, err := json.Marshal(mapping)

	if err != nil {
		log.Fatal(err)
		handleError(ctx, fasthttp.StatusInternalServerError, `{"code": 500, "message": "Error while processing request", "type": "ServerException"}`)
		return
	}

	ctx.SetBody(resp)
}

// HandleNameRequest handles requests for resolving UUIDs to names
func HandleNameRequest(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	uuid := ctx.UserValue("uuid").(string)

	if mojang.ValidLongRegex.MatchString(uuid) {
		uuid = strings.Replace(uuid, "-", "", -1)
	} else if !mojang.ValidShortUuidRegex.MatchString(uuid) {
		handleError(ctx, fasthttp.StatusBadRequest, `{"code": 400, "message": "Provided UUID is not vaild", "type": "InvalidUuidException"}`)
		return
	}

	exists, err := session.UuidEntryExists(uuid)

	if err != nil {
		handleError(ctx, fasthttp.StatusInternalServerError, `{"code": 500, "message": "Error while querying Mojang API", "type": "MojangApiException"}`)
		return
	}

	var data *mojang.PlayerNameMapping

	if exists {
		data, err = retrieveByUuid(uuid)
	} else {
		data, err = api.NameFromUuid(uuid)
	}

	if err != nil {
		handleError(ctx, fasthttp.StatusInternalServerError, `{"code": 500, "message": "Error while processing request", "type": "ServerException"}`)
		return
	}

	resp, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
		handleError(ctx, fasthttp.StatusInternalServerError, `{"code": 500, "message": "Error while processing request", "type": "ServerException"}`)
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

	return entry.Mapping, nil
}

func handleError(ctx *fasthttp.RequestCtx, code int, body string) {
	ctx.SetStatusCode(code)
	ctx.SetBodyString(body)
}
