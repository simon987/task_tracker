package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"src/task_tracker/api"
	"src/task_tracker/storage"
	"testing"
)

func TestCreateGetWorker(t *testing.T) {

	resp, r := createWorker(api.CreateWorkerRequest{
		Alias: "my_worker_alias",
	})

	if r.StatusCode != 200 {
		t.Error()
	}

	if resp.Ok != true {
		t.Error()
	}

	getResp, r := getWorker(resp.Worker.Id)

	if r.StatusCode != 200 {
		t.Error()
	}
	if resp.Worker.Id != getResp.Worker.Id {
		t.Error()
	}

	if len(getResp.Worker.Identity.RemoteAddr) <= 0 {
		t.Error()
	}
	if len(getResp.Worker.Identity.UserAgent) <= 0 {
		t.Error()
	}
	if resp.Worker.Alias != "my_worker_alias" {
		t.Error()
	}
}

func TestGetWorkerNotFound(t *testing.T) {

	resp, r := getWorker(99999999)

	if r.StatusCode != 404 {
		t.Error()
	}
	if resp.Ok != false {
		t.Error()
	}
}

func TestGetWorkerInvalid(t *testing.T) {

	resp, r := getWorker(-1)

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

	resp := grantAccess(wid.Id, 38274593)

	if resp.Ok != false {
		t.Error()
	}
	if len(resp.Message) <= 0 {
		t.Error()
	}
}

func TestRemoveAccessFailedProjectConstraint(t *testing.T) {

	worker := genWid()

	resp := removeAccess(worker.Id, 38274593)

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

	resp := removeAccess(0, pid)

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

	resp := removeAccess(0, pid)

	if resp.Ok != false {
		t.Error()
	}
	if len(resp.Message) <= 0 {
		t.Error()
	}
}

func TestUpdateAliasValid(t *testing.T) {

	wid := genWid()

	updateResp := updateWorker(api.UpdateWorkerRequest{
		Alias: "new alias",
	}, wid)

	if updateResp.Ok != true {
		t.Error()
	}

	w, _ := getWorker(wid.Id)

	if w.Worker.Alias != "new alias" {
		t.Error()
	}
}

func TestCreateWorkerAliasInvalid(t *testing.T) {

	resp, _ := createWorker(api.CreateWorkerRequest{
		Alias: "unassigned", //reserved alias
	})

	if resp.Ok != false {
		t.Error()
	}

	if len(resp.Message) <= 0 {
		t.Error()
	}

}

func createWorker(req api.CreateWorkerRequest) (*api.CreateWorkerResponse, *http.Response) {
	r := Post("/worker/create", req, nil)

	var resp *api.CreateWorkerResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return resp, r
}

func getWorker(id int64) (*api.GetWorkerResponse, *http.Response) {

	r := Get(fmt.Sprintf("/worker/get/%d", id), nil)

	var resp *api.GetWorkerResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return resp, r
}

func genWid() *storage.Worker {

	resp, _ := createWorker(api.CreateWorkerRequest{})
	return resp.Worker
}

func grantAccess(wid int64, project int64) *api.WorkerAccessResponse {

	r := Post("/access/grant", api.WorkerAccessRequest{
		WorkerId:  wid,
		ProjectId: project,
	}, nil)

	var resp *api.WorkerAccessResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return resp
}

func removeAccess(wid int64, project int64) *api.WorkerAccessResponse {

	r := Post("/access/remove", api.WorkerAccessRequest{
		WorkerId:  wid,
		ProjectId: project,
	}, nil)

	var resp *api.WorkerAccessResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return resp
}

func updateWorker(request api.UpdateWorkerRequest, w *storage.Worker) *api.UpdateWorkerResponse {

	r := Post("/worker/update", request, w)

	var resp *api.UpdateWorkerResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return resp
}
