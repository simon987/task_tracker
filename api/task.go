package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/dchest/siphash"
	"github.com/simon987/task_tracker/storage"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

func (api *WebAPI) SubmitTask(r *Request) {

	worker, err := api.validateSecret(r)
	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: err.Error(),
		}, 401)
		return
	}

	createReq := &SubmitTaskRequest{}
	err = json.Unmarshal(r.Ctx.Request.Body(), createReq)
	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}

	if !createReq.IsValid() {
		logrus.WithFields(logrus.Fields{
			"req": createReq,
		}).Warn("Invalid task")
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid task",
		}, 400)
		return
	}

	task := &storage.Task{
		MaxRetries:        createReq.MaxRetries,
		Recipe:            createReq.Recipe,
		Priority:          createReq.Priority,
		AssignTime:        0,
		MaxAssignTime:     createReq.MaxAssignTime,
		VerificationCount: createReq.VerificationCount,
	}

	reservation := api.ReserveSubmit(createReq.Project, 1)
	if reservation == nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Project not found",
		}, 404)
		return
	}
	delay := reservation.DelayFrom(time.Now()).Seconds()
	if delay > 0 {
		r.Json(JsonResponse{
			Ok:             false,
			Message:        "Too many requests",
			RateLimitDelay: delay,
		}, 429)
		reservation.Cancel()
		return
	}

	if createReq.UniqueString != "" {
		createReq.Hash64 = int64(siphash.Hash(1, 2, []byte(createReq.UniqueString)))
	}

	err = api.Database.SaveTask(task, createReq.Project, createReq.Hash64, worker.Id)

	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: err.Error(),
		}, 400)
		reservation.Cancel()
		return
	}

	r.OkJson(JsonResponse{
		Ok: true,
	})
}

func (api *WebAPI) BulkSubmitTask(r *Request) {

	worker, err := api.validateSecret(r)
	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: err.Error(),
		}, 401)
		return
	}

	createReq := &BulkSubmitTaskRequest{}
	err = json.Unmarshal(r.Ctx.Request.Body(), createReq)
	if err != nil || createReq.Requests == nil || len(createReq.Requests) == 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}
	if !createReq.IsValid() {
		logrus.WithFields(logrus.Fields{
			"req": createReq,
		}).Warn("Invalid request")
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid request",
		}, 400)
		return
	}

	saveRequests := make([]storage.SaveTaskRequest, len(createReq.Requests))
	projectId := createReq.Requests[0].Project
	for i, req := range createReq.Requests {

		if req.Project != projectId {
			r.Json(JsonResponse{
				Ok:      false,
				Message: "All the tasks in a bulk submit must be of the same project",
			}, 400)
			return
		}

		if req.UniqueString != "" {
			req.Hash64 = int64(siphash.Hash(1, 2, []byte(req.UniqueString)))
		}

		saveRequests[i] = storage.SaveTaskRequest{
			Task: &storage.Task{
				MaxRetries:        req.MaxRetries,
				Recipe:            req.Recipe,
				Priority:          req.Priority,
				AssignTime:        0,
				MaxAssignTime:     req.MaxAssignTime,
				VerificationCount: req.VerificationCount,
			},
			Project:  projectId,
			WorkerId: worker.Id,
			Hash64:   req.Hash64,
		}
	}

	reservation := api.ReserveSubmit(projectId, len(saveRequests))
	if reservation == nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Project not found",
		}, 404)
		return
	}
	delay := reservation.DelayFrom(time.Now()).Seconds()
	if delay > 0 {
		r.Json(JsonResponse{
			Ok:             false,
			Message:        "Too many requests",
			RateLimitDelay: delay,
		}, 429)
		reservation.Cancel()
		return
	}

	saveErrors := api.Database.BulkSaveTask(saveRequests)

	if saveErrors == nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Fatal error during bulk insert, see server logs",
		}, 400)
		reservation.Cancel()
		return
	}

	r.OkJson(JsonResponse{
		Ok: true,
	})
}

func (api *WebAPI) GetTaskFromProject(r *Request) {

	worker, err := api.validateSecret(r)
	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: err.Error(),
		}, 403)
		return
	}

	if worker.Paused {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "A manager has paused you",
		}, 400)
		return
	}

	project, err := strconv.ParseInt(r.Ctx.UserValue("project").(string), 10, 64)
	if err != nil || project <= 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid project id",
		}, 400)
		return
	}

	reservation := api.ReserveAssign(project)
	if reservation == nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Project not found",
		}, 404)
		return
	}
	delay := reservation.DelayFrom(time.Now()).Seconds()
	if delay > 0 {
		r.Json(JsonResponse{
			Ok:             false,
			Message:        "Too many requests",
			RateLimitDelay: delay,
		}, 429)
		reservation.Cancel()
		return
	}

	task := api.Database.GetTaskFromProject(worker, project)

	if task == nil {
		r.OkJson(JsonResponse{
			Ok:      false,
			Message: "No task available",
		})
		reservation.CancelAt(time.Now())
		return
	}

	r.OkJson(JsonResponse{
		Ok: true,
		Content: GetTaskResponse{
			Task: task,
		},
	})
}

func (api *WebAPI) validateSecret(r *Request) (*storage.Worker, error) {

	widStr := string(r.Ctx.Request.Header.Peek("X-Worker-Id"))
	secretHeader := r.Ctx.Request.Header.Peek("X-Secret")

	if widStr == "" {
		return nil, errors.New("worker id not specified")
	}
	if bytes.Equal(secretHeader, []byte("")) {
		return nil, errors.New("secret is not specified")
	}

	wid, err := strconv.ParseInt(widStr, 10, 64)
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
		}).Warn("Worker id does not match any valid worker")

		return nil, errors.New("worker id does not match any valid worker")
	}

	secret := make([]byte, base64.StdEncoding.EncodedLen(len(worker.Secret)))
	secretLen, _ := base64.StdEncoding.Decode(secret, secretHeader)
	matches := bytes.Equal(worker.Secret, secret[:secretLen])

	logrus.WithFields(logrus.Fields{
		"expected": string(worker.Secret),
		"header":   string(secretHeader),
		"matches":  matches,
	}).Trace("Validating Worker secret")

	if !matches {
		return nil, errors.New("invalid secret")
	}

	return worker, nil
}

func (api *WebAPI) ReleaseTask(r *Request) {

	worker, err := api.validateSecret(r)
	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: err.Error(),
		}, 403)
		return
	}

	req := &ReleaseTaskRequest{}
	err = json.Unmarshal(r.Ctx.Request.Body(), req)
	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}

	if !req.IsValid() {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid request",
		}, 400)
		return
	}

	res := api.Database.ReleaseTask(req.TaskId, worker.Id, req.Result, req.Verification)

	response := JsonResponse{
		Ok: true,
		Content: ReleaseTaskResponse{
			Updated: res,
		},
	}

	if !res {
		response.Message = "Task was not marked as closed"
	}

	logrus.WithFields(logrus.Fields{
		"releaseTaskRequest": req,
		"taskUpdated":        res,
	}).Trace("Release task")

	r.OkJson(response)
}
