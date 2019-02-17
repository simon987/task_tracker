package test

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"encoding/hex"
	"github.com/simon987/task_tracker/api"
	"github.com/simon987/task_tracker/config"
	"net/http"
	"testing"
)

func TestWebHookNoSignature(t *testing.T) {

	r := Post("/git/receivehook", api.GitPayload{}, nil, nil)

	if r.StatusCode != 403 {
		t.Error()
	}
}

func TestWebHookInvalidSignature(t *testing.T) {

	req, _ := http.NewRequest("POST", "http://"+config.Cfg.ServerAddr+"/git/receivehook", nil)
	req.Header.Add("X-Hub-Signature", "invalid")

	client := http.Client{}
	r, _ := client.Do(req)

	if r.StatusCode != 403 {
		t.Error()
	}
}

func TestWebHookDontUpdateVersion(t *testing.T) {

	resp := createProjectAsAdmin(api.CreateProjectRequest{
		Name:    "My version should not be updated",
		Version: "old",
		GitRepo: "username/not_this_one",
	}).Content

	body := []byte(`{"ref": "refs/heads/master", "after": "new", "repository": {"full_name": "username/repo_name"}}`)
	bodyReader := bytes.NewReader(body)

	mac := hmac.New(crypto.SHA1.New, config.Cfg.WebHookSecret)
	mac.Write(body)
	signature := hex.EncodeToString(mac.Sum(nil))
	signature = "sha1=" + signature

	req, _ := http.NewRequest("POST", "http://"+config.Cfg.ServerAddr+"/git/receivehook", bodyReader)
	req.Header.Add("X-Hub-Signature", signature)

	client := http.Client{}
	r, _ := client.Do(req)

	if r.StatusCode != 200 {
		t.Error()
	}

	getResp := getProjectAsAdmin(resp.Id).Content

	if getResp.Project.Version != "old" {
		t.Error()
	}
}
func TestWebHookUpdateVersion(t *testing.T) {

	resp := createProjectAsAdmin(api.CreateProjectRequest{
		Name:    "My version should be updated",
		Version: "old",
		GitRepo: "username/repo_name",
	}).Content

	body := []byte(`{"ref": "refs/heads/master", "after": "new", "repository": {"full_name": "username/repo_name"}}`)
	bodyReader := bytes.NewReader(body)

	mac := hmac.New(crypto.SHA1.New, config.Cfg.WebHookSecret)
	mac.Write(body)
	signature := hex.EncodeToString(mac.Sum(nil))
	signature = "sha1=" + signature

	req, _ := http.NewRequest("POST", "http://"+config.Cfg.ServerAddr+"/git/receivehook", bodyReader)
	req.Header.Add("X-Hub-Signature", signature)

	client := http.Client{}
	r, _ := client.Do(req)

	if r.StatusCode != 200 {
		t.Error()
	}

	getResp := getProjectAsAdmin(resp.Id).Content

	if getResp.Project.Version != "new" {
		t.Error()
	}
}
