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
}

type UpdateProjectRequest struct {
	Name     string `json:"name"`
	CloneUrl string `json:"clone_url"`
	GitRepo  string `json:"git_repo"`
	Priority int64  `json:"priority"`
	Motd     string `json:"motd"`
	Public   bool   `json:"public"`
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

func (api *WebAPI) ProjectCreate(r *Request) {

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
	}

	if isValidProject(project) {
		id, err := api.Database.SaveProject(project)

		if err != nil {
			r.Json(CreateProjectResponse{
				Ok:      false,
				Message: err.Error(),
			}, 500)
		} else {
			r.OkJson(CreateProjectResponse{
				Ok: true,
				Id: id,
			})
			logrus.WithFields(logrus.Fields{
				"project": project,
			}).Debug("Created project")
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

	projects := api.Database.GetAllProjects()

	r.OkJson(GetAllProjectsResponse{
		Ok:       true,
		Projects: projects,
	})
}
