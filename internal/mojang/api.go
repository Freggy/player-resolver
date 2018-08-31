package mojang

import (
	"net/http"
	"time"
	"io/ioutil"
	"encoding/json"
	"regexp"
)

var (
	ValidShortUuidRegex = regexp.MustCompile("(\\w{8})(\\w{4})(\\w{4})(\\w{4})(\\w{12})")
	ValidLongRegex      = regexp.MustCompile(`[0-9a-f]{32}|[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	ValidUserNameRegex  = regexp.MustCompile(`[a-zA-Z0-9_]{1,16}`)
)

// This struct will hold the UUID and the name of a player.
type PlayerNameMapping struct {
	Uuid string `json:"id"`
	Name string
}

// This struct can be used for accessing the Mojang API for resolving names to UUIDs and vice versa.
type Api struct {
	client *http.Client
}

// Create a new instance of Api.
func NewApi() *Api {
	return &Api{
		&http.Client{Timeout: time.Second * 10},
	}
}

// Resolves the given player name to a UUID.
// This is done by GET https://api.mojang.com/users/profiles/minecraft/<name>.
// The return value of this method contains the resolved UUID and the name of the player in the correct spelling.
func (api *Api) UuidFromName(name string) (response *PlayerNameMapping, err error) {
	req, err := http.NewRequest("GET", "https://api.mojang.com/users/profiles/minecraft/"+name, nil)

	if err != nil {
		return nil, err
	}

	// Use other user-agent because apparently the Go user-agent is somehow blocked by Mojang for what ever reasons
	req.Header.Set("User-Agent", " Luxor (https://www.luxor.cloud)")
	resp, err := api.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var obj PlayerNameMapping
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &obj); err != nil {
		return nil, err
	}

	obj.Uuid = ValidShortUuidRegex.ReplaceAllString(obj.Uuid, "$1-$2-$3-$4-$5")

	return &obj, nil
}
