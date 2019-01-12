package test

import (
	"src/task_tracker/api"
	"testing"
	"time"
)


func TestTraceValid(t *testing.T) {

	r := Post("/log/trace", api.LogEntry{
		Scope:"test",
		Message:"This is a test message",
		TimeStamp: time.Now().Unix(),
	})

	if r.StatusCode != 200 {
		t.Fail()
	}
}

func TestTraceInvalidScope(t *testing.T) {
	r := Post("/log/trace", api.LogEntry{
		Message:"this is a test message",
		TimeStamp: time.Now().Unix(),
	})

	if r.StatusCode != 500 {
		t.Fail()
	}

	r = Post("/log/trace", api.LogEntry{
		Scope:"",
		Message:"this is a test message",
		TimeStamp: time.Now().Unix(),
	})

	if r.StatusCode != 500 {
		t.Fail()
	}
	if GenericJson(r.Body)["message"] != "invalid scope" {
		t.Fail()
	}
}

func TestTraceInvalidMessage(t *testing.T) {
	r := Post("/log/trace", api.LogEntry{
		Scope:"test",
		Message:"",
		TimeStamp: time.Now().Unix(),
	})

	if r.StatusCode != 500 {
		t.Fail()
	}
	if GenericJson(r.Body)["message"] != "invalid message" {
		t.Fail()
	}
}

func TestTraceInvalidTime(t *testing.T) {
	r := Post("/log/trace", api.LogEntry{
		Scope: "test",
		Message:"test",

	})
	if r.StatusCode != 500 {
		t.Fail()
	}
	if GenericJson(r.Body)["message"] != "invalid timestamp" {
		t.Fail()
	}
}

func TestWarnValid(t *testing.T) {

	r := Post("/log/warn", api.LogEntry{
		Scope: "test",
		Message:"test",
		TimeStamp:time.Now().Unix(),
	})
	if r.StatusCode != 200 {
		t.Fail()
	}
}

func TestInfoValid(t *testing.T) {

	r := Post("/log/info", api.LogEntry{
		Scope: "test",
		Message:"test",
		TimeStamp:time.Now().Unix(),
	})
	if r.StatusCode != 200 {
		t.Fail()
	}
}

func TestErrorValid(t *testing.T) {

	r := Post("/log/error", api.LogEntry{
		Scope: "test",
		Message:"test",
		TimeStamp:time.Now().Unix(),
	})
	if r.StatusCode != 200 {
		t.Fail()
	}
}
