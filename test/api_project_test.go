package test

import (
	"encoding/json"
	"fmt"
	"github.com/simon987/task_tracker/api"
	"io/ioutil"
	"net/http"
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
		Public:   true,
		Hidden:   true,
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
	if getResp.Project.Public != true {
		t.Error()
	}
	if getResp.Project.Hidden != true {
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

func TestUpdateProjectValid(t *testing.T) {

	pid := createProject(api.CreateProjectRequest{
		Public:   true,
		Version:  "versionA",
		Motd:     "MotdA",
		Name:     "NameA",
		CloneUrl: "CloneUrlA",
		GitRepo:  "GitRepoA",
		Priority: 1,
	}).Id

	updateResp := updateProject(api.UpdateProjectRequest{
		Priority: 2,
		GitRepo:  "GitRepoB",
		CloneUrl: "CloneUrlB",
		Name:     "NameB",
		Motd:     "MotdB",
		Public:   false,
		Hidden:   true,
	}, pid)

	if updateResp.Ok != true {
		t.Error()
	}

	proj, _ := getProject(pid)

	if proj.Project.Public != false {
		t.Error()
	}
	if proj.Project.Motd != "MotdB" {
		t.Error()
	}
	if proj.Project.CloneUrl != "CloneUrlB" {
		t.Error()
	}
	if proj.Project.GitRepo != "GitRepoB" {
		t.Error()
	}
	if proj.Project.Priority != 2 {
		t.Error()
	}
	if proj.Project.Hidden != true {
		t.Error()
	}
}

func TestUpdateProjectInvalid(t *testing.T) {

	pid := createProject(api.CreateProjectRequest{
		Public:   true,
		Version:  "lllllllllllll",
		Motd:     "2wwwwwwwwwwwwwww",
		Name:     "aaaaaaaaaaaaaaaaaaaaaa",
		CloneUrl: "333333333333333",
		GitRepo:  "llllllllllllllllllls",
		Priority: 1,
	}).Id

	updateResp := updateProject(api.UpdateProjectRequest{
		Priority: -1,
		GitRepo:  "GitRepo------",
		CloneUrl: "CloneUrlB000000",
		Name:     "NameB-0",
		Motd:     "MotdB000000",
		Public:   false,
	}, pid)

	if updateResp.Ok != false {
		t.Error()
	}

	if len(updateResp.Message) <= 0 {
		t.Error()
	}
}

func TestUpdateProjectConstraintFail(t *testing.T) {

	pid := createProject(api.CreateProjectRequest{
		Public:   true,
		Version:  "testUpdateProjectConstraintFail",
		Motd:     "testUpdateProjectConstraintFail",
		Name:     "testUpdateProjectConstraintFail",
		CloneUrl: "testUpdateProjectConstraintFail",
		GitRepo:  "testUpdateProjectConstraintFail",
		Priority: 1,
	}).Id

	createProject(api.CreateProjectRequest{
		Public:   true,
		Version:  "testUpdateProjectConstraintFail_d",
		Motd:     "testUpdateProjectConstraintFail_d",
		Name:     "testUpdateProjectConstraintFail_d",
		CloneUrl: "testUpdateProjectConstraintFail_d",
		GitRepo:  "testUpdateProjectConstraintFail_d",
		Priority: 1,
	})

	updateResp := updateProject(api.UpdateProjectRequest{
		Priority: 1,
		GitRepo:  "testUpdateProjectConstraintFail_d",
		CloneUrl: "testUpdateProjectConstraintFail_d",
		Name:     "testUpdateProjectConstraintFail_d",
		Motd:     "testUpdateProjectConstraintFail_d",
	}, pid)

	if updateResp.Ok != false {
		t.Error()
	}

	if len(updateResp.Message) <= 0 {
		t.Error()
	}
}

func createProject(req api.CreateProjectRequest) *api.CreateProjectResponse {

	r := Post("/projectChange/create", req, nil)

	var resp api.CreateProjectResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return &resp
}

func getProject(id int64) (*api.GetProjectResponse, *http.Response) {

	r := Get(fmt.Sprintf("/projectChange/get/%d", id), nil)

	var getResp api.GetProjectResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &getResp)
	handleErr(err)

	return &getResp, r
}

func updateProject(request api.UpdateProjectRequest, pid int64) *api.UpdateProjectResponse {

	r := Post(fmt.Sprintf("/projectChange/update/%d", pid), request, nil)

	var resp api.UpdateProjectResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return &resp
}
