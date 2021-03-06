package test

import (
	"fmt"
	"github.com/simon987/task_tracker/api"
	"github.com/simon987/task_tracker/client"
	"github.com/simon987/task_tracker/storage"
	"math"
	"testing"
)

func TestCreateTaskValid(t *testing.T) {

	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "Some Test name",
		Version:  "Test Version",
		CloneUrl: "http://github.com/test/test",
		GitRepo:  "Some git repo",
	}).Content.Id

	worker := genWid()
	requestAccess(api.CreateWorkerAccessRequest{
		Project: pid,
		Submit:  true,
		Assign:  false,
	}, worker)
	acceptAccessRequest(pid, worker.Id, testAdminCtx)

	resp := createTask(api.SubmitTaskRequest{
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

	resp := createTask(api.SubmitTaskRequest{
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

	resp := getTaskFromProject(testProject, genWid())

	if resp.Ok != false {
		t.Error()
	}

	if len(resp.Message) <= 0 {
		t.Error()
	}
}

func TestGetTaskInvalidWorker(t *testing.T) {

	resp := getTaskFromProject(testProject, &storage.Worker{
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

	resp := createTask(api.SubmitTaskRequest{
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

	resp := createTask(api.SubmitTaskRequest{
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

	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Name:       "My project",
		Version:    "1.0",
		CloneUrl:   "http://github.com/test/test",
		GitRepo:    "myrepo",
		Priority:   999,
		Public:     true,
		AssignRate: 2,
		SubmitRate: 2,
	}).Content.Id

	worker := genWid()
	requestAccess(api.CreateWorkerAccessRequest{
		Submit:  true,
		Assign:  true,
		Project: pid,
	}, worker)
	acceptAccessRequest(pid, worker.Id, testAdminCtx)

	createTask(api.SubmitTaskRequest{
		Project:           pid,
		Recipe:            "{\"url\":\"test\"}",
		MaxRetries:        3,
		Priority:          9999,
		VerificationCount: 12,
	}, worker)

	taskResp := getTaskFromProject(pid, worker).Content

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
	if taskResp.Task.Project.Id != pid {
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
	if taskResp.Task.Project.AssignRate == 1 {
		t.Error()
	}
	if taskResp.Task.Project.SubmitRate != 2 {
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
	}).Content
	highP := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     prefix + "high",
		Version:  "1.0",
		CloneUrl: "http://github.com/test/test",
		GitRepo:  prefix + "high1",
		Priority: 999,
		Public:   true,
	}).Content
	worker := genWid()
	requestAccess(api.CreateWorkerAccessRequest{
		Submit:  true,
		Assign:  false,
		Project: highP.Id,
	}, worker)
	acceptAccessRequest(highP.Id, worker.Id, testAdminCtx)
	requestAccess(api.CreateWorkerAccessRequest{
		Submit:  true,
		Assign:  false,
		Project: lowP.Id,
	}, worker)
	acceptAccessRequest(lowP.Id, worker.Id, testAdminCtx)

	createTask(api.SubmitTaskRequest{
		Project:  lowP.Id,
		Recipe:   "low1",
		Priority: 0,
	}, worker)
	createTask(api.SubmitTaskRequest{
		Project:  lowP.Id,
		Recipe:   "low2",
		Priority: 1,
	}, worker)
	createTask(api.SubmitTaskRequest{
		Project:  highP.Id,
		Recipe:   "high1",
		Priority: 100,
	}, worker)
	createTask(api.SubmitTaskRequest{
		Project:  highP.Id,
		Recipe:   "high2",
		Priority: 101,
	}, worker)

	return lowP.Id, highP.Id
}

func TestTaskProjectPriority(t *testing.T) {

	wid := genWid()
	l, h := createTasks("withProject")

	t1 := getTaskFromProject(l, wid).Content
	t2 := getTaskFromProject(l, wid).Content
	t3 := getTaskFromProject(h, wid).Content
	t4 := getTaskFromProject(h, wid).Content

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
	}).Content.Id

	requestAccess(api.CreateWorkerAccessRequest{
		Project: pid,
		Assign:  true,
		Submit:  true,
	}, worker)
	acceptAccessRequest(pid, worker.Id, testAdminCtx)

	createResp := createTask(api.SubmitTaskRequest{
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
	if tResp.Content.Task != nil {
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
	}).Content.Id

	requestAccess(api.CreateWorkerAccessRequest{
		Submit:  true,
		Assign:  true,
		Project: pid,
	}, worker)
	acceptAccessRequest(pid, worker.Id, testAdminCtx)

	createResp := createTask(api.SubmitTaskRequest{
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
	if tResp.Content.Task == nil {
		t.Error()
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
	}).Content.Id

	requestAccess(api.CreateWorkerAccessRequest{
		Project: pid,
		Assign:  true,
		Submit:  true,
	}, worker)
	acceptAccessRequest(pid, worker.Id, testAdminCtx)

	createTask(api.SubmitTaskRequest{
		Priority:   0,
		Project:    pid,
		Recipe:     "{}",
		MaxRetries: 3,
		Hash64:     math.MaxInt64,
	}, worker)

	task := getTaskFromProject(pid, worker).Content.Task

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
	}).Content.Id

	w := genWid()
	requestAccess(api.CreateWorkerAccessRequest{
		Project: pid,
		Assign:  true,
		Submit:  true,
	}, w)
	acceptAccessRequest(pid, w.Id, testAdminCtx)

	if createTask(api.SubmitTaskRequest{
		Project:  pid,
		Hash64:   123,
		Priority: 1,
		Recipe:   "{}",
	}, w).Ok != true {
		t.Error()
	}

	resp := createTask(api.SubmitTaskRequest{
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
	}).Content.Id

	w := genWid()
	requestAccess(api.CreateWorkerAccessRequest{
		Project: pid,
		Assign:  true,
		Submit:  true,
	}, w)
	acceptAccessRequest(pid, w.Id, testAdminCtx)

	if createTask(api.SubmitTaskRequest{
		Project:      pid,
		UniqueString: "Hello, world",
		Priority:     1,
		Recipe:       "{}",
	}, w).Ok != true {
		t.Error()
	}

	resp := createTask(api.SubmitTaskRequest{
		Project:      pid,
		UniqueString: "Hello, world",
		Priority:     1,
		Recipe:       "{}",
	}, w)

	if !createTask(api.SubmitTaskRequest{
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
	}).Content.Id

	w := genWid()
	requestAccess(api.CreateWorkerAccessRequest{
		Project: pid,
		Assign:  true,
		Submit:  true,
	}, w)
	acceptAccessRequest(pid, w.Id, testAdminCtx)

	createTask(api.SubmitTaskRequest{
		VerificationCount: 2,
		Project:           pid,
		Recipe:            "verifysametasktwice",
	}, w)

	task := getTaskFromProject(pid, w).Content.Task
	rlr := releaseTask(api.ReleaseTaskRequest{
		Result:       storage.TR_OK,
		TaskId:       task.Id,
		Verification: 123,
	}, w).Content

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
	}).Content.Id

	w := genWid()
	w2 := genWid()
	w3 := genWid()
	requestAccess(api.CreateWorkerAccessRequest{
		Project: pid,
		Assign:  true,
		Submit:  true,
	}, w)
	requestAccess(api.CreateWorkerAccessRequest{
		Project: pid,
		Assign:  true,
		Submit:  true,
	}, w2)
	requestAccess(api.CreateWorkerAccessRequest{
		Project: pid,
		Assign:  true,
		Submit:  true,
	}, w3)
	acceptAccessRequest(pid, w.Id, testAdminCtx)
	acceptAccessRequest(pid, w2.Id, testAdminCtx)
	acceptAccessRequest(pid, w3.Id, testAdminCtx)

	createTask(api.SubmitTaskRequest{
		VerificationCount: 2,
		Project:           pid,
		Recipe:            "verify2",
	}, w)

	task := getTaskFromProject(pid, w).Content.Task
	rlr := releaseTask(api.ReleaseTaskRequest{
		Result:       storage.TR_OK,
		TaskId:       task.Id,
		Verification: 123,
	}, w).Content

	if rlr.Updated != false {
		t.Error()
	}

	task2 := getTaskFromProject(pid, w2).Content.Task
	rlr2 := releaseTask(api.ReleaseTaskRequest{
		Result:       storage.TR_OK,
		Verification: 1,
		TaskId:       task2.Id,
	}, w2).Content

	if rlr2.Updated != false {
		t.Error()
	}

	task3 := getTaskFromProject(pid, w3).Content.Task
	rlr3 := releaseTask(api.ReleaseTaskRequest{
		Result:       storage.TR_OK,
		Verification: 123,
		TaskId:       task3.Id,
	}, w3).Content

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
	}).Content.Id

	w := genWid()
	requestAccess(api.CreateWorkerAccessRequest{
		Project: pid,
		Assign:  true,
		Submit:  true,
	}, w)
	acceptAccessRequest(pid, w.Id, testAdminCtx)

	createTask(api.SubmitTaskRequest{
		MaxRetries:        0,
		Project:           pid,
		VerificationCount: 1,
		Recipe:            "releasefail",
	}, w)

	task := getTaskFromProject(pid, w).Content.Task

	resp := releaseTask(api.ReleaseTaskRequest{
		Result:       storage.TR_FAIL,
		TaskId:       task.Id,
		Verification: 1,
	}, w)

	if resp.Content.Updated != true {
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
	}).Content.Id

	p2 := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "testtaskchain2",
		Public:   true,
		GitRepo:  "testtaskchain2",
		CloneUrl: "testtaskchain2",
		Chain:    p1,
	}).Content.Id
	requestAccess(api.CreateWorkerAccessRequest{
		Project: p1,
		Assign:  true,
		Submit:  true,
	}, w)
	requestAccess(api.CreateWorkerAccessRequest{
		Project: p2,
		Assign:  true,
		Submit:  true,
	}, w)
	acceptAccessRequest(p1, w.Id, testAdminCtx)
	acceptAccessRequest(p2, w.Id, testAdminCtx)

	createTask(api.SubmitTaskRequest{
		Project:           p2,
		Recipe:            "###",
		VerificationCount: 0,
	}, w)

	t1 := getTaskFromProject(p2, w).Content.Task

	releaseTask(api.ReleaseTaskRequest{
		TaskId: t1.Id,
		Result: storage.TR_OK,
	}, w)

	chained := getTaskFromProject(p1, w).Content.Task

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

func TestTaskReleaseBigInt(t *testing.T) {

	createTask(api.SubmitTaskRequest{
		Project:           testProject,
		VerificationCount: 1,
		Recipe:            "bigint",
	}, testWorker)
	createTask(api.SubmitTaskRequest{
		Project:           testProject,
		VerificationCount: 1,
		Recipe:            "smallint",
	}, testWorker)

	tid := getTaskFromProject(testProject, testWorker).Content.Task.Id
	tid2 := getTaskFromProject(testProject, testWorker).Content.Task.Id

	r := releaseTask(api.ReleaseTaskRequest{
		Verification: math.MaxInt64,
		Result:       storage.TR_OK,
		TaskId:       tid,
	}, testWorker)

	r2 := releaseTask(api.ReleaseTaskRequest{
		Verification: math.MinInt64,
		Result:       storage.TR_OK,
		TaskId:       tid2,
	}, testWorker)

	if r.Content.Updated != true {
		t.Error()
	}
	if r2.Content.Updated != true {
		t.Error()
	}
}

func TestTaskSubmitUnauthorized(t *testing.T) {

	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "testtasksubmitunauthorized",
		GitRepo:  "testtasksubmitunauthorized",
		CloneUrl: "testtasksubmitunauthorized",
	}).Content.Id

	w := genWid()

	requestAccess(api.CreateWorkerAccessRequest{
		Project: pid,
		Submit:  true,
		Assign:  true,
	}, w)

	resp := createTask(api.SubmitTaskRequest{
		Project: pid,
		Recipe:  "ssss",
	}, w)

	if resp.Ok != false {
		t.Error()
	}
}

func TestTaskGetUnauthorized(t *testing.T) {

	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "testtaskgetunauthorized",
		GitRepo:  "testtaskgetunauthorized",
		CloneUrl: "testtaskgettunauthorized",
		Hidden:   true,
	}).Content.Id

	w := genWid()
	wWithAccess := genWid()

	requestAccess(api.CreateWorkerAccessRequest{
		Project: pid,
		Submit:  true,
		Assign:  true,
	}, wWithAccess)
	acceptAccessRequest(pid, wWithAccess.Id, testAdminCtx)

	createTask(api.SubmitTaskRequest{
		Project: pid,
		Recipe:  "ssss",
	}, wWithAccess)

	requestAccess(api.CreateWorkerAccessRequest{
		Project: pid,
		Submit:  true,
		Assign:  true,
	}, w)

	resp := getTaskFromProject(pid, w)

	fmt.Println(resp.Message)
	if resp.Ok != false {
		t.Error()
	}
}

func TestTaskChainCausesConflict(t *testing.T) {
	p1 := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "testtaskchainconflict",
		CloneUrl: "testtaskchainconfflict",
		Public:   false,
	}).Content.Id

	p2 := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "testtaskchainconflict2",
		CloneUrl: "testtaskchainconfflict2",
		Public:   false,
		Chain:    p1,
	}).Content.Id

	w := genWid()
	requestAccess(api.CreateWorkerAccessRequest{
		Project: p1,
		Assign:  true,
		Submit:  true,
	}, w)
	requestAccess(api.CreateWorkerAccessRequest{
		Project: p2,
		Assign:  true,
		Submit:  true,
	}, w)
	acceptAccessRequest(p1, w.Id, testAdminCtx)
	acceptAccessRequest(p2, w.Id, testAdminCtx)

	createTask(api.SubmitTaskRequest{
		Project: p2,
		Recipe:  "  ",
		Hash64:  1,
	}, w)
	createTask(api.SubmitTaskRequest{
		Project: p1,
		Recipe:  "  ",
		Hash64:  1,
	}, w)
	tid := getTaskFromProject(p2, w).Content.Task.Id
	resp := releaseTask(api.ReleaseTaskRequest{
		TaskId: tid,
		Result: storage.TR_OK,
	}, w)

	if resp.Ok != true {
		t.Error()
	}
}

func TestTaskAssignInvalidDoesntGiveRateLimit(t *testing.T) {

	task := getTaskFromProject(13247, testWorker)

	if task.RateLimitDelay != 0 {
		t.Error()
	}
}

func TestTaskSubmitInvalidDoesntGiveRateLimit(t *testing.T) {

	resp := createTask(api.SubmitTaskRequest{
		Recipe:  "  ",
		Project: 133453,
	}, testWorker)

	if resp.RateLimitDelay != 0 {
		t.Error()
	}
}

func TestBulkTaskSubmitValid(t *testing.T) {

	r := bulkSubmitTask(api.BulkSubmitTaskRequest{
		Requests: []api.SubmitTaskRequest{
			{
				Recipe:  "1234",
				Project: testProject,
			},
			{
				Recipe:  "1234",
				Project: testProject,
			},
			{
				Recipe:  "1234",
				Project: testProject,
				Hash64:  8565956259293726066,
			},
		},
	}, testWorker)

	if r.Ok != true {
		t.Error()
	}
}

func TestBulkTaskSubmitNotTheSameProject(t *testing.T) {

	proj := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "testbulksubmitnotprj",
		CloneUrl: "testbulkprojectsubmitnotprj",
		GitRepo:  "testbulkprojectsubmitnotprj",
	}).Content.Id

	r := bulkSubmitTask(api.BulkSubmitTaskRequest{
		Requests: []api.SubmitTaskRequest{
			{
				Recipe:  "1234",
				Project: proj,
			},
			{
				Recipe:  "1234",
				Project: 348729,
			},
		},
	}, testWorker)

	if r.Ok != false {
		t.Error()
	}
}

func TestBulkTaskSubmitInvalid(t *testing.T) {

	proj := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "testbulksubmitinvalid",
		CloneUrl: "testbulkprojectsubmitinvalid",
		GitRepo:  "testbulkprojectsubmitinvalid",
	}).Content.Id

	r := bulkSubmitTask(api.BulkSubmitTaskRequest{
		Requests: []api.SubmitTaskRequest{
			{
				Recipe:  "1234",
				Project: proj,
			},
			{

				Recipe:  "",
				Project: proj,
			},
		},
	}, testWorker)

	if r.Ok != false {
		t.Error()
	}
}

func TestBulkTaskSubmitInvalid2(t *testing.T) {

	r := bulkSubmitTask(api.BulkSubmitTaskRequest{
		Requests: []api.SubmitTaskRequest{},
	}, testWorker)

	if r.Ok != false {
		t.Error()
	}
}

func TestTaskGetUnauthorizedWithCache(t *testing.T) {

	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "testtaskgetunauthorizedcache",
		GitRepo:  "testtaskgetunauthorizedcache",
		CloneUrl: "testtaskgettunauthorizedcache",
	}).Content.Id

	w := genWid()

	requestAccess(api.CreateWorkerAccessRequest{
		Project: pid,
		Submit:  true,
		Assign:  true,
	}, w)
	acceptAccessRequest(pid, w.Id, testAdminCtx)

	r1 := createTask(api.SubmitTaskRequest{
		Project: pid,
		Recipe:  "ssss",
	}, w)

	// removed access, cache should be invalidated
	rejectAccessRequest(pid, w.Id, testAdminCtx)

	r2 := createTask(api.SubmitTaskRequest{
		Project: pid,
		Recipe:  "ssss",
	}, w)

	if r1.Ok != true {
		t.Error()
	}
	if r2.Ok != false {
		t.Error()
	}
}

func bulkSubmitTask(request api.BulkSubmitTaskRequest, worker *storage.Worker) (ar api.JsonResponse) {
	r := Post("/task/bulk_submit", request, worker, nil)
	UnmarshalResponse(r, &ar)
	return
}

func createTask(request api.SubmitTaskRequest, worker *storage.Worker) (ar api.JsonResponse) {
	r := Post("/task/submit", request, worker, nil)
	UnmarshalResponse(r, &ar)
	return
}

func getTask(worker *storage.Worker) (ar client.AssignTaskResponse) {
	r := Get("/task/get", worker, nil)
	UnmarshalResponse(r, &ar)
	return
}

func getTaskFromProject(project int64, worker *storage.Worker) (ar client.AssignTaskResponse) {
	r := Get(fmt.Sprintf("/task/get/%d", project), worker, nil)
	UnmarshalResponse(r, &ar)
	return
}

func releaseTask(request api.ReleaseTaskRequest, worker *storage.Worker) (ar client.ReleaseTaskResponse) {
	r := Post("/task/release", request, worker, nil)
	UnmarshalResponse(r, &ar)
	return
}
