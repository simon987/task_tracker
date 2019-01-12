package test

import (
	"encoding/json"
	"io/ioutil"
	"src/task_tracker/api"
	"testing"
)

func TestCreateTaskValid(t *testing.T) {

	//Make sure there is always a project for id:1
	createProject(api.CreateProjectRequest{
		Name: "Some Test name",
		Version: "Test Version",
		GitUrl: "http://github.com/test/test",

	})

	resp := createTask(api.CreateTaskRequest{
		Project:1,
		Recipe: "{}",
		MaxRetries:3,
	})

	if resp.Ok != true {
		t.Fail()
	}
}

func TestCreateTaskInvalidProject(t *testing.T) {

	resp := createTask(api.CreateTaskRequest{
		Project:123456,
		Recipe: "{}",
		MaxRetries:3,
	})

	if resp.Ok != false {
		t.Fail()
	}

	if len(resp.Message) <= 0 {
		t.Fail()
	}
}

func TestCreateTaskInvalidRetries(t *testing.T) {

	resp := createTask(api.CreateTaskRequest{
		Project:1,
		MaxRetries:-1,
	})

	if resp.Ok != false {
		t.Fail()
	}

	if len(resp.Message) <= 0 {
		t.Fail()
	}
}

func createTask(request api.CreateTaskRequest) *api.CreateTaskResponse {

	r := Post("/task/create", request)

	var resp api.CreateTaskResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return &resp
}
