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

	//Make sure there is always a project for id:1
	createProject(api.CreateProjectRequest{
		Name:     "Some Test name",
		Version:  "Test Version",
		CloneUrl: "http://github.com/test/test",
	})

	worker := genWid()

	resp := createTask(api.CreateTaskRequest{
		Project:    1,
		Recipe:     "{}",
		MaxRetries: 3,
	}, worker)

	if resp.Ok != true {
		t.Fail()
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
	resp := createProject(api.CreateProjectRequest{
		Name:     "My project",
		Version:  "1.0",
		CloneUrl: "http://github.com/test/test",
		GitRepo:  "myrepo",
		Priority: 999,
		Public:   true,
	})

	worker := genWid()

	createTask(api.CreateTaskRequest{
		Project:    resp.Id,
		Recipe:     "{\"url\":\"test\"}",
		MaxRetries: 3,
		Priority:   9999,
	}, worker)

	taskResp := getTaskFromProject(resp.Id, worker)

	if taskResp.Ok != true {
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

	lowP := createProject(api.CreateProjectRequest{
		Name:     prefix + "low",
		Version:  "1.0",
		CloneUrl: "http://github.com/test/test",
		GitRepo:  prefix + "low1",
		Priority: 1,
		Public:   true,
	})
	highP := createProject(api.CreateProjectRequest{
		Name:     prefix + "high",
		Version:  "1.0",
		CloneUrl: "http://github.com/test/test",
		GitRepo:  prefix + "high1",
		Priority: 999,
		Public:   true,
	})
	worker := genWid()
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

	pid := createProject(api.CreateProjectRequest{
		Name:     "This is a private proj",
		Motd:     "private",
		Version:  "private",
		Priority: 1,
		CloneUrl: "fjkslejf cesl",
		GitRepo:  "fffffffff",
		Public:   false,
	}).Id

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

	grantAccess(worker.Id, pid)
	removeAccess(worker.Id, pid)

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

	pid := createProject(api.CreateProjectRequest{
		Name:     "This is a private proj1",
		Motd:     "private1",
		Version:  "private1",
		Priority: 1,
		CloneUrl: "josaeiuf cesl",
		GitRepo:  "wewwwwwwwwwwwwwwwwwwwwww",
		Public:   false,
	}).Id

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

	grantAccess(worker.Id, pid)

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

	pid := createProject(api.CreateProjectRequest{
		Priority: 0,
		GitRepo:  "testreleasetask",
		CloneUrl: "lllllllll",
		Version:  "11111111111111111",
		Name:     "testreleasetask",
		Motd:     "",
		Public:   true,
	}).Id

	createTask(api.CreateTaskRequest{
		Priority:   0,
		Project:    pid,
		Recipe:     "{}",
		MaxRetries: 3,
	}, worker)

	task := getTaskFromProject(pid, worker).Task

	releaseResp := releaseTask(api.ReleaseTaskRequest{
		TaskId:  task.Id,
		Success: true,
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

	pid := createProject(api.CreateProjectRequest{
		Priority: 1,
		GitRepo:  "testcreateintcollision",
		CloneUrl: "testcreateintcollision",
		Motd:     "testcreateintcollision",
		Public:   true,
		Name:     "testcreateintcollision",
		Version:  "testcreateintcollision",
	}).Id

	w := genWid()

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

	pid := createProject(api.CreateProjectRequest{
		Priority: 1,
		GitRepo:  "testcreatestringcollision",
		CloneUrl: "testcreatestringcollision",
		Motd:     "testcreatestringcollision",
		Public:   true,
		Name:     "testcreatestringcollision",
		Version:  "testcreatestringcollision",
	}).Id

	w := genWid()

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

func createTask(request api.CreateTaskRequest, worker *storage.Worker) *api.CreateTaskResponse {

	r := Post("/task/create", request, worker)

	var resp api.CreateTaskResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return &resp
}

func getTask(worker *storage.Worker) *api.GetTaskResponse {

	r := Get(fmt.Sprintf("/task/get"), worker)

	var resp api.GetTaskResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return &resp
}

func getTaskFromProject(project int64, worker *storage.Worker) *api.GetTaskResponse {

	r := Get(fmt.Sprintf("/task/get/%d", project), worker)

	var resp api.GetTaskResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return &resp
}

func releaseTask(request api.ReleaseTaskRequest, worker *storage.Worker) *api.ReleaseTaskResponse {

	r := Post("/task/release", request, worker)

	var resp api.ReleaseTaskResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return &resp
}
