package test

import (
	"github.com/simon987/task_tracker/api"
	"github.com/simon987/task_tracker/client"
	"github.com/simon987/task_tracker/config"
	"github.com/simon987/task_tracker/storage"
	"testing"
)

func TestClientMakeWorker(t *testing.T) {

	c := client.New(config.Cfg.ServerAddr)
	w, err := c.MakeWorker("test")

	if err != nil {
		t.Error()
	}

	if w.Alias != "test" {
		t.Error()
	}
}

func TestClientFetchTaskNoTaskAvailable(t *testing.T) {

	c := client.New(config.Cfg.ServerAddr)
	w, _ := c.MakeWorker("test")
	c.SetWorker(w)

	_, err := c.FetchTask(89988)

	if err != nil {
		t.Error()
	}
}

func TestClientFetchTask(t *testing.T) {

	c := client.New(config.Cfg.ServerAddr)
	w, _ := c.MakeWorker("test")
	c.SetWorker(w)

	createTask(api.SubmitTaskRequest{
		Project: testProject,
		Recipe:  "   ",
	}, testWorker)

	requestAccess(api.CreateWorkerAccessRequest{
		Project: testProject,
		Submit:  false,
		Assign:  true,
	}, &storage.Worker{
		Secret: w.Secret,
		Id:     w.Id,
	})
	acceptAccessRequest(testProject, w.Id, testAdminCtx)

	resp, err := c.FetchTask(int(testProject))

	if err != nil {
		t.Error()
	}

	if resp.Content.Task == nil {
		t.Error()
	}
}

func TestClientReleaseTask(t *testing.T) {

	c := client.New(config.Cfg.ServerAddr)
	w, _ := c.MakeWorker("test")
	c.SetWorker(w)

	createTask(api.SubmitTaskRequest{
		Project: testProject,
		Recipe:  "   ",
	}, testWorker)

	requestAccess(api.CreateWorkerAccessRequest{
		Project: testProject,
		Submit:  false,
		Assign:  true,
	}, &storage.Worker{
		Secret: w.Secret,
		Id:     w.Id,
	})
	acceptAccessRequest(testProject, w.Id, testAdminCtx)

	fetchResp, _ := c.FetchTask(int(testProject))

	resp, err := c.ReleaseTask(api.ReleaseTaskRequest{
		TaskId: fetchResp.Content.Task.Id,
		Result: storage.TR_OK,
	})

	if err != nil {
		t.Error()
	}

	if resp.Content.Updated != true {
		t.Error()
	}
}
