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

func TestCreateGetWorker(t *testing.T) {

	resp, r := createWorker(api.CreateWorkerRequest{})

	if r.StatusCode != 200 {
		t.Error()
	}

	if resp.Ok != true {
		t.Error()
	}

	getResp, r := getWorker(resp.WorkerId.String())

	if r.StatusCode != 200 {
		t.Error()
	}
	if resp.WorkerId != getResp.Worker.Id {
		t.Error()
	}

	if len(getResp.Worker.Identity.RemoteAddr) <= 0 {
		t.Error()
	}
	if len(getResp.Worker.Identity.UserAgent) <= 0 {
		t.Error()
	}
}

func TestGetWorkerNotFound(t *testing.T) {

	resp, r := getWorker("8bfc0ccd-d5ce-4dc5-a235-3a7ae760d9c6")

	if r.StatusCode != 404 {
		t.Error()
	}
	if resp.Ok != false {
		t.Error()
	}
}

func TestGetWorkerInvalid(t *testing.T) {

	resp, r := getWorker("invalid-uuid")

	if r.StatusCode != 400 {
		t.Error()
	}
	if resp.Ok != false {
		t.Error()
	}
	if len(resp.Message) <= 0 {
		t.Error()
	}
}

func TestGrantAccessFailedProjectConstraint(t *testing.T) {

	wid := genWid()

	resp := grantAccess(wid, 38274593)

	if resp.Ok != false {
		t.Error()
	}
	if len(resp.Message) <= 0 {
		t.Error()
	}
}

func TestRemoveAccessFailedProjectConstraint(t *testing.T) {

	wid := genWid()

	resp := removeAccess(wid, 38274593)

	if resp.Ok != false {
		t.Error()
	}
	if len(resp.Message) <= 0 {
		t.Error()
	}
}

func TestRemoveAccessFailedWorkerConstraint(t *testing.T) {

	pid := createProject(api.CreateProjectRequest{
		Priority: 1,
		GitRepo:  "dfffffffffff",
		CloneUrl: "fffffffffff23r",
		Version:  "f83w9rw",
		Motd:     "ddddddddd",
		Name:     "removeaccessfailedworkerconstraint",
		Public:   true,
	}).Id

	resp := removeAccess(&uuid.Nil, pid)

	if resp.Ok != false {
		t.Error()
	}
	if len(resp.Message) <= 0 {
		t.Error()
	}
}

func TestGrantAccessFailedWorkerConstraint(t *testing.T) {

	pid := createProject(api.CreateProjectRequest{
		Priority: 1,
		GitRepo:  "dfffffffffff1",
		CloneUrl: "fffffffffff23r1",
		Version:  "f83w9rw1",
		Motd:     "ddddddddd1",
		Name:     "grantaccessfailedworkerconstraint",
		Public:   true,
	}).Id

	resp := removeAccess(&uuid.Nil, pid)

	if resp.Ok != false {
		t.Error()
	}
	if len(resp.Message) <= 0 {
		t.Error()
	}
}

func createWorker(req api.CreateWorkerRequest) (*api.CreateWorkerResponse, *http.Response) {
	r := Post("/worker/create", req)

	var resp *api.CreateWorkerResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return resp, r
}

func getWorker(id string) (*api.GetWorkerResponse, *http.Response) {

	r := Get(fmt.Sprintf("/worker/get/%s", id))

	var resp *api.GetWorkerResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return resp, r
}

func genWid() *uuid.UUID {

	resp, _ := createWorker(api.CreateWorkerRequest{})
	return &resp.WorkerId
}

func grantAccess(wid *uuid.UUID, project int64) *api.WorkerAccessResponse {

	r := Post("/access/grant", api.WorkerAccessRequest{
		WorkerId:  wid,
		ProjectId: project,
	})

	var resp *api.WorkerAccessResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return resp
}

func removeAccess(wid *uuid.UUID, project int64) *api.WorkerAccessResponse {

	r := Post("/access/remove", api.WorkerAccessRequest{
		WorkerId:  wid,
		ProjectId: project,
	})

	var resp *api.WorkerAccessResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return resp
}
