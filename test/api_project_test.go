package test

import (
	"encoding/json"
	"fmt"
	"github.com/simon987/task_tracker/api"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestCreateGetProject(t *testing.T) {

	resp := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "Test name",
		CloneUrl: "http://github.com/test/test",
		GitRepo:  "drone/webhooktest",
		Version:  "Test Version",
		Priority: 123,
		Motd:     "motd",
		Public:   true,
		Hidden:   false,
	})

	id := resp.Id

	if id == 0 {
		t.Fail()
	}
	if resp.Ok != true {
		t.Fail()
	}

	getResp, _ := getProjectAsAdmin(id)

	if getResp.Project.Id != id {
		t.Error()
	}

	if getResp.Project.Name != "Test name" {
		t.Error()
	}

	if getResp.Project.Version != "Test Version" {
		t.Error()
	}

	if getResp.Project.CloneUrl != "http://github.com/test/test" {
		t.Error()
	}
	if getResp.Project.GitRepo != "drone/webhooktest" {
		t.Error()
	}
	if getResp.Project.Priority != 123 {
		t.Error()
	}
	if getResp.Project.Motd != "motd" {
		t.Error()
	}
	if getResp.Project.Public != true {
		t.Error()
	}
	if getResp.Project.Hidden != false {
		t.Error()
	}
}

func TestCreateProjectInvalid(t *testing.T) {
	resp := createProjectAsAdmin(api.CreateProjectRequest{})

	if resp.Ok != false {
		t.Fail()
	}
}

func TestCreateDuplicateProjectName(t *testing.T) {
	createProjectAsAdmin(api.CreateProjectRequest{
		Name: "duplicate name",
	})
	resp := createProjectAsAdmin(api.CreateProjectRequest{
		Name: "duplicate name",
	})

	if resp.Ok != false {
		t.Fail()
	}

	if len(resp.Message) <= 0 {
		t.Fail()
	}
}

func TestCreateDuplicateProjectRepo(t *testing.T) {
	createProjectAsAdmin(api.CreateProjectRequest{
		Name:    "different name",
		GitRepo: "user/same",
	})
	resp := createProjectAsAdmin(api.CreateProjectRequest{
		Name:    "but same repo",
		GitRepo: "user/same",
	})

	if resp.Ok != false {
		t.Error()
	}

	if len(resp.Message) <= 0 {
		t.Error()
	}
}

func TestGetProjectNotFound(t *testing.T) {

	getResp, r := getProjectAsAdmin(12345)

	if getResp.Ok != false {
		t.Fail()
	}

	if len(getResp.Message) <= 0 {
		t.Fail()
	}

	if r.StatusCode != 404 {
		t.Fail()
	}
}

func TestUpdateProjectValid(t *testing.T) {

	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Public:   true,
		Version:  "versionA",
		Motd:     "MotdA",
		Name:     "NameA",
		CloneUrl: "CloneUrlA",
		GitRepo:  "GitRepoA",
		Priority: 1,
	}).Id

	updateResp := updateProject(api.UpdateProjectRequest{
		Priority: 2,
		GitRepo:  "GitRepoB",
		CloneUrl: "CloneUrlB",
		Name:     "NameB",
		Motd:     "MotdB",
		Public:   false,
		Hidden:   true,
	}, pid, testAdminCtx)

	if updateResp.Ok != true {
		t.Error()
	}

	proj, _ := getProjectAsAdmin(pid)

	if proj.Project.Public != false {
		t.Error()
	}
	if proj.Project.Motd != "MotdB" {
		t.Error()
	}
	if proj.Project.CloneUrl != "CloneUrlB" {
		t.Error()
	}
	if proj.Project.GitRepo != "GitRepoB" {
		t.Error()
	}
	if proj.Project.Priority != 2 {
		t.Error()
	}
	if proj.Project.Hidden != true {
		t.Error()
	}
}

func TestUpdateProjectInvalid(t *testing.T) {

	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Public:   true,
		Version:  "lllllllllllll",
		Motd:     "2wwwwwwwwwwwwwww",
		Name:     "aaaaaaaaaaaaaaaaaaaaaa",
		CloneUrl: "333333333333333",
		GitRepo:  "llllllllllllllllllls",
		Priority: 1,
	}).Id

	updateResp := updateProject(api.UpdateProjectRequest{
		Priority: -1,
		GitRepo:  "GitRepo------",
		CloneUrl: "CloneUrlB000000",
		Name:     "NameB-0",
		Motd:     "MotdB000000",
		Public:   false,
	}, pid, testAdminCtx)

	if updateResp.Ok != false {
		t.Error()
	}

	if len(updateResp.Message) <= 0 {
		t.Error()
	}
}

func TestUpdateProjectConstraintFail(t *testing.T) {

	pid := createProjectAsAdmin(api.CreateProjectRequest{
		Public:   true,
		Version:  "testUpdateProjectConstraintFail",
		Motd:     "testUpdateProjectConstraintFail",
		Name:     "testUpdateProjectConstraintFail",
		CloneUrl: "testUpdateProjectConstraintFail",
		GitRepo:  "testUpdateProjectConstraintFail",
		Priority: 1,
	}).Id

	createProjectAsAdmin(api.CreateProjectRequest{
		Public:   true,
		Version:  "testUpdateProjectConstraintFail_d",
		Motd:     "testUpdateProjectConstraintFail_d",
		Name:     "testUpdateProjectConstraintFail_d",
		CloneUrl: "testUpdateProjectConstraintFail_d",
		GitRepo:  "testUpdateProjectConstraintFail_d",
		Priority: 1,
	})

	updateResp := updateProject(api.UpdateProjectRequest{
		Priority: 1,
		GitRepo:  "testUpdateProjectConstraintFail_d",
		CloneUrl: "testUpdateProjectConstraintFail_d",
		Name:     "testUpdateProjectConstraintFail_d",
		Motd:     "testUpdateProjectConstraintFail_d",
	}, pid, testAdminCtx)

	if updateResp.Ok != false {
		t.Error()
	}

	if len(updateResp.Message) <= 0 {
		t.Error()
	}
}

func TestNotLoggedProjectCreate(t *testing.T) {

	r := createProject(api.CreateProjectRequest{
		Hidden:   false,
		Name:     "testnotlogged",
		Priority: 1,
		CloneUrl: "testnotlogged",
		GitRepo:  "testnotlogged",
	}, nil)

	if r.Ok != false {
		t.Error()
	}

	if len(r.Message) <= 0 {
		t.Error()
	}
}

func TestUserCanCreatePrivateProject(t *testing.T) {

	r := createProject(api.CreateProjectRequest{
		Hidden:   false,
		Name:     "testuserprivate",
		Priority: 1,
		CloneUrl: "testuserprivate",
		GitRepo:  "testuserprivate",
		Public:   false,
	}, testUserCtx)

	if r.Ok != true {
		t.Error()
	}
}

func TestUserCannotCreatePublicProject(t *testing.T) {

	r := createProject(api.CreateProjectRequest{
		Hidden:   false,
		Name:     "testuserprivate",
		Priority: 1,
		CloneUrl: "testuserprivate",
		GitRepo:  "testuserprivate",
		Public:   true,
	}, testUserCtx)

	if r.Ok != false {
		t.Error()
	}
	if len(r.Message) <= 0 {
		t.Error()
	}
}

func TestHiddenProjectsNotShownInList(t *testing.T) {

	r := createProject(api.CreateProjectRequest{
		Hidden:   true,
		Name:     "testhiddenprojectlist",
		Priority: 1,
		CloneUrl: "testhiddenprojectlist",
		GitRepo:  "testhiddenprojectlist",
		Public:   false,
	}, testUserCtx)

	if r.Ok != true {
		t.Error()
	}

	list := getProjectList(nil)

	for _, p := range *list.Projects {
		if p.Id == r.Id {
			t.Error()
		}
	}
}

func TestHiddenProjectCannotBePublic(t *testing.T) {

	r := createProject(api.CreateProjectRequest{
		Hidden:   true,
		Name:     "testhiddencannotbepublic",
		Priority: 1,
		CloneUrl: "testhiddencannotbepublic",
		GitRepo:  "testhiddencannotbepublic",
		Public:   true,
	}, testAdminCtx)

	if r.Ok != false {
		t.Error()
	}
	if len(r.Message) <= 0 {
		t.Error()
	}
}

func TestHiddenProjectNotAccessible(t *testing.T) {

	otherUser := getSessionCtx("otheruser", "otheruser", false)
	r := createProject(api.CreateProjectRequest{
		Hidden:   true,
		Name:     "testhiddenprojectaccess",
		Priority: 1,
		CloneUrl: "testhiddenprojectaccess",
		GitRepo:  "testhiddenprojectaccess",
		Public:   false,
	}, testUserCtx)

	if r.Ok != true {
		t.Error()
	}

	pAdmin, _ := getProject(r.Id, testAdminCtx)
	pUser, _ := getProject(r.Id, testUserCtx)
	pOtherUser, _ := getProject(r.Id, otherUser)
	pGuest, _ := getProject(r.Id, nil)

	if pAdmin.Ok != true {
		t.Error()
	}
	if pUser.Ok != true {
		t.Error()
	}
	if pOtherUser.Ok != false {
		t.Error()
	}
	if pGuest.Ok != false {
		t.Error()
	}
}

func TestUpdateProjectPermissions(t *testing.T) {

	p := createProjectAsAdmin(api.CreateProjectRequest{
		GitRepo:  "updateprojectpermissions",
		CloneUrl: "updateprojectpermissions",
		Name:     "updateprojectpermissions",
		Version:  "updateprojectpermissions",
	})

	r := updateProject(api.UpdateProjectRequest{
		GitRepo:  "newupdateprojectpermissions",
		CloneUrl: "newupdateprojectpermissions",
		Name:     "newupdateprojectpermissions",
	}, p.Id, nil)

	if r.Ok != false {
		t.Error()
	}
	if len(r.Message) <= 0 {
		t.Error()
	}
}

func TestUserWithReadAccessShouldSeeHiddenProjectInList(t *testing.T) {

	pHidden := createProject(api.CreateProjectRequest{
		GitRepo:  "testUserHidden",
		CloneUrl: "testUserHidden",
		Name:     "testUserHidden",
		Version:  "testUserHidden",
		Hidden:   true,
	}, testUserCtx)

	list := getProjectList(testUserCtx)

	found := false
	for _, p := range *list.Projects {
		if p.Id == pHidden.Id {
			found = true
		}
	}

	if !found {
		t.Error()
	}
}

func TestAdminShouldSeeHiddenProjectInList(t *testing.T) {

	pHidden := createProject(api.CreateProjectRequest{
		GitRepo:  "testAdminHidden",
		CloneUrl: "testAdminHidden",
		Name:     "testAdminHidden",
		Version:  "testAdminHidden",
		Hidden:   true,
	}, testUserCtx)

	list := getProjectList(testAdminCtx)

	found := false
	for _, p := range *list.Projects {
		if p.Id == pHidden.Id {
			found = true
		}
	}

	if !found {
		t.Error()
	}
}

func createProjectAsAdmin(req api.CreateProjectRequest) *api.CreateProjectResponse {
	return createProject(req, testAdminCtx)
}

func createProject(req api.CreateProjectRequest, s *http.Client) *api.CreateProjectResponse {
	r := Post("/project/create", req, nil, s)

	var resp api.CreateProjectResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return &resp
}

func getProjectAsAdmin(id int64) (*api.GetProjectResponse, *http.Response) {
	return getProject(id, testAdminCtx)
}

func getProject(id int64, s *http.Client) (*api.GetProjectResponse, *http.Response) {

	r := Get(fmt.Sprintf("/project/get/%d", id), nil, s)

	var getResp api.GetProjectResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &getResp)
	handleErr(err)

	return &getResp, r
}

func updateProject(request api.UpdateProjectRequest, pid int64, s *http.Client) *api.UpdateProjectResponse {

	r := Post(fmt.Sprintf("/project/update/%d", pid), request, nil, s)

	var resp api.UpdateProjectResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return &resp
}

func getProjectList(s *http.Client) *api.GetAllProjectsResponse {
	r := Get("/project/list", nil, s)

	var resp api.GetAllProjectsResponse
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, &resp)
	handleErr(err)

	return &resp
}
