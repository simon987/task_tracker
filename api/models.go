package api

import (
	"encoding/json"
	"github.com/simon987/task_tracker/storage"
	"golang.org/x/time/rate"
)

const (
	MinPasswordLength = 8
	MinUsernameLength = 3
	MaxUsernameLength = 16
)

type JsonResponse struct {
	Ok             bool        `json:"ok"`
	Message        string      `json:"message,omitempty"`
	RateLimitDelay float64     `json:"rate_limit_delay,omitempty"`
	Content        interface{} `json:"content,omitempty"`
}

type GitPayload struct {
	Ref        string `json:"ref"`
	Before     string `json:"before"`
	After      string `json:"after"`
	Repository struct {
		Id    int64 `json:"id"`
		Owner struct {
			Id       int64  `json:"id"`
			Username string `json:"username"`
			Login    string `json:"login"`
			FullName string `json:"full_name"`
			Email    string `json:"email"`
		} `json:"owner"`
		Name          string `json:"name"`
		FullName      string `json:"full_name"`
		Private       bool   `json:"private"`
		Fork          bool   `json:"fork"`
		Size          int64  `json:"size"`
		HtmlUrl       string `json:"html_url"`
		SshUrl        string `json:"ssh_url"`
		CloneUrl      string `json:"clone_url"`
		DefaultBranch string `json:"default_branch"`
	} `json:"repository"`
}

func (g GitPayload) String() string {
	jsonBytes, _ := json.Marshal(g)
	return string(jsonBytes)
}

type ErrorResponse struct {
	Message    string `json:"message"`
	StackTrace string `json:"stack_trace"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Manager *storage.Manager `json:"manager"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (r *RegisterRequest) isValid() bool {
	return MinUsernameLength <= len(r.Username) &&
		len(r.Username) <= MaxUsernameLength &&
		MinPasswordLength <= len(r.Password)
}

type GetAccountDetailsResponse struct {
	LoggedIn bool             `json:"logged_in"`
	Manager  *storage.Manager `json:"manager,omitempty"`
}

type GetManagerListResponse struct {
	Managers *[]storage.Manager `json:"managers"`
}

type GetManagerListWithRoleOnResponse struct {
	Managers *[]storage.ManagerRoleOn `json:"managers"`
}

type GetLogRequest struct {
	Level storage.LogLevel `json:"level"`
	Since int64            `json:"since"`
}

func (r GetLogRequest) isValid() bool {

	if r.Since <= 0 {
		return false
	}

	return true
}

type LogRequest struct {
	Scope     string `json:"scope"`
	Message   string `json:"Message"`
	TimeStamp int64  `json:"timestamp"`
	worker    *storage.Worker
}

type GetLogResponse struct {
	Logs *[]storage.LogEntry `json:"logs"`
}

type GetSnapshotsResponse struct {
	Snapshots *[]storage.ProjectMonitoringSnapshot `json:"snapshots,omitempty"`
}

type CreateProjectRequest struct {
	Name       string     `json:"name"`
	CloneUrl   string     `json:"clone_url"`
	GitRepo    string     `json:"git_repo"`
	Version    string     `json:"version"`
	Priority   int64      `json:"priority"`
	Motd       string     `json:"motd"`
	Public     bool       `json:"public"`
	Hidden     bool       `json:"hidden"`
	Chain      int64      `json:"chain"`
	AssignRate rate.Limit `json:"assign_rate"`
	SubmitRate rate.Limit `json:"submit_rate"`
}

func (req *CreateProjectRequest) isValid() bool {
	if len(req.Name) <= 0 {
		return false
	}
	if req.Priority < 0 {
		return false
	}
	if req.Hidden && req.Public {
		return false
	}
	return true
}

type UpdateProjectRequest struct {
	Name       string     `json:"name"`
	CloneUrl   string     `json:"clone_url"`
	GitRepo    string     `json:"git_repo"`
	Priority   int64      `json:"priority"`
	Motd       string     `json:"motd"`
	Public     bool       `json:"public"`
	Hidden     bool       `json:"hidden"`
	Chain      int64      `json:"chain"`
	Paused     bool       `json:"paused"`
	AssignRate rate.Limit `json:"assign_rate"`
	SubmitRate rate.Limit `json:"submit_rate"`
	Version    string     `json:"version"`
}

func (req *UpdateProjectRequest) isValid(pid int64) bool {
	if len(req.Name) <= 0 {
		return false
	}
	if req.Priority < 0 {
		return false
	}
	if req.Hidden && req.Public {
		return false
	}
	if req.Chain == pid {
		return false
	}
	return true
}

type CreateProjectResponse struct {
	Id int64 `json:"id,omitempty"`
}

type GetProjectResponse struct {
	Project *storage.Project `json:"project,omitempty"`
}

type GetProjectListResponse struct {
	Projects *[]storage.Project `json:"projects,omitempty"`
}

type GetAssigneeStatsForProjectResponse struct {
	Assignees *[]storage.AssignedTasks `json:"assignees"`
}

type GetWorkerAccessListForProjectResponse struct {
	Accesses *[]storage.WorkerAccess `json:"accesses,omitempty"`
}

type SubmitTaskRequest struct {
	Project           int64  `json:"project"`
	MaxRetries        int16  `json:"max_retries"`
	Recipe            string `json:"recipe"`
	Priority          int16  `json:"priority"`
	MaxAssignTime     int64  `json:"max_assign_time"`
	Hash64            int64  `json:"hash_u64"`
	UniqueString      string `json:"unique_string"`
	VerificationCount int16  `json:"verification_count"`
}

func (req *SubmitTaskRequest) IsValid() bool {
	if req.MaxRetries < 0 {
		return false
	}
	if len(req.Recipe) <= 0 {
		return false
	}
	if req.Hash64 != 0 && len(req.UniqueString) != 0 {
		return false
	}
	if req.Project == 0 {
		return false
	}

	return true
}

type BulkSubmitTaskRequest struct {
	Requests []SubmitTaskRequest `json:"requests"`
}

func (reqs BulkSubmitTaskRequest) IsValid() bool {

	if reqs.Requests == nil {
		return false
	}

	if len(reqs.Requests) == 0 {
		return false
	}

	for _, req := range reqs.Requests {
		if !req.IsValid() {
			return false
		}
	}
	return true
}

type ReleaseTaskRequest struct {
	TaskId       int64              `json:"task_id"`
	Result       storage.TaskResult `json:"result"`
	Verification int64              `json:"verification"`
}

func (r *ReleaseTaskRequest) IsValid() bool {
	return r.TaskId != 0
}

type ReleaseTaskResponse struct {
	Updated bool `json:"updated"`
}

type CreateTaskResponse struct {
}

type GetTaskResponse struct {
	Task *storage.Task `json:"task,omitempty"`
}

type UpdateWorkerRequest struct {
	Alias string `json:"alias"`
}

type WorkerSetPausedRequest struct {
	Worker int64 `json:"worker"`
	Paused bool  `json:"paused"`
}

type CreateWorkerRequest struct {
	Alias string `json:"alias"`
}

func (req *CreateWorkerRequest) isValid() bool {
	if req.Alias == "unassigned" {
		//Reserved alias
		return false
	}

	return true
}

type CreateWorkerResponse struct {
	Worker *storage.Worker `json:"worker,omitempty"`
}

type GetWorkerResponse struct {
	Worker *storage.Worker `json:"worker,omitempty"`
}

type GetAllWorkerStatsResponse struct {
	Stats *[]storage.WorkerStats `json:"stats"`
}

type CreateWorkerAccessRequest struct {
	Assign  bool  `json:"assign"`
	Submit  bool  `json:"submit"`
	Project int64 `json:"project"`
}

func (w *CreateWorkerAccessRequest) isValid() bool {
	if !w.Assign && !w.Submit {
		return false
	}
	return true
}

type SetManagerRoleOnProjectRequest struct {
	Manager int64               `json:"manager"`
	Role    storage.ManagerRole `json:"role"`
}

type Info struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type SetSecretRequest struct {
	Secret string `json:"secret"`
}

type GetSecretResponse struct {
	Secret string `json:"secret"`
}

type SetWebhookSecretRequest struct {
	WebhookSecret string `json:"webhook_secret"`
}

type GetWebhookSecretResponse struct {
	WebhookSecret string `json:"webhook_secret"`
}

type ResetFailedTaskResponse struct {
	AffectedTasks int64 `json:"affected_tasks"`
}

type HardResetResponse struct {
	AffectedTasks int64 `json:"affected_tasks"`
}

type ReclaimAssignedTasksResponse struct {
	AffectedTasks int64 `json:"affected_tasks"`
}
