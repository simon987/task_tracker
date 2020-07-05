package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/simon987/task_tracker/api"
	"github.com/simon987/task_tracker/storage"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type TaskTrackerClient struct {
	worker        *Worker
	httpClient    http.Client
	serverAddress string

	secretB64   string
	workerIdStr string
}

func New(serverAddress string) *TaskTrackerClient {

	client := new(TaskTrackerClient)
	client.serverAddress = serverAddress

	return client
}

func (c *TaskTrackerClient) SetWorker(worker *Worker) {
	c.worker = worker
	c.secretB64 = base64.StdEncoding.EncodeToString(worker.Secret)
	c.workerIdStr = strconv.FormatInt(worker.Id, 10)
}

func (c *TaskTrackerClient) get(path string) *http.Response {

	url := c.serverAddress + path
	req, err := http.NewRequest("GET", url, nil)

	if c.worker != nil {
		req.Header.Add("X-Worker-Id", c.workerIdStr)
		req.Header.Add("X-Secret", c.secretB64)
	}

	r, err := c.httpClient.Do(req)
	handleErr(err)

	return r
}

func (c *TaskTrackerClient) post(path string, x interface{}) *http.Response {

	body, err := json.Marshal(x)
	buf := bytes.NewBuffer(body)

	req, err := http.NewRequest("POST", c.serverAddress+path, buf)
	handleErr(err)

	if c.worker != nil {
		req.Header.Add("X-Worker-Id", c.workerIdStr)
		req.Header.Add("X-Secret", c.secretB64)
	}

	r, err := c.httpClient.Do(req)
	handleErr(err)

	return r
}

func unmarshalResponse(r *http.Response, result interface{}) error {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, result)
	if err != nil {
		return err
	}
	return nil
}

func (c TaskTrackerClient) MakeWorker(alias string) (*Worker, error) {

	httpResp := c.post("/worker/create", api.CreateWorkerRequest{
		Alias: alias,
	})
	var jsonResp CreateWorkerResponse
	err := unmarshalResponse(httpResp, &jsonResp)

	if err == nil {
		clientWorker := Worker{
			Alias:  jsonResp.Content.Worker.Alias,
			Secret: jsonResp.Content.Worker.Secret,
			Id:     jsonResp.Content.Worker.Id,
		}
		return &clientWorker, nil
	}

	return nil, err
}

func (c TaskTrackerClient) FetchTask(projectId int) (*AssignTaskResponse, error) {

	httpResp := c.get("/task/get/" + strconv.Itoa(projectId))
	var jsonResp AssignTaskResponse
	err := unmarshalResponse(httpResp, &jsonResp)

	//TODO: Handle rate limiting here?

	return &jsonResp, err
}

func (c TaskTrackerClient) ReleaseTask(req api.ReleaseTaskRequest) (*ReleaseTaskResponse, error) {

	httpResp := c.post("/task/release", req)
	var jsonResp ReleaseTaskResponse
	err := unmarshalResponse(httpResp, &jsonResp)

	return &jsonResp, err
}

func (c TaskTrackerClient) SubmitTask(req api.SubmitTaskRequest) (AssignTaskResponse, error) {

	httpResp := c.post("/task/submit", req)
	var jsonResp AssignTaskResponse
	err := unmarshalResponse(httpResp, &jsonResp)

	//TODO: Handle rate limiting here?

	return jsonResp, err
}

func (c TaskTrackerClient) GetProjectSecret(projectId int) (string, error) {

	httpResp := c.get("/project/secret/" + strconv.Itoa(projectId))
	var jsonResp ProjectSecretResponse
	err := unmarshalResponse(httpResp, &jsonResp)

	return jsonResp.Content.Secret, err
}

func (c TaskTrackerClient) GetProjectList() ([]storage.Project, error) {

	httpResp := c.get("/project/list")
	var jsonResp ProjectListResponse
	err := unmarshalResponse(httpResp, &jsonResp)

	return jsonResp.Content.Projects, err
}

func (c TaskTrackerClient) RequestAccess(req api.CreateWorkerAccessRequest) (api.JsonResponse, error) {

	httpResp := c.post("/project/request_access", req)
	var jsonResp api.JsonResponse
	err := unmarshalResponse(httpResp, &jsonResp)

	return jsonResp, err
}

func (c TaskTrackerClient) Log(level storage.LogLevel, message string) error {

	var levelString string
	req := api.LogRequest{
		Scope:     "task_tracker go client",
		Message:   message,
		TimeStamp: time.Now().Unix(),
	}

	switch level {
	case storage.ERROR:
		levelString = "error"
	case storage.WARN:
		levelString = "warn"
	case storage.INFO:
		levelString = "info"
	case storage.TRACE:
		levelString = "trace"
	default:
		return errors.New("this log level is not implemented")
	}

	httpResp := c.post("/log/"+levelString, req)
	if httpResp.StatusCode == http.StatusNoContent {
		return nil
	}

	var jsonResp api.JsonResponse
	err := unmarshalResponse(httpResp, &jsonResp)

	if err != nil {
		return err
	}

	if !jsonResp.Ok {
		return errors.New(jsonResp.Message)
	}
	return nil
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
