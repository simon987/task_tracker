package api

import (
	"encoding/json"
	"github.com/simon987/task_tracker/storage"
	"github.com/sirupsen/logrus"
	"strconv"
)

func (api *WebAPI) Login(r *Request) {

	req := &LoginRequest{}
	err := json.Unmarshal(r.Ctx.Request.Body(), req)

	if err != nil {
		r.Json(JsonResponse{
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

		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid username/password",
		}, 403)
		return
	}

	sess, _ := api.Session.Get(r.Ctx)
	sess.Set("manager", manager)
	api.Session.Save(r.Ctx, sess)

	r.OkJson(JsonResponse{
		Content: LoginResponse{
			Manager: manager,
		},
		Ok: true,
	})

	logrus.WithFields(logrus.Fields{
		"username": string(manager.Username),
	}).Info("Logged in")
}

func (api *WebAPI) Logout(r *Request) {

	sess, _ := api.Session.Get(r.Ctx)
	sess.Flush()
	r.Ctx.Response.SetStatusCode(204)
}

func (api *WebAPI) Register(r *Request) {

	req := &RegisterRequest{}
	err := json.Unmarshal(r.Ctx.Request.Body(), req)

	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}

	if !req.isValid() {
		r.Json(JsonResponse{
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

		r.Json(JsonResponse{
			Ok:      false,
			Message: err.Error(),
		}, 400)
		return
	}

	sess, _ := api.Session.Get(r.Ctx)
	sess.Set("manager", manager)
	api.Session.Save(r.Ctx, sess)

	r.OkJson(JsonResponse{
		Ok: true,
	})

	logrus.WithFields(logrus.Fields{
		"username": string(manager.Username),
	}).Info("Registered")
}

func (api *WebAPI) GetAccountDetails(r *Request) {

	sess, _ := api.Session.Get(r.Ctx)
	manager := sess.Get("manager")

	logrus.WithFields(logrus.Fields{
		"manager": manager,
	}).Trace("Account details request")

	if manager == nil {
		r.OkJson(JsonResponse{
			Ok: false,
			Content: GetAccountDetailsResponse{
				LoggedIn: false,
			},
		})
	} else {
		r.OkJson(JsonResponse{
			Ok: true,
			Content: GetAccountDetailsResponse{
				LoggedIn: true,
				Manager:  manager.(*storage.Manager),
			},
		})
	}
}

func (api *WebAPI) GetManagerList(r *Request) {

	sess, _ := api.Session.Get(r.Ctx)
	manager := sess.Get("manager")

	if manager == nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Unauthorized",
		}, 401)
		return
	}

	managers := api.Database.GetManagerList()

	r.OkJson(JsonResponse{
		Ok: true,
		Content: GetManagerListResponse{
			Managers: managers,
		},
	})
}

func (api *WebAPI) GetManagerListWithRoleOn(r *Request) {

	pid, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	if err != nil || pid <= 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid project id",
		}, 400)
		return
	}

	sess, _ := api.Session.Get(r.Ctx)
	manager := sess.Get("manager")

	if manager == nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Unauthorized",
		}, 401)
		return
	}

	managers := api.Database.GetManagerListWithRoleOn(pid)

	r.OkJson(JsonResponse{
		Ok: true,
		Content: GetManagerListWithRoleOnResponse{
			Managers: managers,
		},
	})
}

func (api *WebAPI) PromoteManager(r *Request) {

	id, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	if err != nil || id <= 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid manager id",
		}, 400)
		return
	}

	sess, _ := api.Session.Get(r.Ctx)
	manager := sess.Get("manager")

	if !manager.(*storage.Manager).WebsiteAdmin || manager.(*storage.Manager).Id == id {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Unauthorized",
		}, 401)
		return
	}

	if !manager.(*storage.Manager).WebsiteAdmin {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Unauthorized",
		}, 403)
		return
	}

	api.Database.UpdateManager(&storage.Manager{
		Id:           id,
		WebsiteAdmin: true,
	})

	r.Ctx.Response.SetStatusCode(204)
}

func (api *WebAPI) DemoteManager(r *Request) {

	id, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	if err != nil || id <= 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid manager id",
		}, 400)
		return
	}

	sess, _ := api.Session.Get(r.Ctx)
	manager := sess.Get("manager")

	if manager == nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Unauthorized",
		}, 401)
		return
	}

	if !manager.(*storage.Manager).WebsiteAdmin || manager.(*storage.Manager).Id == id {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Unauthorized",
		}, 403)
		return
	}

	api.Database.UpdateManager(&storage.Manager{
		Id:           id,
		WebsiteAdmin: false,
	})

	r.Ctx.Response.SetStatusCode(204)
}
