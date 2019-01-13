package test

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"src/task_tracker/api"
	"testing"
)

func TestCreateTaskValid(t *testing.T) {

	//Make sure there is always a project for id:1
	createProject(api.CreateProjectRequest{
		Name:    "Some Test name",
		Version: "Test Version",
		GitUrl:  "http://github.com/test/test",
	})

	resp := createTask(api.CreateTaskRequest{
		Project:    1,
		Recipe:     "{}",
		MaxRetries: 3,
	})

	if resp.Ok != true {
		t.Fail()
	}
}

func TestCreateTaskInvalidProject(t *testing.T) {

	resp := createTask(api.CreateTaskRequest{
		Project:    123456,
		Recipe:     "{}",
		MaxRetries: 3,
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
		Project:    1,
		MaxRetries: -1,
	})

	if resp.Ok != false {
		t.Fail()
	}

	if len(resp.Message) <= 0 {
		t.Fail()
	}
}

func TestCreateGetTask(t *testing.T) {

	//Make sure there is always a project for id:1
	resp := createProject(api.CreateProjectRequest{
		Name:     "My project",
		Version:  "1.0",
		GitUrl:   "http://github.com/test/test",
		Priority: 999,
	})

	createTask(api.CreateTaskRequest{
		Project:    resp.Id,
		Recipe:     "{\"url\":\"test\"}",
		MaxRetries: 3,
		Priority:   9999,
	})

	taskResp := getTaskFromProject(resp.Id, genWid())

	if taskResp.Ok != true {
		t.Error()
	}
	if taskResp.Task.Priority != 9999 {
		t.Error()
	}
	if taskResp.Task.Id == 0 {
		t.Error()
	}
	if taskResp.Task.Recipe != "{\"url\":\"test\"}" {
		t.Error()
	}
	if taskResp.Task.Status != "new" {
		t.Error()
	}
	if taskResp.Task.MaxRetries != 3 {
		t.Error()
	}
	if taskResp.Task.Project.Id != resp.Id {
		t.Error()
	}
	if taskResp.Task.Project.Priority != 999 {
		t.Error()
	}
	if taskResp.Task.Project.Version != "1.0" {
		t.Error()
	}
	if taskResp.Task.Project.GitUrl != "http://github.com/test/test" {
		t.Error()
	}
}

func createTasks(prefix string) (int64, int64) {

	lowP := createProject(api.CreateProjectRequest{
		Name:     prefix + "low",
		Version:  "1.0",
		GitUrl:   "http://github.com/test/test",
		Priority: 1,
	})
	highP := createProject(api.CreateProjectRequest{
		Name:     prefix + "high",
		Version:  "1.0",
		GitUrl:   "http://github.com/test/test",
		Priority: 999,
	})
	createTask(api.CreateTaskRequest{
		Project:  lowP.Id,
		Recipe:   "low1",
		Priority: 0,
	})
	createTask(api.CreateTaskRequest{
		Project:  lowP.Id,
		Recipe:   "low2",
		Priority: 1,
	})
	createTask(api.CreateTaskRequest{
		Project:  highP.Id,
		Recipe:   "high1",
		Priority: 100,
	})
	createTask(api.CreateTaskRequest{
		Project:  highP.Id,
		Recipe:   "high2",
		Priority: 101,
	})

	return lowP.Id, highP.Id
}

func TestTaskProjectPriority(t *testing.T) {

	wid := genWid()
	l, h := createTasks("withProject")

	t1 := getTaskFromProject(l, wid)
	t2 := getTaskFromProject(l, wid)
	t3 := getTaskFromProject(h, wid)
	t4 := getTaskFromProject(h, wid)

	if t1.Task.Recipe != "low2" {
		t.Error()
	}
	if t2.Task.Recipe != "low1" {
		t.Error()
	}
	if t3.Task.Recipe != "high2" {
		t.Error()
	}
	if t4.Task.Recipe != "high1" {
		t.Error()
	}
}

func TestTaskPriority(t *testing.T) {

	wid := genWid()

	// Clean other tasks
	for i := 0; i < 20; i++ {
		getTask(wid)
	}

	createTasks("")

	t1 := getTask(wid)
	t2 := getTask(wid)
	t3 := getTask(wid)
	t4 := getTask(wid)

	if t1.Task.Recipe != "high2" {
		t.Error()
	}
	if t2.Task.Recipe != "high1" {
		t.Error()
	}
	if t3.Task.Recipe != "low2" {
		t.Error()
	}
	if t4.Task.Recipe != "low1" {
		t.Error()
	}
}

func TestNoMoreTasks(t *testing.T) {

	wid := genWid()

	for i := 0; i < 30; i++ {
		getTask(wid)
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

func getTask(wid *uuid.UUID) *api.GetTaskResponse {

	r := Get(fmt.Sprintf("/task/get?wid=%s", wid))

	var resp api.GetTaskResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return &resp
}

func getTaskFromProject(project int64, wid *uuid.UUID) *api.GetTaskResponse {

	r := Get(fmt.Sprintf("/task/get/%d?wid=%s", project, wid))

	var resp api.GetTaskResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return &resp
}
