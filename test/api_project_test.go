package test

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"src/task_tracker/api"
	"testing"
)

func TestCreateGetProject(t *testing.T) {

	resp := createProject(api.CreateProjectRequest{
		Name:     "Test name",
		CloneUrl: "http://github.com/test/test",
		GitRepo:  "drone/webhooktest",
		Version:  "Test Version",
		Priority: 123,
		Motd:     "motd",
	})

	id := resp.Id

	if id == 0 {
		t.Fail()
	}
	if resp.Ok != true {
		t.Fail()
	}

	getResp, _ := getProject(id)

	if getResp.Project.Id != id {
		t.Error()
	}

	if getResp.Project.Name != "Test name" {
		t.Error()
	}

	if getResp.Project.Version != "Test Version" {
		t.Error()
	}

	if getResp.Project.CloneUrl != "http://github.com/test/test" {
		t.Error()
	}
	if getResp.Project.GitRepo != "drone/webhooktest" {
		t.Error()
	}
	if getResp.Project.Priority != 123 {
		t.Error()
	}
	if getResp.Project.Motd != "motd" {
		t.Error()
	}
}

func TestCreateProjectInvalid(t *testing.T) {
	resp := createProject(api.CreateProjectRequest{})

	if resp.Ok != false {
		t.Fail()
	}
}

func TestCreateDuplicateProjectName(t *testing.T) {
	createProject(api.CreateProjectRequest{
		Name: "duplicate name",
	})
	resp := createProject(api.CreateProjectRequest{
		Name: "duplicate name",
	})

	if resp.Ok != false {
		t.Fail()
	}

	if len(resp.Message) <= 0 {
		t.Fail()
	}
}

func TestCreateDuplicateProjectRepo(t *testing.T) {
	createProject(api.CreateProjectRequest{
		Name:    "different name",
		GitRepo: "user/same",
	})
	resp := createProject(api.CreateProjectRequest{
		Name:    "but same repo",
		GitRepo: "user/same",
	})

	if resp.Ok != false {
		t.Error()
	}

	if len(resp.Message) <= 0 {
		t.Error()
	}
}

func TestGetProjectNotFound(t *testing.T) {

	getResp, r := getProject(12345)

	if getResp.Ok != false {
		t.Fail()
	}

	if len(getResp.Message) <= 0 {
		t.Fail()
	}

	if r.StatusCode != 404 {
		t.Fail()
	}
}

func TestGetProjectStats(t *testing.T) {

	r := createProject(api.CreateProjectRequest{
		Motd:     "motd",
		Name:     "Name",
		Version:  "version",
		CloneUrl: "http://github.com/drone/test",
		GitRepo:  "drone/test",
		Priority: 3,
	})

	pid := r.Id

	createTask(api.CreateTaskRequest{
		Priority:   1,
		Project:    pid,
		MaxRetries: 0,
		Recipe:     "{}",
	})
	createTask(api.CreateTaskRequest{
		Priority:   2,
		Project:    pid,
		MaxRetries: 0,
		Recipe:     "{}",
	})
	createTask(api.CreateTaskRequest{
		Priority:   3,
		Project:    pid,
		MaxRetries: 0,
		Recipe:     "{}",
	})

	stats := getProjectStats(pid)

	if stats.Ok != true {
		t.Error()
	}

	if stats.Stats.Project.Id != pid {
		t.Error()
	}

	if stats.Stats.NewTaskCount != 3 {
		t.Error()
	}

	if stats.Stats.Assignees[0].Assignee != uuid.Nil {
		t.Error()
	}
	if stats.Stats.Assignees[0].TaskCount != 3 {
		t.Error()
	}
}

func TestGetProjectStatsNotFound(t *testing.T) {

	r := createProject(api.CreateProjectRequest{
		Motd:     "eeeeeeeeej",
		Name:     "Namaaaaaaaaaaaa",
		Version:  "versionsssssssss",
		CloneUrl: "http://github.com/drone/test1",
		GitRepo:  "drone/test1",
		Priority: 1,
	})
	s := getProjectStats(r.Id)

	if s.Ok != false {
		t.Error()
	}

	if len(s.Message) <= 0 {
		t.Error()
	}

}

func createProject(req api.CreateProjectRequest) *api.CreateProjectResponse {

	r := Post("/project/create", req)

	var resp api.CreateProjectResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return &resp
}

func getProject(id int64) (*api.GetProjectResponse, *http.Response) {

	r := Get(fmt.Sprintf("/project/get/%d", id))

	var getResp api.GetProjectResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &getResp)
	handleErr(err)

	return &getResp, r
}

func getProjectStats(id int64) *api.GetProjectStatsResponse {

	r := Get(fmt.Sprintf("/project/stats/%d", id))

	var getResp api.GetProjectStatsResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &getResp)
	handleErr(err)

	return &getResp
}
