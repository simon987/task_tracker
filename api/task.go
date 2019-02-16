package api

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/dchest/siphash"
	"github.com/simon987/task_tracker/storage"
	"strconv"
)

type CreateTaskRequest struct {
	Project           int64  `json:"projectChange"`
	MaxRetries        int64  `json:"max_retries"`
	Recipe            string `json:"recipe"`
	Priority          int64  `json:"priority"`
	MaxAssignTime     int64  `json:"max_assign_time"`
	Hash64            int64  `json:"hash_u64"`
	UniqueString      string `json:"unique_string"`
	VerificationCount int64  `json:"verification_count"`
}

type ReleaseTaskRequest struct {
	TaskId       int64              `json:"task_id"`
	Result       storage.TaskResult `json:"result"`
	Verification int64              `json:"verification"`
}

type ReleaseTaskResponse struct {
	Ok      bool   `json:"ok"`
	Updated bool   `json:"updated"`
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

	createReq := &CreateTaskRequest{}
	err := json.Unmarshal(r.Ctx.Request.Body(), createReq)
	if err != nil {
		r.Json(CreateProjectResponse{
			Ok:      false,
			Message: "Could not parse request",
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

	if createReq.IsValid() && isTaskValid(task) {

		if createReq.UniqueString != "" {
			//TODO: Load key from config
			createReq.Hash64 = int64(siphash.Hash(1, 2, []byte(createReq.UniqueString)))
		}

		err := api.Database.SaveTask(task, createReq.Project, createReq.Hash64)

		if err != nil {
			r.Json(CreateTaskResponse{
				Ok:      false,
				Message: err.Error(),
			}, 400)
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

func (req *CreateTaskRequest) IsValid() bool {
	return req.Hash64 == 0 || req.UniqueString == ""
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

	worker, err := api.validateSignature(r)
	if err != nil {
		r.Json(GetTaskResponse{
			Ok:      false,
			Message: err.Error(),
		}, 403)
		return
	}

	project, err := strconv.ParseInt(r.Ctx.UserValue("projectChange").(string), 10, 64)
	handleErr(err, r)
	task := api.Database.GetTaskFromProject(worker, project)

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

	worker, err := api.validateSignature(r)
	if err != nil {
		r.Json(GetTaskResponse{
			Ok:      false,
			Message: err.Error(),
		}, 403)
		return
	}

	task := api.Database.GetTask(worker)
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

func (api WebAPI) validateSignature(r *Request) (*storage.Worker, error) {

	widStr := string(r.Ctx.Request.Header.Peek("X-Worker-Id"))
	signature := r.Ctx.Request.Header.Peek("X-Signature")

	if widStr == "" {
		return nil, errors.New("worker id not specified")
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

	var body []byte
	if r.Ctx.Request.Header.IsGet() {
		body = r.Ctx.Request.RequestURI()
	} else {
		body = r.Ctx.Request.Body()
	}

	mac := hmac.New(crypto.SHA256.New, worker.Secret)
	mac.Write(body)

	expectedMac := make([]byte, 64)
	hex.Encode(expectedMac, mac.Sum(nil))
	matches := bytes.Compare(expectedMac, signature) == 0

	logrus.WithFields(logrus.Fields{
		"expected":  string(expectedMac),
		"signature": string(signature),
		"matches":   matches,
	}).Trace("Validating Worker signature")

	if !matches {
		return nil, errors.New("invalid signature")
	}

	return worker, nil
}

func (api *WebAPI) TaskRelease(r *Request) {

	worker, err := api.validateSignature(r)
	if err != nil {
		r.Json(GetTaskResponse{
			Ok:      false,
			Message: err.Error(),
		}, 403)
		return
	}

	req := &ReleaseTaskRequest{}
	err = json.Unmarshal(r.Ctx.Request.Body(), req)
	if err != nil {
		r.Json(CreateProjectResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
	}
	res := api.Database.ReleaseTask(req.TaskId, worker.Id, req.Result, req.Verification)

	response := ReleaseTaskResponse{
		Updated: res,
		Ok:      true,
	}

	if !res {
		response.Message = "Task was not marked as closed"
	} else {

		logrus.WithFields(logrus.Fields{
			"releaseTaskRequest": req,
			"taskUpdated":        res,
		}).Trace("Release task")
	}

	r.OkJson(response)
}
