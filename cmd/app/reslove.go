package app

import (
	"github.com/luxordynamics/player-resolver/util/cassandra"
	"github.com/luxordynamics/player-resolver/util/mojang"
	"log"
	"strings"
	"time"
)

func ResolveNameToUuid(
	name string,
	session *cassandra.Session,
	api *mojang.Api,
	queryInterval time.Duration) (*mojang.PlayerNameMapping, error) {

	if !mojang.ValidUserNameRegex.MatchString(name) {
		log.Println("Given name is not valid. (" + name + ")")
		return nil, NewServerBadRequestError("Provided name is not valid", "InvalidNameException")
	}

	exists, err := session.NameEntryExists(name)

	if err != nil {
		return nil, NewInternalServerError("Error while processing request", "ServerException")
	}

	var data *mojang.PlayerNameMapping

	if exists {
		entry, err := session.EntryByName(name)
		if err != nil {
			return nil, NewInternalServerError("Error while processing request", "ServerException")
		}
		data, err = tryNameRemapping(entry.Mapping.Name, session, api, queryInterval)
	} else {
		data, err = api.UuidFromName(name)
	}
	return data, nil
}

func ResolveUuidToName(
	uuid string,
	session *cassandra.Session,
	api *mojang.Api,
	queryInterval time.Duration) (*mojang.PlayerNameMapping, error) {

	if mojang.ValidLongRegex.MatchString(uuid) {
		uuid = strings.Replace(uuid, "-", "", -1)
	} else if !mojang.ValidShortUuidRegex.MatchString(uuid) {
		return nil, NewServerBadRequestError("Provided UUID is not vaild", "InvalidUUIDException")
	}

	exists, err := session.UuidEntryExists(uuid)

	if err != nil {
		return nil, NewInternalServerError("Error while querying Mojang API", "MojangAPIExcpetion")
	}

	var data *mojang.PlayerNameMapping

	if exists {
		data, err = tryNameRemapping(uuid, session, api, queryInterval)
	} else {
		data, err = api.NameFromUuid(uuid)
	}

	return data, nil
}

// tryNameRemapping TODO: add doc
func tryNameRemapping(uuid string,
	session *cassandra.Session,
	api *mojang.Api,
	queryInterval time.Duration) (mapping *mojang.PlayerNameMapping, err error) {

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
		if time.Unix(entry.LastQuery/1000, 0).After(time.Now().Add(queryInterval)) {
			mapping, err := api.NameFromUuid(entry.Mapping.Uuid)

			if err != nil {
				return nil, err
			}

			if err := session.UpdateLastQuery(time.Now().UnixNano()/1000000, mapping.Uuid); err != nil {
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
