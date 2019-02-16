package test

import (
	"github.com/simon987/task_tracker/api"
	"github.com/simon987/task_tracker/config"
	"net/http"
	"testing"
	"time"
)

var testApi *api.WebAPI
var testAdminCtx *http.Client
var testUserCtx *http.Client

func TestMain(m *testing.M) {

	config.SetupConfig()

	testApi = api.New()
	testApi.SetupLogger()
	testApi.Database.Reset()
	go testApi.Run()

	time.Sleep(time.Millisecond * 100)

	testAdminCtx = getSessionCtx("testadmin", "testadmin", true)
	testUserCtx = getSessionCtx("testuser", "testuser", false)

	m.Run()
}
