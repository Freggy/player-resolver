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

type UuidForNameRequest struct {
	Uuid string `json:"id"`
	Name string
}

type Api struct {
	client *http.Client
}

func NewApi() *Api {
	return &Api{
		&http.Client{Timeout: time.Second * 10},
	}
}

func (api *Api) UuidFromName(name string) (response *UuidForNameRequest, err error) {
	req, err := http.NewRequest("GET", "https://api.mojang.com/users/profiles/minecraft/"+name, nil)

	if err != nil {
		return nil, err
	}

	// Use other user-agent because apparently the Go user-agent is somehow blocked by Mojang for what ever reasons
	req.Header.Set("User-Agent", " Luxor")
	resp, err := api.client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var obj UuidForNameRequest
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
