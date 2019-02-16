package test

import (
	"bytes"
	"encoding/json"
	"github.com/simon987/task_tracker/api"
	"github.com/simon987/task_tracker/config"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"testing"
)

func TestLoginAndAccountInfo(t *testing.T) {

	regResp := register(&api.RegisterRequest{
		Username: "testusername",
		Password: "testpassword",
	})

	if regResp.Ok != true {
		t.Error()
	}

	loginResp, r := login(&api.LoginRequest{
		Username: "testusername",
		Password: "testpassword",
	})

	if loginResp.Ok != true {
		t.Error()
	}
	if loginResp.Manager.Username != "testusername" {
		t.Error()
	}
	if loginResp.Manager.Id == 0 {
		t.Error()
	}

	ok := false
	for _, c := range r.Cookies() {
		if c.Name == config.Cfg.SessionCookieName {
			ok = true
		}
	}
	if ok != true {
		t.Error()
	}

	url := "http://" + config.Cfg.ServerAddr + "/account"
	req, err := http.NewRequest("GET", url, nil)
	for _, c := range r.Cookies() {
		req.AddCookie(c)
	}

	client := http.Client{}
	r, err = client.Do(req)
	handleErr(err)
	details := &api.AccountDetails{}
	data, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(data, details)
	handleErr(err)

	if details.LoggedIn != true {
		t.Error()
	}
	if details.Manager.Username != "testusername" {
		t.Error()
	}
	if details.Manager.Id != loginResp.Manager.Id {
		t.Error()
	}
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

	r, _ := login(&api.LoginRequest{
		Username: "testinvalidcreds",
		Password: "wrong",
	})

	if r.Ok != false || len(r.Message) <= 0 {
		t.Error()
	}
}

func register(request *api.RegisterRequest) *api.RegisterResponse {

	r := Post("/register", request, nil, nil)

	resp := &api.RegisterResponse{}
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, resp)
	handleErr(err)

	return resp
}

func login(request *api.LoginRequest) (*api.LoginResponse, *http.Response) {

	r := Post("/login", request, nil, nil)

	resp := &api.LoginResponse{}
	data, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(data, resp)
	handleErr(err)

	return resp, r
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
