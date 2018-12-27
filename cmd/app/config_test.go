package app

import (
	"encoding/json"
	"log"
	"testing"
	"time"
)

func TestDuration_UnmarshalJSON(t *testing.T) {
	v := `{ "mojangApiQueryInterval": "1h" }`

	var conf Config

	if err := json.Unmarshal([]byte(v), &conf); err != nil {
		log.Println(err)
		t.Fail()
	}

	if conf.MojangAPIQueryInterval.Duration != time.Duration(1 * time.Hour) {
		t.Fail()
	}
}
