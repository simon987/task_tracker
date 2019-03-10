package test

import (
	"fmt"
	"github.com/simon987/task_tracker/api"
	"github.com/simon987/task_tracker/storage"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestTraceValid(t *testing.T) {

	w := genWid()
	r := Post("/log/trace", api.LogRequest{
		Scope:     "test",
		Message:   "This is a test message",
		TimeStamp: time.Now().Unix(),
	}, w, nil)

	if r.StatusCode != 200 {
		t.Fail()
	}
}

func TestTraceInvalidScope(t *testing.T) {
	w := genWid()
	r := Post("/log/trace", api.LogRequest{
		Message:   "this is a test message",
		TimeStamp: time.Now().Unix(),
	}, w, nil)

	if r.StatusCode == 200 {
		t.Error()
	}

	r = Post("/log/trace", api.LogRequest{
		Scope:     "",
		Message:   "this is a test message",
		TimeStamp: time.Now().Unix(),
	}, w, nil)

	if r.StatusCode == 200 {
		t.Error()
	}
	if len(GenericJson(r.Body)["message"].(string)) <= 0 {
		t.Error()
	}
}

func TestTraceInvalidMessage(t *testing.T) {
	w := genWid()
	r := Post("/log/trace", api.LogRequest{
		Scope:     "test",
		Message:   "",
		TimeStamp: time.Now().Unix(),
	}, w, nil)

	if r.StatusCode == 200 {
		t.Error()
	}
	if len(GenericJson(r.Body)["message"].(string)) <= 0 {
		t.Error()
	}
}

func TestTraceInvalidTime(t *testing.T) {
	w := genWid()
	r := Post("/log/trace", api.LogRequest{
		Scope:   "test",
		Message: "test",
	}, w, nil)
	if r.StatusCode == 200 {
		t.Error()
	}
	if len(GenericJson(r.Body)["message"].(string)) <= 0 {
		t.Error()
	}
}

func TestWarnValid(t *testing.T) {

	w := genWid()
	r := Post("/log/warn", api.LogRequest{
		Scope:     "test",
		Message:   "test",
		TimeStamp: time.Now().Unix(),
	}, w, nil)
	if r.StatusCode != 200 {
		t.Fail()
	}
}

func TestInfoValid(t *testing.T) {

	w := genWid()
	r := Post("/log/info", api.LogRequest{
		Scope:     "test",
		Message:   "test",
		TimeStamp: time.Now().Unix(),
	}, w, nil)
	if r.StatusCode != 200 {
		t.Fail()
	}
}

func TestErrorValid(t *testing.T) {

	w := genWid()
	r := Post("/log/error", api.LogRequest{
		Scope:     "test",
		Message:   "test",
		TimeStamp: time.Now().Unix(),
	}, w, nil)
	if r.StatusCode != 200 {
		t.Fail()
	}
}

func TestGetLogs(t *testing.T) {

	now := time.Now()

	logrus.WithTime(now.Add(time.Second * -100)).WithFields(logrus.Fields{
		"test": "value",
	}).Debug("This is a test log")

	logrus.WithTime(now.Add(time.Second * -200)).WithFields(logrus.Fields{
		"test": "value",
	}).Debug("This one shouldn't be returned")

	logrus.WithTime(now.Add(time.Second * -100)).WithFields(logrus.Fields{
		"test": "value",
	}).Error("error")

	r := getLogs(time.Now().Add(time.Second*-150).Unix(), storage.DEBUG)

	if r.Ok != true {
		t.Error()
	}

	if len(*r.Content.Logs) <= 0 {
		t.Error()
	}

	debugFound := false
	for _, log := range *r.Content.Logs {
		if log.Message == "This one shouldn't be returned" {
			t.Error()
		} else if log.Message == "error" {
			t.Error()
		} else if log.Message == "This is a test log" {
			debugFound = true
		}
	}

	if !debugFound {
		t.Error()
	}
}

func TestGetLogsInvalid(t *testing.T) {

	r := getLogs(-1, storage.ERROR)

	if r.Ok != false {
		t.Error()
	}

	if len(r.Message) <= 0 {
		t.Error()
	}
}

func getLogs(since int64, level storage.LogLevel) (ar LogsAR) {

	r := Post(fmt.Sprintf("/logs"), api.GetLogRequest{
		Since: since,
		Level: level,
	}, nil, nil)
	UnmarshalResponse(r, &ar)
	return
}
