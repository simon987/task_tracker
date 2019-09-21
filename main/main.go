package main

import (
	"github.com/simon987/task_tracker/api"
	"github.com/simon987/task_tracker/config"
	"math/rand"
	"time"
)

func main() {

	rand.Seed(time.Now().UTC().UnixNano())
	config.SetupConfig()

	webApi := api.New()
	webApi.SetupLogger()
	webApi.Run()
}
