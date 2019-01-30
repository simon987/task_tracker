package test

import (
	"src/task_tracker/api"
	"src/task_tracker/config"
	"src/task_tracker/storage"
	"strconv"
	"testing"
)

func BenchmarkCreateTaskRemote(b *testing.B) {

	resp := createProject(api.CreateProjectRequest{
		Name:     "BenchmarkCreateTask" + strconv.Itoa(b.N),
		GitRepo:  "benchmark_test" + strconv.Itoa(b.N),
		Version:  "f09e8c9r0w839x0c43",
		CloneUrl: "http://localhost",
	})

	worker := genWid()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		createTask(api.CreateTaskRequest{
			Project:    resp.Id,
			Priority:   1,
			Recipe:     "{}",
			MaxRetries: 1,
		}, worker)
	}
}

func BenchmarkCreateTask(b *testing.B) {

	config.SetupConfig()
	db := storage.Database{}

	project, _ := db.SaveProject(&storage.Project{
		Priority: 1,
		Id:       1,
		Version:  "bmcreatetask",
		Public:   true,
		Motd:     "bmcreatetask",
		Name:     "BenchmarkCreateTask" + strconv.Itoa(b.N),
		GitRepo:  "benchmark_test" + strconv.Itoa(b.N),
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = db.SaveTask(&storage.Task{}, project, 0)
	}
}

func TestTest(t *testing.T) {

	config.SetupConfig()
	db := storage.Database{}

	project, _ := db.SaveProject(&storage.Project{
		Priority: 1,
		Id:       1,
		Version:  "bmcreatetask",
		Public:   true,
		Motd:     "bmcreatetask",
		Name:     "BenchmarkCreateTask",
		GitRepo:  "benchmark_test",
	})

	for i := 0; i < 1000000; i++ {
		_ = db.SaveTask(&storage.Task{}, project, 0)
	}
}
