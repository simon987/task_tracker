package test

import (
	"encoding/json"
	"fmt"
	"github.com/simon987/task_tracker/api"
	"github.com/simon987/task_tracker/storage"
	"io/ioutil"
	"testing"
)

func TestCreateTaskValid(t *testing.T) {

	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "Some Test name",
		Version:  "Test Version",
		CloneUrl: "http://github.com/test/test",
		GitRepo:  "Some git repo",
	}).Id

	worker := genWid()
	requestAccess(api.WorkerAccessRequest{
		Project: pid,
		Submit:  true,
		Assign:  false,
	}, worker)
	acceptAccessRequest(pid, worker.Id, testAdminCtx)

	resp := createTask(api.CreateTaskRequest{
		Project:    pid,
		Recipe:     "{}",
		MaxRetries: 3,
	}, worker)

	if resp.Ok != true {
		t.Error()
	}
}

func TestCreateTaskInvalidProject(t *testing.T) {

	worker := genWid()

	resp := createTask(api.CreateTaskRequest{
		Project:    123456,
		Recipe:     "{}",
		MaxRetries: 3,
	}, worker)

	if resp.Ok != false {
		t.Error()
	}

	if len(resp.Message) <= 0 {
		t.Error()
	}
}

func TestGetTaskInvalidWid(t *testing.T) {

	resp := getTask(nil)

	if resp.Ok != false {
		t.Error()
	}

	if len(resp.Message) <= 0 {
		t.Error()
	}
}

func TestGetTaskInvalidWorker(t *testing.T) {

	resp := getTask(&storage.Worker{
		Id: -1,
	})

	if resp.Ok != false {
		t.Error()
	}

	if len(resp.Message) <= 0 {
		t.Error()
	}
}

func TestGetTaskFromProjectInvalidWorker(t *testing.T) {

	resp := getTaskFromProject(1, &storage.Worker{
		Id: 99999999,
	})

	if resp.Ok != false {
		t.Error()
	}

	if len(resp.Message) <= 0 {
		t.Error()
	}
}

func TestCreateTaskInvalidRetries(t *testing.T) {

	worker := genWid()

	resp := createTask(api.CreateTaskRequest{
		Project:    1,
		MaxRetries: -1,
	}, worker)

	if resp.Ok != false {
		t.Error()
	}

	if len(resp.Message) <= 0 {
		t.Error()
	}
}

func TestCreateTaskInvalidRecipe(t *testing.T) {

	worker := genWid()

	resp := createTask(api.CreateTaskRequest{
		Project:    1,
		Recipe:     "",
		MaxRetries: 3,
	}, worker)

	if resp.Ok != false {
		t.Error()
	}

	if len(resp.Message) <= 0 {
		t.Error()
	}
}

func TestCreateGetTask(t *testing.T) {

	//Make sure there is always a project for id:1
	resp := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "My project",
		Version:  "1.0",
		CloneUrl: "http://github.com/test/test",
		GitRepo:  "myrepo",
		Priority: 999,
		Public:   true,
	})

	worker := genWid()
	requestAccess(api.WorkerAccessRequest{
		Submit:  true,
		Assign:  true,
		Project: resp.Id,
	}, worker)
	acceptAccessRequest(resp.Id, worker.Id, testAdminCtx)

	createTask(api.CreateTaskRequest{
		Project:           resp.Id,
		Recipe:            "{\"url\":\"test\"}",
		MaxRetries:        3,
		Priority:          9999,
		VerificationCount: 12,
	}, worker)

	taskResp := getTaskFromProject(resp.Id, worker)

	if taskResp.Ok != true {
		t.Error()
	}
	if taskResp.Task.VerificationCount != 12 {
		t.Error()
	}
	if taskResp.Task.Priority != 9999 {
		t.Error()
	}
	if taskResp.Task.Id == 0 {
		t.Error()
	}
	if string(taskResp.Task.Recipe) != "{\"url\":\"test\"}" {
		t.Error()
	}
	if taskResp.Task.Status != 1 {
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
	if taskResp.Task.Project.CloneUrl != "http://github.com/test/test" {
		t.Error()
	}
	if taskResp.Task.Project.Public != true {
		t.Error()
	}
}

func createTasks(prefix string) (int64, int64) {

	lowP := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     prefix + "low",
		Version:  "1.0",
		CloneUrl: "http://github.com/test/test",
		GitRepo:  prefix + "low1",
		Priority: 1,
		Public:   true,
	})
	highP := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     prefix + "high",
		Version:  "1.0",
		CloneUrl: "http://github.com/test/test",
		GitRepo:  prefix + "high1",
		Priority: 999,
		Public:   true,
	})
	worker := genWid()
	requestAccess(api.WorkerAccessRequest{
		Submit:  true,
		Assign:  false,
		Project: highP.Id,
	}, worker)
	acceptAccessRequest(highP.Id, worker.Id, testAdminCtx)
	requestAccess(api.WorkerAccessRequest{
		Submit:  true,
		Assign:  false,
		Project: lowP.Id,
	}, worker)
	acceptAccessRequest(lowP.Id, worker.Id, testAdminCtx)

	createTask(api.CreateTaskRequest{
		Project:  lowP.Id,
		Recipe:   "low1",
		Priority: 0,
	}, worker)
	createTask(api.CreateTaskRequest{
		Project:  lowP.Id,
		Recipe:   "low2",
		Priority: 1,
	}, worker)
	createTask(api.CreateTaskRequest{
		Project:  highP.Id,
		Recipe:   "high1",
		Priority: 100,
	}, worker)
	createTask(api.CreateTaskRequest{
		Project:  highP.Id,
		Recipe:   "high2",
		Priority: 101,
	}, worker)

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

func TestTaskNoAccess(t *testing.T) {

	worker := genWid()

	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "This is a private proj",
		Motd:     "private",
		Version:  "private",
		Priority: 1,
		CloneUrl: "fjkslejf cesl",
		GitRepo:  "fffffffff",
		Public:   false,
	}).Id

	requestAccess(api.WorkerAccessRequest{
		Project: pid,
		Assign:  true,
		Submit:  true,
	}, worker)
	acceptAccessRequest(worker.Id, pid, testAdminCtx)

	createResp := createTask(api.CreateTaskRequest{
		Project:       pid,
		Priority:      1,
		MaxAssignTime: 10,
		MaxRetries:    2,
		Recipe:        "---",
	}, worker)

	if createResp.Ok != true {
		t.Error()
	}

	rejectAccessRequest(pid, worker.Id, testAdminCtx)

	tResp := getTaskFromProject(pid, worker)

	if tResp.Ok != false {
		t.Error()
	}
	if len(tResp.Message) <= 0 {
		t.Error()
	}
	if tResp.Task != nil {
		t.Error()
	}
}

func TestTaskHasAccess(t *testing.T) {

	worker := genWid()

	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "This is a private proj1",
		Motd:     "private1",
		Version:  "private1",
		Priority: 1,
		CloneUrl: "josaeiuf cesl",
		GitRepo:  "wewwwwwwwwwwwwwwwwwwwwww",
		Public:   false,
	}).Id

	requestAccess(api.WorkerAccessRequest{
		Submit:  true,
		Assign:  true,
		Project: pid,
	}, worker)
	acceptAccessRequest(worker.Id, pid, testAdminCtx)

	createResp := createTask(api.CreateTaskRequest{
		Project:       pid,
		Priority:      1,
		MaxAssignTime: 10,
		MaxRetries:    2,
		Recipe:        "---",
	}, worker)

	if createResp.Ok != true {
		t.Error()
	}

	tResp := getTaskFromProject(pid, worker)

	if tResp.Ok != true {
		t.Error()
	}
	if tResp.Task == nil {
		t.Error()
	}
}

func TestNoMoreTasks(t *testing.T) {

	worker := genWid()

	for i := 0; i < 15; i++ {
		getTask(worker)
	}
}

func TestReleaseTaskSuccess(t *testing.T) {

	worker := genWid()

	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Priority: 0,
		GitRepo:  "testreleasetask",
		CloneUrl: "lllllllll",
		Version:  "11111111111111111",
		Name:     "testreleasetask",
		Motd:     "",
		Public:   true,
	}).Id

	requestAccess(api.WorkerAccessRequest{
		Project: pid,
		Assign:  true,
		Submit:  true,
	}, worker)
	acceptAccessRequest(pid, worker.Id, testAdminCtx)

	createTask(api.CreateTaskRequest{
		Priority:   0,
		Project:    pid,
		Recipe:     "{}",
		MaxRetries: 3,
	}, worker)

	task := getTaskFromProject(pid, worker).Task

	releaseResp := releaseTask(api.ReleaseTaskRequest{
		TaskId: task.Id,
		Result: storage.TR_OK,
	}, worker)

	if releaseResp.Ok != true {
		t.Error()
	}

	otherTask := getTaskFromProject(pid, worker)

	//Shouldn't have more tasks available
	if otherTask.Ok != false {
		t.Error()
	}
}

func TestCreateIntCollision(t *testing.T) {

	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Priority: 1,
		GitRepo:  "testcreateintcollision",
		CloneUrl: "testcreateintcollision",
		Motd:     "testcreateintcollision",
		Public:   true,
		Name:     "testcreateintcollision",
		Version:  "testcreateintcollision",
	}).Id

	w := genWid()
	requestAccess(api.WorkerAccessRequest{
		Project: pid,
		Assign:  true,
		Submit:  true,
	}, w)
	acceptAccessRequest(pid, w.Id, testAdminCtx)

	if createTask(api.CreateTaskRequest{
		Project:  pid,
		Hash64:   123,
		Priority: 1,
		Recipe:   "{}",
	}, w).Ok != true {
		t.Error()
	}

	resp := createTask(api.CreateTaskRequest{
		Project:  pid,
		Hash64:   123,
		Priority: 1,
		Recipe:   "{}",
	}, w)

	if resp.Ok != false {
		t.Error()
	}

	fmt.Println(resp.Message)
	if len(resp.Message) <= 0 {
		t.Error()
	}
}

func TestCreateStringCollision(t *testing.T) {

	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Priority: 1,
		GitRepo:  "testcreatestringcollision",
		CloneUrl: "testcreatestringcollision",
		Motd:     "testcreatestringcollision",
		Public:   true,
		Name:     "testcreatestringcollision",
		Version:  "testcreatestringcollision",
	}).Id

	w := genWid()
	requestAccess(api.WorkerAccessRequest{
		Project: pid,
		Assign:  true,
		Submit:  true,
	}, w)
	acceptAccessRequest(pid, w.Id, testAdminCtx)

	if createTask(api.CreateTaskRequest{
		Project:      pid,
		UniqueString: "Hello, world",
		Priority:     1,
		Recipe:       "{}",
	}, w).Ok != true {
		t.Error()
	}

	resp := createTask(api.CreateTaskRequest{
		Project:      pid,
		UniqueString: "Hello, world",
		Priority:     1,
		Recipe:       "{}",
	}, w)

	if !createTask(api.CreateTaskRequest{
		Project:      pid,
		UniqueString: "This one should work",
		Priority:     1,
		Recipe:       "{}",
	}, w).Ok {
		t.Error()
	}

	if resp.Ok != false {
		t.Error()
	}

	fmt.Println(resp.Message)
	if len(resp.Message) <= 0 {
		t.Error()
	}
}

func TestCannotVerifySameTaskTwice(t *testing.T) {

	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Priority: 1,
		GitRepo:  "verifysametasktwice",
		CloneUrl: "verifysametasktwice",
		Motd:     "verifysametasktwice",
		Public:   true,
		Name:     "verifysametasktwice",
		Version:  "verifysametasktwice",
	}).Id

	w := genWid()
	requestAccess(api.WorkerAccessRequest{
		Project: pid,
		Assign:  true,
		Submit:  true,
	}, w)
	acceptAccessRequest(pid, w.Id, testAdminCtx)

	createTask(api.CreateTaskRequest{
		VerificationCount: 2,
		Project:           pid,
		Recipe:            "verifysametasktwice",
	}, w)

	task := getTaskFromProject(pid, w).Task
	rlr := releaseTask(api.ReleaseTaskRequest{
		Result:       storage.TR_OK,
		TaskId:       task.Id,
		Verification: 123,
	}, w)

	if rlr.Updated != false {
		t.Error()
	}

	sameTask := getTaskFromProject(pid, w)

	if sameTask.Ok != false {
		t.Error()
	}
}

func TestVerification2(t *testing.T) {

	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Priority: 1,
		GitRepo:  "verify2",
		CloneUrl: "verify2",
		Motd:     "verify2",
		Public:   true,
		Name:     "verify2",
		Version:  "verify2",
	}).Id

	w := genWid()
	w2 := genWid()
	w3 := genWid()
	requestAccess(api.WorkerAccessRequest{
		Project: pid,
		Assign:  true,
		Submit:  true,
	}, w)
	requestAccess(api.WorkerAccessRequest{
		Project: pid,
		Assign:  true,
		Submit:  true,
	}, w2)
	requestAccess(api.WorkerAccessRequest{
		Project: pid,
		Assign:  true,
		Submit:  true,
	}, w3)
	acceptAccessRequest(pid, w.Id, testAdminCtx)
	acceptAccessRequest(pid, w2.Id, testAdminCtx)
	acceptAccessRequest(pid, w3.Id, testAdminCtx)

	createTask(api.CreateTaskRequest{
		VerificationCount: 2,
		Project:           pid,
		Recipe:            "verify2",
	}, w)

	task := getTaskFromProject(pid, w).Task
	rlr := releaseTask(api.ReleaseTaskRequest{
		Result:       storage.TR_OK,
		TaskId:       task.Id,
		Verification: 123,
	}, w)

	if rlr.Updated != false {
		t.Error()
	}

	task2 := getTaskFromProject(pid, w2).Task
	rlr2 := releaseTask(api.ReleaseTaskRequest{
		Result:       storage.TR_OK,
		Verification: 1,
		TaskId:       task2.Id,
	}, w2)

	if rlr2.Updated != false {
		t.Error()
	}

	task3 := getTaskFromProject(pid, w3).Task
	rlr3 := releaseTask(api.ReleaseTaskRequest{
		Result:       storage.TR_OK,
		Verification: 123,
		TaskId:       task3.Id,
	}, w3)

	if rlr3.Updated != true {
		t.Error()
	}
}

func TestReleaseTaskFail(t *testing.T) {

	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Priority: 1,
		GitRepo:  "releasefail",
		CloneUrl: "releasefail",
		Motd:     "releasefail",
		Public:   true,
		Name:     "releasefail",
		Version:  "releasefail",
	}).Id

	w := genWid()
	requestAccess(api.WorkerAccessRequest{
		Project: pid,
		Assign:  true,
		Submit:  true,
	}, w)
	acceptAccessRequest(pid, w.Id, testAdminCtx)

	createTask(api.CreateTaskRequest{
		MaxRetries:        0,
		Project:           pid,
		VerificationCount: 1,
		Recipe:            "releasefail",
	}, w)

	task := getTaskFromProject(pid, w).Task

	resp := releaseTask(api.ReleaseTaskRequest{
		Result:       storage.TR_FAIL,
		TaskId:       task.Id,
		Verification: 1,
	}, w)

	if resp.Updated != true {
		t.Error()
	}
	if resp.Ok != true {
		t.Error()
	}

}

func TestTaskChain(t *testing.T) {

	w := genWid()

	p1 := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "testtaskchain1",
		Public:   true,
		GitRepo:  "testtaskchain1",
		CloneUrl: "testtaskchain1",
	}).Id

	p2 := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "testtaskchain2",
		Public:   true,
		GitRepo:  "testtaskchain2",
		CloneUrl: "testtaskchain2",
		Chain:    p1,
	}).Id
	requestAccess(api.WorkerAccessRequest{
		Project: p1,
		Assign:  true,
		Submit:  true,
	}, w)
	requestAccess(api.WorkerAccessRequest{
		Project: p2,
		Assign:  true,
		Submit:  true,
	}, w)
	acceptAccessRequest(p1, w.Id, testAdminCtx)
	acceptAccessRequest(p2, w.Id, testAdminCtx)

	createTask(api.CreateTaskRequest{
		Project:           p2,
		Recipe:            "###",
		VerificationCount: 0,
	}, w)

	t1 := getTaskFromProject(p2, w).Task

	releaseTask(api.ReleaseTaskRequest{
		TaskId: t1.Id,
		Result: storage.TR_OK,
	}, w)

	chained := getTaskFromProject(p1, w).Task

	if chained.VerificationCount != t1.VerificationCount {
		t.Error()
	}
	if chained.Recipe != t1.Recipe {
		t.Error()
	}
	if chained.MaxRetries != t1.MaxRetries {
		t.Error()
	}
	if chained.Priority != t1.Priority {
		t.Error()
	}
	if chained.Status != storage.NEW {
		t.Error()
	}
}

func createTask(request api.CreateTaskRequest, worker *storage.Worker) *api.CreateTaskResponse {

	r := Post("/task/create", request, worker, nil)

	var resp api.CreateTaskResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return &resp
}

func getTask(worker *storage.Worker) *api.GetTaskResponse {

	r := Get(fmt.Sprintf("/task/get"), worker, nil)

	var resp api.GetTaskResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return &resp
}

func getTaskFromProject(project int64, worker *storage.Worker) *api.GetTaskResponse {

	r := Get(fmt.Sprintf("/task/get/%d", project), worker, nil)

	var resp api.GetTaskResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return &resp
}

func releaseTask(request api.ReleaseTaskRequest, worker *storage.Worker) *api.ReleaseTaskResponse {

	r := Post("/task/release", request, worker, nil)

	var resp api.ReleaseTaskResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return &resp
}
