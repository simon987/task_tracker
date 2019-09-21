package test

import (
	"github.com/simon987/task_tracker/api"
	"github.com/simon987/task_tracker/storage"
	"strconv"
	"testing"
)

func BenchmarkCreateTaskRemote(b *testing.B) {

	resp := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "BenchmarkCreateTask" + strconv.Itoa(b.N),
		GitRepo:  "benchmark_test" + strconv.Itoa(b.N),
		Version:  "f09e8c9r0w839x0c43",
		CloneUrl: "http://localhost",
	})

	worker := genWid()

	requestAccess(api.CreateWorkerAccessRequest{
		Submit:  true,
		Assign:  false,
		Project: resp.Content.Id,
	}, worker)
	acceptAccessRequest(resp.Content.Id, worker.Id, testAdminCtx)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		createTask(api.SubmitTaskRequest{
			Project:    resp.Content.Id,
			Priority:   1,
			Recipe:     "{}",
			MaxRetries: 1,
		}, worker)
	}
}

func BenchmarkCreateTask(b *testing.B) {

	resp := createProjectAsAdmin(api.CreateProjectRequest{
		Name:     "BenchmarkCreateTask" + strconv.Itoa(b.N),
		GitRepo:  "benchmark_test" + strconv.Itoa(b.N),
		Version:  "f09e8c9r0w839x0c43",
		CloneUrl: "http://localhost",
	})

	worker := genWid()

	requestAccess(api.CreateWorkerAccessRequest{
		Submit:  true,
		Assign:  false,
		Project: resp.Content.Id,
	}, worker)
	acceptAccessRequest(resp.Content.Id, worker.Id, testAdminCtx)

	db := storage.New()

	b.ResetTimer()

	p := db.GetProject(resp.Content.Id)
	for i := 0; i < b.N; i++ {
		db.SaveTask(&storage.Task{
			Project:    p,
			Priority:   0,
			Recipe:     "{}",
			MaxRetries: 1,
		}, resp.Content.Id, 0, worker.Id)
	}
}
