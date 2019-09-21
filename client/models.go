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
		Task *struct {
			Id                int64              `json:"id"`
			Priority          int16              `json:"priority"`
			Project           *storage.Project   `json:"project"`
			Assignee          int64              `json:"assignee"`
			Retries           int16              `json:"retries"`
			MaxRetries        int64              `json:"max_retries"`
			Status            storage.TaskStatus `json:"status"`
			Recipe            string             `json:"recipe"`
			MaxAssignTime     int64              `json:"max_assign_time"`
			AssignTime        int64              `json:"assign_time"`
			VerificationCount int16              `json:"verification_count"`
		} `json:"task"`
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
