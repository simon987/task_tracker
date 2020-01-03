package api

import (
	"encoding/json"
	"github.com/simon987/task_tracker/storage"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strconv"
	"time"
)

func (api *WebAPI) CreateWorker(r *Request) {

	workerReq := &CreateWorkerRequest{}
	err := json.Unmarshal(r.Ctx.Request.Body(), workerReq)
	if err != nil {
		return
	}

	if !workerReq.isValid() {

		logrus.WithFields(logrus.Fields{
			"createWorkerRequest": workerReq,
		}).Warn("Failed CreateWorkerRequest")

		r.Json(JsonResponse{
			Ok:      false,
			Message: "You are now allowed to create a worker",
		}, 403)
		return
	}

	worker, err := api.workerCreate(workerReq)
	if err != nil {
		handleErr(err, r)
	} else {
		r.OkJson(JsonResponse{
			Ok: true,
			Content: CreateWorkerResponse{
				Worker: worker,
			},
		})
	}
}

func (api *WebAPI) GetWorker(r *Request) {

	id, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"id": id,
		}).Warn("Invalid worker id")

		r.Json(JsonResponse{
			Ok:      false,
			Message: err.Error(),
		}, 400)
		return
	} else if id <= 0 {
		logrus.WithFields(logrus.Fields{
			"id": id,
		}).Warn("Invalid worker id")

		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid worker id",
		}, 400)
		return
	}

	worker := api.Database.GetWorker(id)

	if worker != nil {

		sess, _ := api.Session.Get(r.Ctx)
		manager := sess.Get("manager")

		var secret []byte = nil
		if manager != nil && manager.(*storage.Manager).WebsiteAdmin {
			secret = worker.Secret
		}

		r.OkJson(JsonResponse{
			Ok: true,
			Content: GetWorkerResponse{
				Worker: &storage.Worker{
					Alias:   worker.Alias,
					Id:      worker.Id,
					Created: worker.Created,
					Paused:  worker.Paused,
					Secret:  secret,
				},
			},
		})
	} else {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Worker not found",
		}, 404)
	}
}

func (api *WebAPI) UpdateWorker(r *Request) {

	worker, err := api.validateSecret(r)
	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: err.Error(),
		}, 401)
		return
	}

	req := &UpdateWorkerRequest{}
	err = json.Unmarshal(r.Ctx.Request.Body(), req)
	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}
	worker.Alias = req.Alias

	ok := api.Database.UpdateWorker(worker)

	if ok {
		r.OkJson(JsonResponse{
			Ok: true,
		})
	} else {
		r.OkJson(JsonResponse{
			Ok:      false,
			Message: "Could not update worker",
		})
	}
}

func (api *WebAPI) WorkerSetPaused(r *Request) {

	sess, _ := api.Session.Get(r.Ctx)
	manager := sess.Get("manager")

	if manager == nil || !manager.(*storage.Manager).WebsiteAdmin {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Unauthorized",
		}, 403)
		return
	}

	req := &WorkerSetPausedRequest{}
	err := json.Unmarshal(r.Ctx.Request.Body(), req)
	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}

	worker := api.Database.GetWorker(req.Worker)
	if worker == nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid worker",
		}, 400)
		return
	}

	worker.Paused = req.Paused

	ok := api.Database.UpdateWorker(worker)

	if ok {
		r.OkJson(JsonResponse{
			Ok: true,
		})
	} else {
		r.OkJson(JsonResponse{
			Ok:      false,
			Message: "Could not update worker",
		})
	}
}

func (api *WebAPI) GetAllWorkerStats(r *Request) {

	stats := api.Database.GetAllWorkerStats()

	r.OkJson(JsonResponse{
		Ok: true,
		Content: GetAllWorkerStatsResponse{
			Stats: stats,
		},
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

func makeSecret() []byte {

	secret := make([]byte, 32)
	for i := 0; i < 32; i++ {
		secret[i] = byte(rand.Int31())
	}

	return secret
}
