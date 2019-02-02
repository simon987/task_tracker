package main

import (
	"github.com/simon987/task_tracker/api"
	"github.com/simon987/task_tracker/config"
	"github.com/simon987/task_tracker/storage"
	"math/rand"
	"time"
)

func tmpDebugSetup() {

	db := storage.Database{}
	db.Reset()

}

func main() {

	rand.Seed(time.Now().UTC().UnixNano())
	config.SetupConfig()

	webApi := api.New()
	webApi.SetupLogger()
	tmpDebugSetup()
	webApi.Run()
}
