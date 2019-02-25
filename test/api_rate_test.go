package test

import (
	"fmt"
	"github.com/simon987/task_tracker/api"
	"testing"
)

func TestAssignRateLimit(t *testing.T) {

	project := createProjectAsAdmin(api.CreateProjectRequest{
		SubmitRate: 2,
		AssignRate: 2,
		Name:       "testassignratelimit",
		GitRepo:    "testassignratelimit",
		CloneUrl:   "testassignratelimit",
	}).Content.Id

	w := genWid()
	requestAccess(api.CreateWorkerAccessRequest{
		Project: project,
		Submit:  true,
		Assign:  true,
	}, w)
	acceptAccessRequest(project, w.Id, testAdminCtx)

	for i := 0; i < 3; i++ {
		createTask(api.SubmitTaskRequest{
			Project: project,
			Recipe:  fmt.Sprintf("%d", i),
		}, w)
	}

	var lastResp TaskAR
	for i := 0; i < 3; i++ {
		lastResp = getTaskFromProject(project, w)
	}

	if lastResp.Ok != false {
		t.Error()
	}
	if len(lastResp.Message) <= 0 {
		t.Error()
	}
}

func TestSubmitRateLimit(t *testing.T) {

	project := createProjectAsAdmin(api.CreateProjectRequest{
		SubmitRate: 2,
		AssignRate: 2,
		Name:       "testsubmitratlimit",
		GitRepo:    "testsubmitratelimit",
		CloneUrl:   "testsubmitratelimit",
	}).Content.Id

	w := genWid()
	requestAccess(api.CreateWorkerAccessRequest{
		Project: project,
		Submit:  true,
		Assign:  true,
	}, w)
	acceptAccessRequest(project, w.Id, testAdminCtx)

	var lastResp api.JsonResponse
	for i := 0; i < 2; i++ {
		lastResp = createTask(api.SubmitTaskRequest{
			Project: project,
			Recipe:  fmt.Sprintf("%d", i),
		}, w)
	}

	if lastResp.Ok != false {
		t.Error()
	}
	if len(lastResp.Message) <= 0 {
		t.Error()
	}

}
