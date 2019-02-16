package api

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/simon987/task_tracker/storage"
	"math/rand"
	"strconv"
	"time"
)

type UpdateWorkerRequest struct {
	Alias string `json:"alias"`
}

type UpdateWorkerResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}

type CreateWorkerRequest struct {
	Alias string `json:"alias"`
}

type CreateWorkerResponse struct {
	Ok      bool            `json:"ok"`
	Message string          `json:"message,omitempty"`
	Worker  *storage.Worker `json:"worker,omitempty"`
}

type GetWorkerResponse struct {
	Ok      bool            `json:"ok"`
	Message string          `json:"message,omitempty"`
	Worker  *storage.Worker `json:"worker,omitempty"`
}

type GetAllWorkerStatsResponse struct {
	Ok      bool                   `json:"ok"`
	Message string                 `json:"message,omitempty"`
	Stats   *[]storage.WorkerStats `json:"stats"`
}

type WorkerAccessRequest struct {
	Assign  bool  `json:"assign"`
	Submit  bool  `json:"submit"`
	Project int64 `json:"project"`
}

func (w *WorkerAccessRequest) isValid() bool {
	if !w.Assign && !w.Submit {
		return false
	}
	return true
}

func (api *WebAPI) WorkerCreate(r *Request) {

	workerReq := &CreateWorkerRequest{}
	err := json.Unmarshal(r.Ctx.Request.Body(), workerReq)
	if err != nil {
		return
	}

	if !canCreateWorker(r, workerReq) {

		logrus.WithFields(logrus.Fields{
			"createWorkerRequest": workerReq,
		}).Warn("Failed CreateWorkerRequest")

		r.Json(CreateWorkerResponse{
			Ok:      false,
			Message: "You are now allowed to create a worker",
		}, 403)
		return
	}

	worker, err := api.workerCreate(workerReq)
	if err != nil {
		handleErr(err, r)
	} else {
		r.OkJson(CreateWorkerResponse{
			Ok:     true,
			Worker: worker,
		})
	}
}

func (api *WebAPI) WorkerGet(r *Request) {

	id, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"id": id,
		}).Warn("Invalid worker id")

		r.Json(GetWorkerResponse{
			Ok:      false,
			Message: err.Error(),
		}, 400)
		return
	} else if id <= 0 {
		logrus.WithFields(logrus.Fields{
			"id": id,
		}).Warn("Invalid worker id")

		r.Json(GetWorkerResponse{
			Ok:      false,
			Message: "Invalid worker id",
		}, 400)
		return
	}

	worker := api.Database.GetWorker(id)

	if worker != nil {

		worker.Secret = nil

		r.OkJson(GetWorkerResponse{
			Ok:     true,
			Worker: worker,
		})
	} else {
		r.Json(GetWorkerResponse{
			Ok:      false,
			Message: "Worker not found",
		}, 404)
	}
}

func (api *WebAPI) WorkerUpdate(r *Request) {

	worker, err := api.validateSignature(r)
	if err != nil {
		r.Json(GetTaskResponse{
			Ok:      false,
			Message: err.Error(),
		}, 401)
		return
	}

	req := &UpdateWorkerRequest{}
	err = json.Unmarshal(r.Ctx.Request.Body(), req)
	if err != nil {
		r.Json(GetTaskResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}
	worker.Alias = req.Alias

	ok := api.Database.UpdateWorker(worker)

	if ok {
		r.OkJson(UpdateWorkerResponse{
			Ok: true,
		})
	} else {
		r.OkJson(UpdateWorkerResponse{
			Ok:      false,
			Message: "Could not update worker",
		})
	}
}

func (api *WebAPI) GetAllWorkerStats(r *Request) {

	stats := api.Database.GetAllWorkerStats()

	r.OkJson(GetAllWorkerStatsResponse{
		Ok:    true,
		Stats: stats,
	})
}

func (api *WebAPI) workerCreate(request *CreateWorkerRequest) (*storage.Worker, error) {

	if request.Alias == "" {
		request.Alias = "default_alias"
	}

	worker := storage.Worker{
		Created: time.Now().Unix(),
		Secret:  makeSecret(),
		Alias:   request.Alias,
	}

	api.Database.SaveWorker(&worker)
	return &worker, nil
}

func canCreateWorker(r *Request, cwr *CreateWorkerRequest) bool {

	if cwr.Alias == "unassigned" {
		//Reserved alias
		return false
	}

	return true
}

func makeSecret() []byte {

	secret := make([]byte, 32)
	for i := 0; i < 32; i++ {
		secret[i] = byte(rand.Int31())
	}

	return secret
}
