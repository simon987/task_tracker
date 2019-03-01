package api

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/google/uuid"
	"github.com/simon987/task_tracker/storage"
	"golang.org/x/time/rate"
	"strconv"
)

func (api *WebAPI) GetProject(r *Request) {

	id, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	if err != nil || id <= 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid worker id",
		}, 400)
		return
	}

	sess := api.Session.StartFasthttp(r.Ctx)
	manager := sess.Get("manager")

	project := api.Database.GetProject(id)

	if project == nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Project not found",
		}, 404)
		return
	}

	if !isProjectReadAuthorized(project, manager, api.Database) {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Unauthorized",
		}, 403)
		return
	}

	r.OkJson(JsonResponse{
		Ok: true,
		Content: GetProjectResponse{
			Project: project,
		},
	})
}

func (api *WebAPI) CreateProject(r *Request) {

	sess := api.Session.StartFasthttp(r.Ctx)
	manager := sess.Get("manager")

	createReq := &CreateProjectRequest{}
	err := json.Unmarshal(r.Ctx.Request.Body(), createReq)
	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}

	if createReq.AssignRate == 0 {
		createReq.AssignRate = rate.Inf
	}
	if createReq.SubmitRate == 0 {
		createReq.SubmitRate = rate.Inf
	}
	project := &storage.Project{
		Name:       createReq.Name,
		Version:    createReq.Version,
		CloneUrl:   createReq.CloneUrl,
		GitRepo:    createReq.GitRepo,
		Priority:   createReq.Priority,
		Motd:       createReq.Motd,
		Public:     createReq.Public,
		Hidden:     createReq.Hidden,
		Chain:      createReq.Chain,
		AssignRate: createReq.AssignRate,
		SubmitRate: createReq.SubmitRate,
	}

	if !createReq.isValid() {
		logrus.WithFields(logrus.Fields{
			"project": project,
		}).Warn("Invalid project")

		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid project",
		}, 400)
		return
	}

	if !api.isProjectCreationAuthorized(project, manager) {
		logrus.WithFields(logrus.Fields{
			"project": project,
		}).Warn("Unauthorized project creation")

		r.Json(JsonResponse{
			Ok:      false,
			Message: "You are not permitted to create a project with this configuration",
		}, 400)
		return
	}

	webhookSecret := makeWebhookSecret()

	id, err := api.Database.SaveProject(project, webhookSecret)
	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: err.Error(),
		}, 500)
		return
	}

	api.Database.SetManagerRoleOn(manager.(*storage.Manager).Id, id,
		storage.RoleManageAccess|storage.RoleRead|storage.RoleEdit|storage.RoleSecret|storage.RoleMaintenance)
	r.OkJson(JsonResponse{
		Ok: true,
		Content: CreateProjectResponse{
			Id: id,
		},
	})
	logrus.WithFields(logrus.Fields{
		"project": project,
	}).Debug("Created project")
}

func makeWebhookSecret() string {
	return uuid.New().String()
}

func (api *WebAPI) UpdateProject(r *Request) {

	id, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	if err != nil || id <= 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid project id",
		}, 400)
		return
	}

	updateReq := &UpdateProjectRequest{}
	err = json.Unmarshal(r.Ctx.Request.Body(), updateReq)
	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}

	if !updateReq.isValid(id) {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid request",
		}, 400)
		return
	}

	project := &storage.Project{
		Id:         id,
		Name:       updateReq.Name,
		CloneUrl:   updateReq.CloneUrl,
		GitRepo:    updateReq.GitRepo,
		Priority:   updateReq.Priority,
		Motd:       updateReq.Motd,
		Public:     updateReq.Public,
		Hidden:     updateReq.Hidden,
		Chain:      updateReq.Chain,
		Paused:     updateReq.Paused,
		AssignRate: updateReq.AssignRate,
		SubmitRate: updateReq.SubmitRate,
		Version:    updateReq.Version,
	}
	sess := api.Session.StartFasthttp(r.Ctx)
	manager := sess.Get("manager")

	if !isActionOnProjectAuthorized(project.Id, manager, storage.RoleEdit, api.Database) {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Unauthorized",
		}, 403)
		return
	}

	if project.Chain != 0 && !isActionOnProjectAuthorized(project.Chain, manager, storage.RoleEdit, api.Database) {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Unauthorized (You need RoleEdit on the project you wish to chain tasks to)",
		}, 403)
		return
	}

	err = api.Database.UpdateProject(project)
	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: err.Error(),
		}, 500)

		logrus.WithError(err).WithFields(logrus.Fields{
			"project": project,
		}).Warn("Error during project update")
	} else {
		api.SubmitLimiters.Delete(project.Id)
		api.AssignLimiters.Delete(project.Id)

		r.OkJson(JsonResponse{
			Ok: true,
		})

		logrus.WithFields(logrus.Fields{
			"project": project,
		}).Debug("Updated project")
	}
}

func (api *WebAPI) isProjectCreationAuthorized(project *storage.Project, manager interface{}) bool {

	if manager == nil {
		return false
	}

	if project.Public && !manager.(*storage.Manager).WebsiteAdmin {
		return false
	}

	if project.Chain != 0 {
		chainsTo := api.Database.GetProject(project.Chain)
		if chainsTo == nil {
			return false
		}

		if !isActionOnProjectAuthorized(chainsTo.Id, manager.(*storage.Manager),
			storage.RoleEdit, api.Database) {
			return false
		}
	}

	return true
}

func isActionOnProjectAuthorized(project int64, manager interface{},
	requiredRole storage.ManagerRole, db *storage.Database) bool {

	if manager == nil {
		return false
	}

	if manager.(*storage.Manager).WebsiteAdmin {
		return true
	}

	role := db.GetManagerRoleOn(manager.(*storage.Manager), project)
	if role&requiredRole != 0 {
		return true
	}

	return false
}

func isProjectReadAuthorized(project *storage.Project, manager interface{}, db *storage.Database) bool {

	if project.Public || !project.Hidden {
		return true
	}
	if manager == nil {
		return false
	}
	if manager.(*storage.Manager).WebsiteAdmin {
		return true
	}
	role := db.GetManagerRoleOn(manager.(*storage.Manager), project.Id)
	if role&storage.RoleRead == 1 {
		return true
	}

	return false
}

func (api *WebAPI) GetProjectList(r *Request) {

	sess := api.Session.StartFasthttp(r.Ctx)
	manager := sess.Get("manager")

	var id int64
	if manager == nil {
		id = 0
	} else {
		id = manager.(*storage.Manager).Id
	}

	projects := api.Database.GetAllProjects(id)

	r.OkJson(JsonResponse{
		Ok: true,
		Content: GetProjectListResponse{
			Projects: projects,
		},
	})
}

func (api *WebAPI) GetAssigneeStatsForProject(r *Request) {

	id, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	if err != nil || id <= 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid worker id",
		}, 400)
		return
	}

	stats := api.Database.GetAssigneeStats(id, 16)

	r.OkJson(JsonResponse{
		Ok: true,
		Content: GetAssigneeStatsForProjectResponse{
			Assignees: stats,
		},
	})
}

func (api *WebAPI) GetWorkerAccessListForProject(r *Request) {

	sess := api.Session.StartFasthttp(r.Ctx)
	manager := sess.Get("manager")

	id, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	if err != nil || id <= 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid worker id",
		}, 400)
		return
	}

	if !isActionOnProjectAuthorized(id, manager, storage.RoleManageAccess, api.Database) {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Unauthorized",
		}, 403)
		return
	}
	accesses := api.Database.GetAllAccesses(id)

	r.OkJson(JsonResponse{
		Ok: true,
		Content: GetWorkerAccessListForProjectResponse{
			Accesses: accesses,
		},
	})
}

func (api *WebAPI) CreateWorkerAccess(r *Request) {

	req := &CreateWorkerAccessRequest{}
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
			Message: "Invalid request",
		}, 400)
		return
	}

	worker, err := api.validateSignature(r)
	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: err.Error(),
		}, 401)
		return
	}

	res := api.Database.SaveAccessRequest(&storage.WorkerAccess{
		Worker:  *worker,
		Submit:  req.Submit,
		Assign:  req.Assign,
		Project: req.Project,
	})

	if res {
		r.OkJson(JsonResponse{
			Ok: true,
		})
	} else {
		r.Json(JsonResponse{
			Ok: false,
			Message: "Project is public, you already have " +
				"an active request or you already have access to this project",
		}, 400)
	}
}

func (api *WebAPI) AcceptAccessRequest(r *Request) {

	pid, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	if err != nil || pid <= 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid worker id",
		}, 400)
		return
	}

	wid, err := strconv.ParseInt(r.Ctx.UserValue("wid").(string), 10, 64)
	if err != nil || wid <= 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid worker id",
		}, 400)
		return
	}

	sess := api.Session.StartFasthttp(r.Ctx)
	manager := sess.Get("manager")

	if !isActionOnProjectAuthorized(pid, manager, storage.RoleManageAccess, api.Database) {
		r.Json(JsonResponse{
			Message: "Unauthorized",
			Ok:      false,
		}, 403)
		return
	}

	ok := api.Database.AcceptAccessRequest(wid, pid)

	if ok {
		r.OkJson(JsonResponse{
			Ok: true,
		})
	} else {
		r.OkJson(JsonResponse{
			Ok:      false,
			Message: "Worker did not have access to this project",
		})
	}
}

func (api *WebAPI) RejectAccessRequest(r *Request) {

	pid, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	if err != nil || pid <= 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid project id",
		}, 400)
		return
	}

	wid, err := strconv.ParseInt(r.Ctx.UserValue("wid").(string), 10, 64)
	if err != nil || wid <= 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid worker id",
		}, 400)
		return
	}

	ok := api.Database.RejectAccessRequest(wid, pid)

	if ok {
		r.OkJson(JsonResponse{
			Ok: true,
		})
	} else {
		r.OkJson(JsonResponse{
			Ok:      false,
			Message: "Worker did not have access to this project",
		})
	}
}

func (api *WebAPI) SetManagerRoleOnProject(r *Request) {

	pid, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	if err != nil || pid <= 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid project id",
		}, 400)
		return
	}

	req := &SetManagerRoleOnProjectRequest{}
	err = json.Unmarshal(r.Ctx.Request.Body(), req)
	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}

	sess := api.Session.StartFasthttp(r.Ctx)
	manager := sess.Get("manager")

	if !isActionOnProjectAuthorized(pid, manager, storage.RoleManageAccess, api.Database) {
		r.Json(JsonResponse{
			Message: "Unauthorized",
			Ok:      false,
		}, 403)
		return
	}

	api.Database.SetManagerRoleOn(req.Manager, pid, req.Role)
	r.OkJson(JsonResponse{
		Ok: true,
	})
}

func (api *WebAPI) SetSecret(r *Request) {

	pid, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	if err != nil || pid <= 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid project id",
		}, 400)
		return
	}

	sess := api.Session.StartFasthttp(r.Ctx)
	manager := sess.Get("manager")

	if !isActionOnProjectAuthorized(pid, manager, storage.RoleSecret, api.Database) {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Unauthorized",
		}, 403)
		return
	}

	req := &SetSecretRequest{}
	err = json.Unmarshal(r.Ctx.Request.Body(), req)
	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}

	api.Database.SetSecret(pid, req.Secret)

	r.OkJson(JsonResponse{
		Ok: true,
	})
}

func (api *WebAPI) GetSecret(r *Request) {

	pid, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	if err != nil || pid <= 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid project id",
		}, 400)
		return
	}

	var secret string

	worker, err := api.validateSignature(r)
	if err == nil {
		secret, err = api.Database.GetSecret(pid, worker.Id)
		if err != nil {
			r.Json(JsonResponse{
				Ok:      false,
				Message: "Unauthorized",
			}, 403)
			return
		}
		r.OkJson(JsonResponse{
			Ok: true,
			Content: GetSecretResponse{
				Secret: secret,
			},
		})
		return
	}

	sess := api.Session.StartFasthttp(r.Ctx)
	manager := sess.Get("manager")

	if !isActionOnProjectAuthorized(pid, manager, storage.RoleSecret, api.Database) {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Unauthorized",
		}, 403)
		return
	}

	secret, _ = api.Database.GetSecret(pid, 0)
	r.OkJson(JsonResponse{
		Ok: true,
		Content: GetSecretResponse{
			Secret: secret,
		},
	})
}

func (api *WebAPI) GetWebhookSecret(r *Request) {

	pid, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	if err != nil || pid <= 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid project id",
		}, 400)
		return
	}

	sess := api.Session.StartFasthttp(r.Ctx)
	manager := sess.Get("manager")

	if !isActionOnProjectAuthorized(pid, manager, storage.RoleSecret, api.Database) {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Unauthorized",
		}, 403)
		return
	}

	secret, err := api.Database.GetWebhookSecret(pid)
	r.OkJson(JsonResponse{
		Ok: true,
		Content: GetWebhookSecretResponse{
			WebhookSecret: secret,
		},
	})
}

func (api *WebAPI) SetWebhookSecret(r *Request) {

	pid, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	if err != nil || pid <= 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid project id",
		}, 400)
		return
	}

	req := &SetWebhookSecretRequest{}
	err = json.Unmarshal(r.Ctx.Request.Body(), req)
	if err != nil {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}

	sess := api.Session.StartFasthttp(r.Ctx)
	manager := sess.Get("manager")

	if !isActionOnProjectAuthorized(pid, manager, storage.RoleSecret, api.Database) {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Unauthorized",
		}, 403)
		return
	}

	err = api.Database.SetWebhookSecret(pid, req.WebhookSecret)
	if err == nil {
		r.OkJson(JsonResponse{
			Ok: true,
		})
	} else {
		r.OkJson(JsonResponse{
			Ok:      false,
			Message: err.Error(),
		})
	}
}

func (api *WebAPI) ResetFailedTasks(r *Request) {

	pid, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	if err != nil || pid <= 0 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid project id",
		}, 400)
		return
	}

	sess := api.Session.StartFasthttp(r.Ctx)
	manager := sess.Get("manager")

	if !isActionOnProjectAuthorized(pid, manager, storage.RoleMaintenance, api.Database) {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Unauthorized",
		}, 403)
		return
	}

	res := api.Database.ResetFailedTasks(pid)

	r.OkJson(JsonResponse{
		Ok: true,
		Content: ResetFailedTaskResponse{
			AffectedTasks: res,
		},
	})
}
