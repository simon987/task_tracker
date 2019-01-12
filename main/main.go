package main

import (
	"src/task_tracker/api"
	"src/task_tracker/config"
	"src/task_tracker/storage"
)

func tmpDebugSetup() {

	db := storage.Database{}
	db.Reset()

}

func main() {

	config.SetupConfig()

	tmpDebugSetup()

	webApi := api.New()
	webApi.Run()
}
