package test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/simon987/task_tracker/api"
	"github.com/simon987/task_tracker/config"
	"github.com/simon987/task_tracker/storage"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

type SessionContext struct {
	Manager       *storage.Manager
	SessionCookie *http.Cookie
}

type ResponseHeader struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

func Post(path string, x interface{}, worker *storage.Worker, s *http.Client) *http.Response {

	if s == nil {
		s = &http.Client{}
	}

	body, err := json.Marshal(x)
	buf := bytes.NewBuffer(body)

	req, err := http.NewRequest("POST", "http://"+config.Cfg.ServerAddr+path, buf)
	handleErr(err)

	if worker != nil {
		req.Header.Add("X-Worker-Id", strconv.FormatInt(worker.Id, 10))
		secretHeader := base64.StdEncoding.EncodeToString(worker.Secret)
		req.Header.Add("X-Secret", string(secretHeader))
	}

	r, err := s.Do(req)
	handleErr(err)

	return r
}

func Get(path string, worker *storage.Worker, s *http.Client) *http.Response {

	if s == nil {
		s = &http.Client{}
	}

	url := "http://" + config.Cfg.ServerAddr + path
	req, err := http.NewRequest("GET", url, nil)

	if worker != nil {
		req.Header.Add("X-Worker-Id", strconv.FormatInt(worker.Id, 10))
		secretHeader := base64.StdEncoding.EncodeToString(worker.Secret)
		req.Header.Add("X-Secret", string(secretHeader))
	}

	r, err := s.Do(req)
	handleErr(err)

	return r
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func GenericJson(body io.ReadCloser) map[string]interface{} {

	var obj map[string]interface{}

	data, _ := ioutil.ReadAll(body)

	err := json.Unmarshal(data, &obj)
	handleErr(err)

	return obj
}

func UnmarshalResponse(r *http.Response, result interface{}) {
	data, err := ioutil.ReadAll(r.Body)
	fmt.Println(string(data))
	err = json.Unmarshal(data, result)
	handleErr(err)
}

type WorkerAR struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Content struct {
		Worker *storage.Worker `json:"worker"`
	} `json:"content"`
}

type RegisterAR struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Content struct {
		Manager *storage.Manager `json:"manager"`
	} `json:"content"`
}

type ProjectAR struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Content struct {
		Project *storage.Project `json:"project"`
	} `json:"content"`
}

type CreateProjectAR struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Content struct {
		Id int64 `json:"id"`
	} `json:"content"`
}

type InfoAR struct {
	Ok       bool   `json:"ok"`
	Message  string `json:"message"`
	api.Info `json:"content"`
}

type LogsAR struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Content struct {
		Logs *[]storage.LogEntry `json:"logs"`
	} `json:"content"`
}

type ProjectListAR struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Content struct {
		Projects *[]storage.Project `json:"projects"`
	} `json:"content"`
}

type TaskAR struct {
	Ok             bool    `json:"ok"`
	Message        string  `json:"message"`
	RateLimitDelay float64 `json:"rate_limit_delay,omitempty"`
	Content        struct {
		Task *storage.Task `json:"task"`
	} `json:"content"`
}

type ReleaseAR struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Content struct {
		Updated bool `json:"updated"`
	} `json:"content"`
}

type WebhookSecretAR struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Content struct {
		WebhookSecret string `json:"webhook_secret"`
	} `json:"content"`
}

type AccountAR struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
	Content struct {
		*storage.Manager `json:"manager"`
	} `json:"content"`
}
