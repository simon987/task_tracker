package api

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/simon987/task_tracker/storage"
)

const MaxUsernameLength = 16

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Ok      bool             `json:"ok"`
	Message string           `json:"message,omitempty"`
	Manager *storage.Manager `json:"manager"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AccountDetails struct {
	LoggedIn bool             `json:"logged_in"`
	Manager  *storage.Manager `json:"manager,omitempty"`
}

func (r *RegisterRequest) isValid() bool {
	return len(r.Username) <= MaxUsernameLength
}

type RegisterResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
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

	manager, err := api.Database.ValidateCredentials([]byte(req.Username), []byte(req.Password))
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"username": req.Username,
		}).Warning("Login attempt")

		r.Json(LoginResponse{
			Ok:      false,
			Message: "Invalid username/password",
		}, 403)
		return
	}

	sess := api.Session.StartFasthttp(r.Ctx)
	sess.Set("manager", manager)

	logrus.Debug("SET")
	logrus.Debug(sess.ID())
	logrus.Debug(manager)

	r.OkJson(LoginResponse{
		Manager: manager,
		Ok:      true,
	})

	logrus.WithFields(logrus.Fields{
		"username": string(manager.Username),
	}).Info("Logged in")
}

func (api *WebAPI) Register(r *Request) {

	req := &RegisterRequest{}
	err := json.Unmarshal(r.Ctx.Request.Body(), req)

	if err != nil {
		r.Json(LoginResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}

	if !req.isValid() {
		r.Json(LoginResponse{
			Ok:      false,
			Message: "Invalid register request",
		}, 400)
		return
	}

	manager := &storage.Manager{
		Username: string(req.Username),
	}

	err = api.Database.SaveManager(manager, []byte(req.Password))
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"username": string(manager.Username),
		}).Warning("Register attempt")

		r.Json(LoginResponse{
			Ok:      false,
			Message: err.Error(),
		}, 400)
		return
	}

	sess := api.Session.StartFasthttp(r.Ctx)
	sess.Set("manager", manager)

	r.OkJson(RegisterResponse{
		Ok: true,
	})

	logrus.WithFields(logrus.Fields{
		"username": string(manager.Username),
	}).Info("Registered")
}

func (api *WebAPI) AccountDetails(r *Request) {

	sess := api.Session.StartFasthttp(r.Ctx)
	manager := sess.Get("manager")
	logrus.Debug("GET")
	logrus.Debug(sess.ID())

	if manager == nil {
		r.OkJson(AccountDetails{
			LoggedIn: false,
		})
	} else {
		r.OkJson(AccountDetails{
			LoggedIn: true,
			Manager:  manager.(*storage.Manager),
		})
	}
}
