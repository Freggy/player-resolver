package main

import (
	"encoding/json"
	"github.com/buaazp/fasthttprouter"
	"github.com/luxordynamics/player-resolver/cmd/app"
	"github.com/luxordynamics/player-resolver/util/cassandra"
	"github.com/luxordynamics/player-resolver/util/mojang"
	"github.com/valyala/fasthttp"
	"log"
	"strings"
	"time"
)

// TODO: change error responses

var api = mojang.NewApi()
var session cassandra.Session
var config app.Config

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
		handleError(nil, ctx, fasthttp.StatusBadRequest, `{"code": 400, "message": "Provided name is not valid", "type": "InvalidNameException"}`)
		return
	}

	exists, err := session.NameEntryExists(name)

	if err != nil {
		handleError(err, ctx, fasthttp.StatusBadRequest, `{"code": 500, "message": "Error while processing request", "type": "ServerException"}`)
		return
	}

	var data *mojang.PlayerNameMapping

	if exists {
		entry, err := session.EntryByName(name)
		if err != nil {
			handleError(err, ctx, fasthttp.StatusBadRequest, `{"code": 500, "message": "Error while processing request", "type": "ServerException"}`)
			return
		}
		data, err = tryNameRemapping(entry.Mapping.Name)
	} else {
		data, err = api.UuidFromName(name)
	}

	if err != nil {
		handleError(err, ctx, fasthttp.StatusInternalServerError, `{"code": 500, "message": "Error while querying Mojang API", "type": "MojangApiException"}`)
		return
	}

	resp, err := json.Marshal(data)

	if err != nil {
		handleError(err, ctx, fasthttp.StatusInternalServerError, `{"code": 500, "message": "Error while processing request", "type": "ServerException"}`)
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
		handleError(nil, ctx, fasthttp.StatusBadRequest, `{"code": 400, "message": "Provided UUID is not vaild", "type": "InvalidUuidException"}`)
		return
	}

	exists, err := session.UuidEntryExists(uuid)

	if err != nil {
		handleError(err, ctx, fasthttp.StatusInternalServerError, `{"code": 500, "message": "Error while querying Mojang API", "type": "MojangApiException"}`)
		return
	}

	var data *mojang.PlayerNameMapping

	if exists {
		data, err = tryNameRemapping(uuid)
	} else {
		data, err = api.NameFromUuid(uuid)
	}

	if err != nil {
		handleError(err, ctx, fasthttp.StatusInternalServerError, `{"code": 500, "message": "Error while processing request", "type": "ServerException"}`)
		return
	}

	resp, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
		handleError(err, ctx, fasthttp.StatusInternalServerError, `{"code": 500, "message": "Error while processing request", "type": "ServerException"}`)
		return
	}

	ctx.SetBody(resp)
}

// tryNameRemapping TODO: add doc
func tryNameRemapping(uuid string) (mapping *mojang.PlayerNameMapping, err error) {
	entry, err := session.EntryByUuid(uuid)

	if err != nil {
		return nil, err
	}

	canChangeDate := time.Unix(entry.Mapping.ChangedToAt/1000, 0).AddDate(0, 1, 0)

	// If the current date is past the date on which
	// the player is able to change their name again
	// query Mojang API
	if canChangeDate.After(time.Now()) {

		// If the last time we queried the Mojang api exceeds the specified interval,
		// we retrieve the newest name in order to have the most up to date values
		if time.Unix(entry.LastQuery/1000, 0).After(time.Now().Add(config.MojangAPIQueryInterval)) {
			mapping, err := api.NameFromUuid(entry.Mapping.Uuid)

			if err != nil {
				return nil, err
			}

			if err := session.UpdateLastQuery(time.Now().UnixNano() / 1000000, mapping.Uuid); err != nil {
				return nil, err
			}

			// Name has changed write changes to database
			if mapping.Name != entry.Mapping.Name {
				if err := session.UpdateName(mapping.Name, mapping.Uuid); err != nil {
					return nil, err
				}
			}
		}
	}

	return entry.Mapping, nil
}

func handleError(err error, ctx *fasthttp.RequestCtx, code int, body string) {
	log.Print(err)
	ctx.SetStatusCode(code)
	ctx.SetBodyString(body)
}
