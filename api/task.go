package api

import (
	"github.com/Sirupsen/logrus"
	"src/task_tracker/storage"
)

type CreateTaskRequest struct {
	Project int64 `json:"project"`
	MaxRetries int64 `json:"max_retries"`
	Recipe string `json:"recipe"`
}

type CreateTaskResponse struct {
	Ok bool
	Message string
}

func (api *WebAPI) TaskCreate(r *Request) {

	var createReq CreateTaskRequest
	if r.GetJson(&createReq) {

		task := &storage.Task{
			Project:createReq.Project,
			MaxRetries: createReq.MaxRetries,
			Recipe:createReq.Recipe,
		}

		if isTaskValid(task) {
			err := api.Database.SaveTask(task)

			if err != nil {
				r.Json(CreateTaskResponse{
					Ok: false,
					Message: err.Error(), //todo: hide sensitive error?
				}, 500)
			} else {
				r.OkJson(CreateTaskResponse{
					Ok: true,
				})
			}
		} else {
			logrus.WithFields(logrus.Fields{
				"task": task,
			}).Warn("Invalid task")
			r.Json(CreateTaskResponse{
				Ok: false,
				Message: "Invalid task",
			}, 400)
		}

	}
}

func isTaskValid(task *storage.Task) bool {
	if task.MaxRetries < 0 {
		return false
	}
	if task.Project <= 0 {
		return false
	}
	if len(task.Recipe) <= 0 {
		return false
	}

	return true
}

func (api *WebAPI) TaskGet(r *Request) {


}