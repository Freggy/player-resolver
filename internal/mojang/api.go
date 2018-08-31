package mojang

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"net/http"
	"regexp"
	"time"
)

var (
	ValidShortUuidRegex = regexp.MustCompile("(\\w{8})(\\w{4})(\\w{4})(\\w{4})(\\w{12})")
	ValidLongRegex      = regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	ValidUserNameRegex  = regexp.MustCompile(`[a-zA-Z0-9_]{1,16}`)
)

// This struct holds the name and the time the player changed to it.
type UuidResolveRequest struct {
	Name        string
	ChangedToAt int64 `json:"changedToAt,omitempty"`
}

// This struct holds the UUID and the name of a player.
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
	_, body, err := fasthttp.Get(nil, "https://api.mojang.com/users/profiles/minecraft/"+name)

	if err != nil {
		return nil, err
	}

	var obj PlayerNameMapping

	if err = json.Unmarshal(body, &obj); err != nil {
		return nil, err
	}

	obj.Uuid = ValidShortUuidRegex.ReplaceAllString(obj.Uuid, "$1-$2-$3-$4-$5")

	return &obj, nil
}

// Resolves the given UUID to the corresponding name.
// This is done by GET https://api.mojang.com/user/profiles/<uuid>/names.
// The given UUID has to be in short form i.e 92de217b8b2b403b86a5fe26fa3a9b5f.
// The return value of this method contains the resolved UUID and the name of the player in the correct spelling.
func (api *Api) NameFromUuid(uuid string) (response *PlayerNameMapping, err error) {
	_, body, err := fasthttp.Get(nil, "https://api.mojang.com/user/profiles/"+uuid+"/names")

	data := make([]UuidResolveRequest, 0)

	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return &PlayerNameMapping{
		ValidShortUuidRegex.ReplaceAllString(uuid, "$1-$2-$3-$4-$5"),
		data[1].Name,
	}, nil
}
