package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/simon987/task_tracker/api"
	"github.com/simon987/task_tracker/config"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"testing"
)

func TestLoginAndAccountInfo(t *testing.T) {

	//c := getSessionCtx("testusername", "testusername", false)
	//
	//r, _ := c.Get(config.Cfg.ServerAddr + "/account")
	//
	//details := &api.GetAccountDetailsResponse{}
	//data, _ := ioutil.ReadAll(r.Body)
	//err := json.Unmarshal(data, details)
	//handleErr(err)
	//
	//if details.LoggedIn != true {
	//	t.Error()
	//}
	//if details.Manager.Username != "testusername" {
	//	t.Error()
	//}
}

func TestInvalidUsernameRegister(t *testing.T) {

	regResp := register(&api.RegisterRequest{
		Username: "12",
		Password: "testpassword",
	})

	if regResp.Ok != false || len(regResp.Message) <= 0 {
		t.Error()
	}

	regResp2 := register(&api.RegisterRequest{
		Username: "12345678901234567",
		Password: "testpassword",
	})

	if regResp2.Ok != false || len(regResp2.Message) <= 0 {
		t.Error()
	}
}

func TestInvalidPasswordRegister(t *testing.T) {

	regResp := register(&api.RegisterRequest{
		Username: "testinvalidpassword1",
		Password: "12345678",
	})

	if regResp.Ok != false || len(regResp.Message) <= 0 {
		t.Error()
	}
}

func TestDuplicateUsernameRegister(t *testing.T) {

	r1 := register(&api.RegisterRequest{
		Password: "testdupeusername",
		Username: "testdupeusername",
	})

	if r1.Ok != true {
		t.Error()
	}

	r2 := register(&api.RegisterRequest{
		Password: "testdupeusername",
		Username: "testdupeusername",
	})
	if r2.Ok != false || len(r2.Message) <= 0 {
		t.Error()
	}
}

func TestInvalidCredentialsLogin(t *testing.T) {

	register(&api.RegisterRequest{
		Password: "testinvalidcreds",
		Username: "testinvalidcreds",
	})

	r := login(&api.LoginRequest{
		Username: "testinvalidcreds",
		Password: "wrong",
	})

	if r.Ok != false || len(r.Message) <= 0 {
		t.Error()
	}
}

func TestRequireManageAccessRole(t *testing.T) {

	user := getSessionCtx("testreqmanrole", "testreqmanrole", false)

	pid := createProject(api.CreateProjectRequest{
		GitRepo:  "testRequireManageAccessRole",
		CloneUrl: "testRequireManageAccessRole",
		Name:     "testRequireManageAccessRole",
		Version:  "testRequireManageAccessRole",
	}, user).Content.Id

	w := genWid()
	requestAccess(api.CreateWorkerAccessRequest{
		Submit:  true,
		Assign:  true,
		Project: pid,
	}, w)

	rGuest := acceptAccessRequest(pid, w.Id, nil)
	rOtherUser := acceptAccessRequest(pid, w.Id, testUserCtx)
	rUser := acceptAccessRequest(pid, w.Id, user)

	if rGuest.Ok != false {
		t.Error()
	}
	if rOtherUser.Ok != false {
		t.Error()
	}
	if rUser.Ok != true {
		t.Error()
	}

}

func register(request *api.RegisterRequest) (ar RegisterAR) {
	r := Post("/register", request, nil, nil)
	UnmarshalResponse(r, &ar)
	return
}

func login(request *api.LoginRequest) (ar api.JsonResponse) {
	r := Post("/login", request, nil, nil)
	UnmarshalResponse(r, &ar)
	return
}

func getSessionCtx(username string, password string, admin bool) *http.Client {

	register(&api.RegisterRequest{
		Username: username,
		Password: password,
	})

	if admin {
		manager, _ := testApi.Database.ValidateCredentials([]byte(username), []byte(password))
		manager.WebsiteAdmin = true
		testApi.Database.UpdateManager(manager)
	}

	body, err := json.Marshal(api.LoginRequest{
		Username: username,
		Password: password,
	})
	buf := bytes.NewBuffer(body)

	req, err := http.NewRequest("POST", "http://"+config.Cfg.ServerAddr+"/login", buf)
	handleErr(err)

	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client := &http.Client{
		Jar: jar,
	}
	_, err = client.Do(req)
	handleErr(err)

	return client
}

func setRoleOnProject(req api.SetManagerRoleOnProjectRequest, pid int64, s *http.Client) (ar api.JsonResponse) {
	r := Post(fmt.Sprintf("/manager/set_role_for_project/%d", pid), req, nil, s)
	UnmarshalResponse(r, &ar)
	return
}

func getAccountDetails(s *http.Client) (ar AccountAR) {
	r := Get("/account", nil, s)
	UnmarshalResponse(r, &ar)
	return
}
