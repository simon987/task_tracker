package api

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"src/task_tracker/storage"
	"strconv"
)

type CreateTaskRequest struct {
	Project       int64  `json:"project"`
	MaxRetries    int64  `json:"max_retries"`
	Recipe        string `json:"recipe"`
	Priority      int64  `json:"priority"`
	MaxAssignTime int64  `json:"max_assign_time"`
}

type ReleaseTaskRequest struct {
	TaskId   int64      `json:"task_id"`
	Success  bool       `json:"success"`
	WorkerId *uuid.UUID `json:"worker_id"`
}

type ReleaseTaskResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
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
			MaxRetries:    createReq.MaxRetries,
			Recipe:        createReq.Recipe,
			Priority:      createReq.Priority,
			AssignTime:    0,
			MaxAssignTime: createReq.MaxAssignTime,
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

	if task == nil {

		r.OkJson(GetTaskResponse{
			Ok:      false,
			Message: "No task available",
		})

	} else {

		r.OkJson(GetTaskResponse{
			Ok:   true,
			Task: task,
		})
	}

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

func (api *WebAPI) TaskRelease(r *Request) {

	req := ReleaseTaskRequest{}
	if r.GetJson(req) {

		res := api.Database.ReleaseTask(req.TaskId, req.WorkerId, req.Success)

		response := ReleaseTaskResponse{
			Ok: res,
		}

		if !res {
			response.Message = "Could not find a task with the specified Id assigned to this workerId"

			logrus.WithFields(logrus.Fields{
				"releaseTaskRequest": req,
				"taskUpdated":        res,
			}).Warn("Release task: NOT FOUND")
		} else {

			logrus.WithFields(logrus.Fields{
				"releaseTaskRequest": req,
				"taskUpdated":        res,
			}).Trace("Release task")
		}

		r.OkJson(response)
	}
}
