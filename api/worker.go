package api

import (
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"math/rand"
	"src/task_tracker/storage"
	"time"
)

type CreateWorkerRequest struct {
	Alias string `json:"alias"`
}

type UpdateWorkerRequest struct {
	Alias string `json:"alias"`
}

type UpdateWorkerResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
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

type WorkerAccessRequest struct {
	WorkerId  *uuid.UUID `json:"worker_id"`
	ProjectId int64      `json:"project_id"`
}

type WorkerAccessResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func (api *WebAPI) WorkerCreate(r *Request) {

	workerReq := &CreateWorkerRequest{}
	if !r.GetJson(workerReq) {
		return
	}

	identity := getIdentity(r)

	if !canCreateWorker(r, workerReq, identity) {

		logrus.WithFields(logrus.Fields{
			"identity":            identity,
			"createWorkerRequest": workerReq,
		}).Warn("Failed CreateWorkerRequest")

		r.Json(CreateWorkerResponse{
			Ok:      false,
			Message: "You are now allowed to create a worker",
		}, 403)
		return
	}

	worker, err := api.workerCreate(workerReq, getIdentity(r))
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

	id, err := uuid.Parse(r.Ctx.UserValue("id").(string))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"id": id,
		}).Warn("Invalid UUID")

		r.Json(GetWorkerResponse{
			Ok:      false,
			Message: err.Error(),
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

func (api *WebAPI) WorkerGrantAccess(r *Request) {

	req := &WorkerAccessRequest{}
	if r.GetJson(req) {

		ok := api.Database.GrantAccess(req.WorkerId, req.ProjectId)

		if ok {
			r.OkJson(WorkerAccessResponse{
				Ok: true,
			})
		} else {
			r.OkJson(WorkerAccessResponse{
				Ok:      false,
				Message: "Worker already has access to this project",
			})
		}
	}
}

func (api *WebAPI) WorkerRemoveAccess(r *Request) {

	req := &WorkerAccessRequest{}
	if r.GetJson(req) {

		ok := api.Database.RemoveAccess(req.WorkerId, req.ProjectId)

		if ok {
			r.OkJson(WorkerAccessResponse{
				Ok: true,
			})
		} else {
			r.OkJson(WorkerAccessResponse{
				Ok:      false,
				Message: "Worker did not have access to this project",
			})
		}
	}
}

func (api *WebAPI) WorkerUpdate(r *Request) {

	worker, err := api.validateSignature(r)
	if err != nil {
		r.Json(GetTaskResponse{
			Ok:      false,
			Message: err.Error(),
		}, 403)
		return
	}

	req := &UpdateWorkerRequest{}
	if r.GetJson(req) {

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
}

func (api *WebAPI) workerCreate(request *CreateWorkerRequest, identity *storage.Identity) (*storage.Worker, error) {

	if request.Alias == "" {
		request.Alias = "default_alias"
	}

	worker := storage.Worker{
		Id:       uuid.New(),
		Created:  time.Now().Unix(),
		Identity: identity,
		Secret:   makeSecret(),
		Alias:    request.Alias,
	}

	api.Database.SaveWorker(&worker)
	return &worker, nil
}

func canCreateWorker(r *Request, cwr *CreateWorkerRequest, identity *storage.Identity) bool {

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

func getIdentity(r *Request) *storage.Identity {

	identity := storage.Identity{
		RemoteAddr: r.Ctx.RemoteAddr().String(),
		UserAgent:  string(r.Ctx.UserAgent()),
	}

	return &identity
}
