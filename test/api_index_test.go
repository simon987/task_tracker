package test

import (
	"encoding/json"
	"io/ioutil"
	"src/task_tracker/api"
	"testing"
)

func TestIndex(t *testing.T) {

	r := Get("/")

	body, _ := ioutil.ReadAll(r.Body)
	var info api.Info
	err := json.Unmarshal(body, &info)

	if err != nil {
		t.Error(err.Error())
	}

	if len(info.Name) <= 0 {
		t.Error()
	}
	if len(info.Version) <= 0 {
		t.Error()
	}
}
