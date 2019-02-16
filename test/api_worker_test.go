package test

import (
	"encoding/json"
	"fmt"
	"github.com/simon987/task_tracker/api"
	"github.com/simon987/task_tracker/storage"
	"io/ioutil"
	"net/http"
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

func TestInvalidAccessRequest(t *testing.T) {

	w := genWid()
	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "testinvalidaccessreq",
		CloneUrl: "testinvalidaccessreq",
		GitRepo:  "testinvalidaccessreq",
	}).Id

	r := requestAccess(api.WorkerAccessRequest{
		Submit:  false,
		Assign:  false,
		Project: pid,
	}, w)

	if r.Ok != false {
		t.Error()
	}

	if len(r.Message) <= 0 {
		t.Error()
	}
}

func createWorker(req api.CreateWorkerRequest) (*api.CreateWorkerResponse, *http.Response) {
	r := Post("/worker/create", req, nil, nil)

	var resp *api.CreateWorkerResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return resp, r
}

func getWorker(id int64) (*api.GetWorkerResponse, *http.Response) {

	r := Get(fmt.Sprintf("/worker/get/%d", id), nil, nil)

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

func requestAccess(req api.WorkerAccessRequest, w *storage.Worker) *api.WorkerAccessRequestResponse {

	r := Post(fmt.Sprintf("/project/request_access"), req, w, nil)

	var resp *api.WorkerAccessRequestResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return resp
}

func acceptAccessRequest(pid int64, wid int64, s *http.Client) *api.WorkerAccessRequestResponse {

	r := Post(fmt.Sprintf("/project/accept_request/%d/%d", pid, wid), nil,
		nil, s)

	var resp *api.WorkerAccessRequestResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return resp
}

func rejectAccessRequest(pid int64, wid int64, s *http.Client) *api.WorkerAccessRequestResponse {

	r := Post(fmt.Sprintf("/project/reject_request/%d/%d", pid, wid), nil,
		nil, s)

	var resp *api.WorkerAccessRequestResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return resp
}

func updateWorker(request api.UpdateWorkerRequest, w *storage.Worker) *api.UpdateWorkerResponse {

	r := Post("/worker/update", request, w, nil)

	var resp *api.UpdateWorkerResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return resp
}
