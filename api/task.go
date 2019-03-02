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
	"math"
	"strconv"
	"time"
)

func (api *WebAPI) SubmitTask(r *Request) {

	worker, err := api.validateSignature(r)
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
	task := &storage.Task{
		MaxRetries:        createReq.MaxRetries,
		Recipe:            createReq.Recipe,
		Priority:          createReq.Priority,
		AssignTime:        0,
		MaxAssignTime:     createReq.MaxAssignTime,
		VerificationCount: createReq.VerificationCount,
	}

	if !createReq.IsValid() {
		logrus.WithFields(logrus.Fields{
			"task": task,
		}).Warn("Invalid task")
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid task",
		}, 400)
		return
	}

	reservation := api.ReserveSubmit(createReq.Project)
	delay := reservation.DelayFrom(time.Now()).Seconds()
	if delay > 0 {
		r.Json(JsonResponse{
			Ok:             false,
			Message:        "Too many requests",
			RateLimitDelay: delay,
		}, 429)
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

func (api *WebAPI) GetTaskFromProject(r *Request) {

	worker, err := api.validateSignature(r)
	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: err.Error(),
		}, 403)
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
	delay := reservation.DelayFrom(time.Now()).Seconds()
	if delay > 0 {
		r.Json(JsonResponse{
			Ok:             false,
			Message:        "Too many requests",
			RateLimitDelay: delay,
		}, 429)
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

func (api WebAPI) validateSignature(r *Request) (*storage.Worker, error) {

	widStr := string(r.Ctx.Request.Header.Peek("X-Worker-Id"))
	timeStampStr := string(r.Ctx.Request.Header.Peek("Timestamp"))
	signature := r.Ctx.Request.Header.Peek("X-Signature")

	if widStr == "" {
		return nil, errors.New("worker id not specified")
	}
	if timeStampStr == "" {
		return nil, errors.New("date is not specified")
	}

	timestamp, err := time.Parse(time.RFC1123, timeStampStr)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"date": timeStampStr,
		}).Warn("Can't parse Timestamp")

		return nil, err
	}

	if math.Abs(float64(timestamp.Unix()-time.Now().Unix())) > 60 {
		logrus.WithError(err).WithFields(logrus.Fields{
			"date": timeStampStr,
		}).Warn("Invalid Timestamp")

		return nil, errors.New("invalid Timestamp")
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
	mac.Write([]byte(timeStampStr))

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

func (api *WebAPI) ReleaseTask(r *Request) {

	worker, err := api.validateSignature(r)
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
