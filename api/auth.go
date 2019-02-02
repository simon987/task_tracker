package api

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/kataras/go-sessions"
	"github.com/simon987/task_tracker/storage"
)

type LoginRequest struct {
	Username []byte
	Password []byte
}

type LoginResponse struct {
	Ok      bool
	Message string
	Manager *storage.Manager
}

func (api *WebAPI) Login(r *Request) {

	req := &LoginRequest{}
	err := json.Unmarshal(r.Ctx.Request.Body(), req)

	if err != nil {
		r.Json(LoginResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}

	manager, err := api.Database.ValidateCredentials(req.Username, req.Password)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"username": string(manager.Username),
		}).Warning("Login attempt")

		r.Json(LoginResponse{
			Ok:      false,
			Message: "Invalid username/password",
		}, 403)
		return
	}

	sess := sessions.StartFasthttp(r.Ctx)
	sess.Set("manager", manager)

	logrus.WithFields(logrus.Fields{
		"username": string(manager.Username),
	}).Info("Logged in")
}
