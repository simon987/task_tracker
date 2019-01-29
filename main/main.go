package main

import (
	"math/rand"
	"src/task_tracker/api"
	"src/task_tracker/config"
	"src/task_tracker/storage"
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
