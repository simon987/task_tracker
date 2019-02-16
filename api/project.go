package api

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/simon987/task_tracker/storage"
	"strconv"
)

type CreateProjectRequest struct {
	Name     string `json:"name"`
	CloneUrl string `json:"clone_url"`
	GitRepo  string `json:"git_repo"`
	Version  string `json:"version"`
	Priority int64  `json:"priority"`
	Motd     string `json:"motd"`
	Public   bool   `json:"public"`
	Hidden   bool   `json:"hidden"`
	Chain    int64  `json:"chain"`
}

type UpdateProjectRequest struct {
	Name     string `json:"name"`
	CloneUrl string `json:"clone_url"`
	GitRepo  string `json:"git_repo"`
	Priority int64  `json:"priority"`
	Motd     string `json:"motd"`
	Public   bool   `json:"public"`
	Hidden   bool   `json:"hidden"`
	Chain    int64  `json:"chain"`
}

type UpdateProjectResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}

type CreateProjectResponse struct {
	Ok      bool   `json:"ok"`
	Id      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

type GetProjectResponse struct {
	Ok      bool             `json:"ok"`
	Message string           `json:"message,omitempty"`
	Project *storage.Project `json:"project,omitempty"`
}

type GetAllProjectsResponse struct {
	Ok       bool               `json:"ok"`
	Message  string             `json:"message,omitempty"`
	Projects *[]storage.Project `json:"projects,omitempty"`
}

type GetAssigneeStatsResponse struct {
	Ok        bool                     `json:"ok"`
	Message   string                   `json:"message,omitempty"`
	Assignees *[]storage.AssignedTasks `json:"assignees"`
}

type WorkerRequestAccessResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
}

type ProjectGetAccessRequestsResponse struct {
	Ok       bool              `json:"ok"`
	Message  string            `json:"message,omitempty"`
	Requests *[]storage.Worker `json:"requests,omitempty"`
}

func (api *WebAPI) ProjectCreate(r *Request) {

	sess := api.Session.StartFasthttp(r.Ctx)
	manager := sess.Get("manager")

	createReq := &CreateProjectRequest{}
	err := json.Unmarshal(r.Ctx.Request.Body(), createReq)
	if err != nil {
		r.Json(CreateProjectResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}
	project := &storage.Project{
		Name:     createReq.Name,
		Version:  createReq.Version,
		CloneUrl: createReq.CloneUrl,
		GitRepo:  createReq.GitRepo,
		Priority: createReq.Priority,
		Motd:     createReq.Motd,
		Public:   createReq.Public,
		Hidden:   createReq.Hidden,
		Chain:    createReq.Chain,
	}

	if !isValidProject(project) {
		logrus.WithFields(logrus.Fields{
			"project": project,
		}).Warn("Invalid project")

		r.Json(CreateProjectResponse{
			Ok:      false,
			Message: "Invalid project",
		}, 400)
		return
	}

	if !isProjectCreationAuthorized(project, manager) {
		logrus.WithFields(logrus.Fields{
			"project": project,
		}).Warn("Unauthorized project creation")

		r.Json(CreateProjectResponse{
			Ok:      false,
			Message: "You are not permitted to create a project with this configuration",
		}, 400)
		return
	}

	id, err := api.Database.SaveProject(project)
	if err != nil {
		r.Json(CreateProjectResponse{
			Ok:      false,
			Message: err.Error(),
		}, 500)
		return
	}
	r.OkJson(CreateProjectResponse{
		Ok: true,
		Id: id,
	})
	logrus.WithFields(logrus.Fields{
		"project": project,
	}).Debug("Created project")
}

func (api *WebAPI) ProjectUpdate(r *Request) {

	id, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	if err != nil || id <= 0 {
		r.Json(CreateProjectResponse{
			Ok:      false,
			Message: "Invalid project id",
		}, 400)
		return
	}

	updateReq := &UpdateProjectRequest{}
	err = json.Unmarshal(r.Ctx.Request.Body(), updateReq)
	if err != nil {
		r.Json(CreateProjectResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}
	project := &storage.Project{
		Id:       id,
		Name:     updateReq.Name,
		CloneUrl: updateReq.CloneUrl,
		GitRepo:  updateReq.GitRepo,
		Priority: updateReq.Priority,
		Motd:     updateReq.Motd,
		Public:   updateReq.Public,
		Hidden:   updateReq.Hidden,
		Chain:    updateReq.Chain,
	}

	if isValidProject(project) {
		err := api.Database.UpdateProject(project)

		if err != nil {
			r.Json(CreateProjectResponse{
				Ok:      false,
				Message: err.Error(),
			}, 500)

			logrus.WithError(err).WithFields(logrus.Fields{
				"project": project,
			}).Warn("Error during project update")
		} else {
			r.OkJson(UpdateProjectResponse{
				Ok: true,
			})

			logrus.WithFields(logrus.Fields{
				"project": project,
			}).Debug("Updated project")
		}

	} else {
		logrus.WithFields(logrus.Fields{
			"project": project,
		}).Warn("Invalid project")

		r.Json(CreateProjectResponse{
			Ok:      false,
			Message: "Invalid project",
		}, 400)
	}
}

func isValidProject(project *storage.Project) bool {
	if len(project.Name) <= 0 {
		return false
	}
	if project.Priority < 0 {
		return false
	}

	return true
}

func isProjectCreationAuthorized(project *storage.Project, manager interface{}) bool {

	return true
	if manager == nil {
		return false
	}

	if project.Public && manager.(*storage.Manager).WebsiteAdmin {
		return false
	}
	return true
}

func isProjectUpdateAuthorized(project *storage.Project, manager interface{}, db *storage.Database) bool {

	var man storage.Manager
	if manager != nil {
		man = manager.(storage.Manager)
	}

	if man.WebsiteAdmin {
		return true
	}

	role := db.ManagerHasRoleOn(&man, project.Id)
	if role&storage.ROLE_EDIT == 1 {
		return true
	}

	return false
}

func (api *WebAPI) ProjectGet(r *Request) {

	id, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	handleErr(err, r) //todo handle invalid id

	project := api.Database.GetProject(id)

	if project != nil {
		r.OkJson(GetProjectResponse{
			Ok:      true,
			Project: project,
		})
	} else {
		r.Json(GetProjectResponse{
			Ok:      false,
			Message: "Project not found",
		}, 404)
	}
}

func (api *WebAPI) ProjectGetAllProjects(r *Request) {

	worker, _ := api.validateSignature(r)

	var id int64
	if worker == nil {
		id = 0
	} else {
		id = worker.Id
	}
	projects := api.Database.GetAllProjects(id)

	r.OkJson(GetAllProjectsResponse{
		Ok:       true,
		Projects: projects,
	})
}

func (api *WebAPI) ProjectGetAssigneeStats(r *Request) {

	id, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	handleErr(err, r) //todo handle invalid id

	stats := api.Database.GetAssigneeStats(id, 16)

	r.OkJson(GetAssigneeStatsResponse{
		Ok:        true,
		Assignees: stats,
	})
}

func (api *WebAPI) ProjectGetAccessRequests(r *Request) {

	sess := api.Session.StartFasthttp(r.Ctx)
	manager := sess.Get("manager")

	id, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	handleErr(err, r) //todo handle invalid id

	if manager == nil {
		r.Json(ProjectGetAccessRequestsResponse{
			Ok:      false,
			Message: "Unauthorized",
		}, 401)
		return
	}

	if !manager.(*storage.Manager).WebsiteAdmin &&
		api.Database.ManagerHasRoleOn(manager.(*storage.Manager), 1)&
			storage.ROLE_MANAGE_ACCESS == 0 {
		r.Json(ProjectGetAccessRequestsResponse{
			Ok:      false,
			Message: "Unauthorized",
		}, 403)
		return
	}
	requests := api.Database.GetAllAccessRequests(id)

	r.OkJson(ProjectGetAccessRequestsResponse{
		Ok:       true,
		Requests: requests,
	})
}

func (api *WebAPI) WorkerRequestAccess(r *Request) {

	id, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	handleErr(err, r) //todo handle invalid id

	worker, err := api.validateSignature(r)
	if err != nil {
		r.Json(WorkerRequestAccessResponse{
			Ok:      false,
			Message: err.Error(),
		}, 401)
	}

	res := api.Database.SaveAccessRequest(worker, id)

	if res {
		r.OkJson(WorkerRequestAccessResponse{
			Ok: true,
		})
	} else {
		r.Json(WorkerRequestAccessResponse{
			Ok: false,
			Message: "Project is public, you already have " +
				"an active request or you already have access to this project",
		}, 400)
	}
}
