package test

import (
	"testing"
)

func TestIndex(t *testing.T) {

	r := Get("/", nil, nil)
	var info InfoAR
	UnmarshalResponse(r, &info)

	if len(info.Info.Name) <= 0 {
		t.Error()
	}
	if len(info.Info.Version) <= 0 {
		t.Error()
	}
}
