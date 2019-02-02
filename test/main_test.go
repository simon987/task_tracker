package test

import (
	"github.com/simon987/task_tracker/api"
	"github.com/simon987/task_tracker/config"
	"testing"
	"time"
)

func TestMain(m *testing.M) {

	config.SetupConfig()

	testApi := api.New()
	testApi.SetupLogger()
	testApi.Database.Reset()
	go testApi.Run()

	time.Sleep(time.Millisecond * 100)
	m.Run()
}
