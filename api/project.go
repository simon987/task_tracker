package api

import (
	"github.com/Sirupsen/logrus"
	"src/task_tracker/storage"
	"strconv"
)

type CreateProjectRequest struct {
	Name     string `json:"name"`
	GitUrl   string `json:"git_url"`
	Version  string `json:"version"`
	Priority int64  `json:"priority"`
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

func (api *WebAPI) ProjectCreate(r *Request) {

	createReq := &CreateProjectRequest{}
	if r.GetJson(createReq) {

		project := &storage.Project{
			Name:     createReq.Name,
			Version:  createReq.Version,
			GitUrl:   createReq.GitUrl,
			Priority: createReq.Priority,
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
}

func isValidProject(project *storage.Project) bool {
	if len(project.Name) <= 0 {
		return false
	}

	return true
}

func (api *WebAPI) ProjectGet(r *Request) {

	id, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	handleErr(err, r)

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
