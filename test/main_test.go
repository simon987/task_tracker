package test

import (
	"src/task_tracker/api"
	"src/task_tracker/config"
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
