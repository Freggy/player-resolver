package main

import (
	"encoding/json"
	"github.com/buaazp/fasthttprouter"
	"github.com/luxordynamics/player-resolver/cmd/app"
	"github.com/luxordynamics/player-resolver/util/cassandra"
	"github.com/luxordynamics/player-resolver/util/mojang"
	"github.com/valyala/fasthttp"
	"log"
	"time"
)

// TODO: make testable

var api = mojang.NewApi()
var config app.Config
var dbSession *cassandra.Session

func main() {

	session, err := cassandra.New()
	defer session.Close()

	if err != nil {
		log.Fatal(err)
	}

	if err = session.Setup(); err != nil {
		log.Fatal(err)
	}

	// TODO: load config

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
	sendData(app.ResolveNameToUuid, name, dbSession, api, config.MojangAPIQueryInterval, ctx)
}

// HandleNameRequest handles requests for resolving UUIDs to names
func HandleNameRequest(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	uuid := ctx.UserValue("uuid").(string)
	sendData(app.ResolveUuidToName, uuid, dbSession, api, config.MojangAPIQueryInterval, ctx)
}

func sendData(
	f func(string, *cassandra.Session, *mojang.Api, time.Duration) (*mojang.PlayerNameMapping, error),
	identifier string,
	session *cassandra.Session,
	api *mojang.Api,
	queryInterval time.Duration,
	ctx *fasthttp.RequestCtx) {

	data, err := f(identifier, session, api, queryInterval)

	if err != nil {
		handleError(nil, ctx, app.NewInternalServerError("Error while processing request", "ServerException"))
		return
	}

	resp, err := json.Marshal(data)

	if err != nil {
		handleError(nil, ctx, app.NewInternalServerError("Error while processing request", "ServerException"))
		return
	}

	ctx.SetBody(resp)
}

func handleError(err error, ctx *fasthttp.RequestCtx, serviceError *app.ServiceError) {
	log.Print(err)
	ctx.SetStatusCode(serviceError.Status)
	data, _ := serviceError.ToJSON()
	ctx.SetBodyString(data)
}
