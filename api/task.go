package api

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"src/task_tracker/storage"
	"strconv"
)

type CreateTaskRequest struct {
	Project    int64  `json:"project"`
	MaxRetries int64  `json:"max_retries"`
	Recipe     string `json:"recipe"`
	Priority   int64  `json:"priority"`
}

type CreateTaskResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}

type GetTaskResponse struct {
	Ok      bool          `json:"ok"`
	Message string        `json:"message,omitempty"`
	Task    *storage.Task `json:"task,omitempty"`
}

func (api *WebAPI) TaskCreate(r *Request) {

	var createReq CreateTaskRequest
	if r.GetJson(&createReq) {

		task := &storage.Task{
			MaxRetries: createReq.MaxRetries,
			Recipe:     createReq.Recipe,
			Priority:   createReq.Priority,
		}

		if isTaskValid(task) {
			err := api.Database.SaveTask(task, createReq.Project)

			if err != nil {
				r.Json(CreateTaskResponse{
					Ok:      false,
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
				Ok:      false,
				Message: "Invalid task",
			}, 400)
		}
	}
}

func isTaskValid(task *storage.Task) bool {
	if task.MaxRetries < 0 {
		return false
	}
	if len(task.Recipe) <= 0 {
		return false
	}

	return true
}

func (api *WebAPI) TaskGetFromProject(r *Request) {

	worker, err := api.workerFromQueryArgs(r)
	if err != nil {
		r.Json(GetTaskResponse{
			Ok:      false,
			Message: err.Error(),
		}, 403)
		return
	}

	project, err := strconv.Atoi(r.Ctx.UserValue("project").(string))
	handleErr(err, r)
	task := api.Database.GetTaskFromProject(worker, int64(project))

	r.OkJson(GetTaskResponse{
		Ok:   true,
		Task: task,
	})
}

func (api *WebAPI) TaskGet(r *Request) {

	worker, err := api.workerFromQueryArgs(r)
	if err != nil {
		r.Json(GetTaskResponse{
			Ok:      false,
			Message: err.Error(),
		}, 403)
		return
	}

	task := api.Database.GetTask(worker)

	r.OkJson(GetTaskResponse{
		Ok:   true,
		Task: task,
	})
}

func (api WebAPI) workerFromQueryArgs(r *Request) (*storage.Worker, error) {

	widStr := string(r.Ctx.QueryArgs().Peek("wid"))
	wid, err := uuid.Parse(widStr)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"wid": widStr,
		}).Warn("Can't parse wid")

		return nil, err
	}

	worker := api.Database.GetWorker(wid)

	if worker == nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"wid": widStr,
		}).Warn("Can't parse wid")

		return nil, errors.New("worker id does not match any valid worker")
	}

	return worker, nil
}
