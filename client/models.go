package client

import "github.com/simon987/task_tracker/storage"

type Worker struct {
	Id     int64  `json:"id"`
	Alias  string `json:"alias,omitempty"`
	Secret []byte `json:"secret"`
}

type CreateWorkerResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Content struct {
		Worker *storage.Worker `json:"worker"`
	} `json:"content"`
}

type AssignTaskResponse struct {
	Ok             bool    `json:"ok"`
	Message        string  `json:"message"`
	RateLimitDelay float64 `json:"rate_limit_delay,omitempty"`
	Content        struct {
		storage.Task `json:"task"`
	} `json:"content"`
}

type ReleaseTaskResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Content struct {
		Updated bool `json:"updated"`
	} `json:"content"`
}

type ProjectSecretResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Content struct {
		Secret string `json:"secret"`
	} `json:"content"`
}

type ProjectListResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Content struct {
		Projects []storage.Project `json:"projects,omitempty"`
	} `json:"content"`
}
