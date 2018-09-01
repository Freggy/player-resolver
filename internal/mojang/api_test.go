package mojang

import "testing"

func TestUuidFromName(t *testing.T) {
	api := NewApi()
	uuid := "92de217b8b2b403b86a5fe26fa3a9b5f"
	mapping, err := api.NameFromUuid(uuid)
	if err != nil {
		t.Error(err)
		return
	}
	failIfNotEqual(mapping, t)
}

func TestNameFromUuid(t *testing.T) {
	api := NewApi()
	name := "freggyy"
	mapping, err := api.UuidFromName(name)
	if err != nil {
		t.Error(err)
		return
	}
	failIfNotEqual(mapping, t)
}

func failIfNotEqual(mapping *PlayerNameMapping, t *testing.T) {
	if mapping.Uuid != "92de217b-8b2b-403b-86a5-fe26fa3a9b5f" {
		t.Fail()
	}

	if mapping.Name != "freggyy" {
		t.Fail()
	}
}
