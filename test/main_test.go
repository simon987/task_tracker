package test

import (
	"github.com/simon987/task_tracker/api"
	"github.com/simon987/task_tracker/config"
	"github.com/simon987/task_tracker/storage"
	"net/http"
	"testing"
	"time"
)

var testApi *api.WebAPI
var testAdminCtx *http.Client
var testUserCtx *http.Client

var testProject int64
var testWorker *storage.Worker

func TestMain(m *testing.M) {

	config.SetupConfig()

	testApi = api.New()
	testApi.SetupLogger()
	testApi.Database.Reset()
	go testApi.Run()

	time.Sleep(time.Millisecond * 100)

	testAdminCtx = getSessionCtx("testadmin", "testadmin", true)
	testUserCtx = getSessionCtx("testuser", "testuser", false)
	testProject = createProjectAsAdmin(api.CreateProjectRequest{
		Name:   "generictestproject",
		Public: false,
	}).Content.Id
	testWorker = createWorker(api.CreateWorkerRequest{
		Alias: "generictestworker",
	}).Content.Worker
	requestAccess(api.CreateWorkerAccessRequest{
		Project: testProject,
		Assign:  true,
		Submit:  true,
	}, testWorker)
	acceptAccessRequest(testProject, testWorker.Id, testAdminCtx)

	m.Run()
}
