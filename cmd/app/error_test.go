package app

import (
	"encoding/json"
	"log"
	"testing"
)

func TestErrorToJsonString(t *testing.T) {
	myErr := NewServiceError(400, "My Error Message", "MyErrorType")
	jsonErr, err := myErr.ToJSON()

	if err != nil {
		log.Println("Error while marshalling")
		t.Fail()
	}

	if string(jsonErr) != `{"code":400,"message":"My Error Message","type":"MyErrorType"}` {
		t.Fail()
	}
}

func TestErrorUnmarshalling(t *testing.T) {
	myErr := NewServiceError(400, "My Error Message", "MyErrorType")
	data, err := myErr.ToJSON()

	if err != nil {
		log.Println("Error while marshalling")
		t.Fail()
	}

	var otherErr ServiceError

	if err := json.Unmarshal([]byte(data), &otherErr); err != nil {
		log.Println("Error while unmarshalling")
		t.Fail()
	}

	if otherErr.Status != myErr.Status {
		t.Fail()
	} else if otherErr.Type != myErr.Type {
		t.Fail()
	} else if otherErr.Message != myErr.Message {
		t.Fail()
	}
}
