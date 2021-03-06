package test

import (
	"fmt"
	"github.com/simon987/task_tracker/api"
	"github.com/simon987/task_tracker/client"
	"github.com/simon987/task_tracker/storage"
	"net/http"
	"strings"
	"testing"
)

func TestCreateGetWorker(t *testing.T) {

	resp := createWorker(api.CreateWorkerRequest{
		Alias: "my_worker_alias",
	})
	w := resp.Content.Worker

	if resp.Ok != true {
		t.Error()
	}

	getResp := getWorker(w.Id)

	if w.Id != getResp.Content.Worker.Id {
		t.Error()
	}

	if w.Alias != "my_worker_alias" {
		t.Error()
	}

	if w.Paused != false {
		t.Error()
	}
}

func TestGetWorkerNotFound(t *testing.T) {

	resp := getWorker(99999999)

	if resp.Ok != false {
		t.Error()
	}
}

func TestGetWorkerInvalid(t *testing.T) {

	resp := getWorker(-1)

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

	w := getWorker(wid.Id).Content.Worker

	if w.Alias != "new alias" {
		t.Error()
	}

	if w.Paused != false {
		t.Error()
	}
}

func TestCreateWorkerAliasInvalid(t *testing.T) {

	resp := createWorker(api.CreateWorkerRequest{
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
	}).Content.Id

	r := requestAccess(api.CreateWorkerAccessRequest{
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

func TestAssignTaskWhenPaused(t *testing.T) {

	w := genWid()

	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "testassigntaskwhenpaused",
		CloneUrl: "testassigntaskwhenpaused",
		GitRepo:  "testassigntaskwhenpaused",
	}).Content.Id

	requestAccess(api.CreateWorkerAccessRequest{
		Submit:  true,
		Assign:  true,
		Project: pid,
	}, w)
	acceptAccessRequest(pid, w.Id, testAdminCtx)

	r := createTask(api.SubmitTaskRequest{
		Project: pid,
		Recipe:  "a",
		Hash64:  1,
	}, w)

	if r.Ok != true {
		t.Error()
	}

	pauseWorker(&api.WorkerSetPausedRequest{
		Paused: true,
		Worker: w.Id,
	}, testAdminCtx)

	resp := getTaskFromProject(pid, w)

	if resp.Ok != false {
		t.Error()
	}
	if !strings.Contains(resp.Message, "paused") {
		t.Error()
	}
}

func TestPauseInvalidWorker(t *testing.T) {

	r := pauseWorker(&api.WorkerSetPausedRequest{
		Paused: true,
		Worker: 9999111,
	}, testAdminCtx)

	if r.Ok != false {
		t.Error()
	}
}

func TestPauseUnauthorized(t *testing.T) {

	w := genWid()

	r := pauseWorker(&api.WorkerSetPausedRequest{
		Paused: true,
		Worker: w.Id,
	}, testUserCtx)

	if r.Ok != false {
		t.Error()
	}
}

func createWorker(req api.CreateWorkerRequest) (ar client.CreateWorkerResponse) {
	r := Post("/worker/create", req, nil, nil)
	UnmarshalResponse(r, &ar)
	return
}

func getWorker(id int64) (ar client.CreateWorkerResponse) {
	r := Get(fmt.Sprintf("/worker/get/%d", id), nil, nil)
	UnmarshalResponse(r, &ar)
	return
}

func genWid() *storage.Worker {
	resp := createWorker(api.CreateWorkerRequest{})
	return resp.Content.Worker
}

func requestAccess(req api.CreateWorkerAccessRequest, w *storage.Worker) (ar client.CreateWorkerResponse) {
	r := Post(fmt.Sprintf("/project/request_access"), req, w, nil)
	UnmarshalResponse(r, &ar)
	return
}

func acceptAccessRequest(pid int64, wid int64, s *http.Client) (ar api.JsonResponse) {
	r := Post(fmt.Sprintf("/project/accept_request/%d/%d", pid, wid), nil,
		nil, s)
	UnmarshalResponse(r, &ar)
	return
}

func rejectAccessRequest(pid int64, wid int64, s *http.Client) (ar api.JsonResponse) {
	r := Post(fmt.Sprintf("/project/reject_request/%d/%d", pid, wid), nil,
		nil, s)
	UnmarshalResponse(r, &ar)
	return
}

func updateWorker(request api.UpdateWorkerRequest, w *storage.Worker) (ar api.JsonResponse) {
	r := Post("/worker/update", request, w, nil)
	UnmarshalResponse(r, &ar)
	return
}

func pauseWorker(request *api.WorkerSetPausedRequest, s *http.Client) (ar api.JsonResponse) {
	r := Post("/worker/set_paused", request, nil, s)
	UnmarshalResponse(r, &ar)
	return
}
