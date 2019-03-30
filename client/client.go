package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/simon987/task_tracker/api"
	"io/ioutil"
	"net/http"
	"strconv"
)

type taskTrackerClient struct {
	worker        *Worker
	httpClient    http.Client
	serverAddress string

	secretB64   string
	workerIdStr string
}

func New(serverAddress string) *taskTrackerClient {

	client := new(taskTrackerClient)
	client.serverAddress = serverAddress

	return client
}

func (c *taskTrackerClient) SetWorker(worker *Worker) {
	c.worker = worker
	c.secretB64 = base64.StdEncoding.EncodeToString(worker.Secret)
	c.workerIdStr = strconv.FormatInt(worker.Id, 10)
}

func (c *taskTrackerClient) get(path string) *http.Response {

	url := "http://" + c.serverAddress + path
	req, err := http.NewRequest("GET", url, nil)

	if c.worker != nil {
		req.Header.Add("X-Worker-Id", c.workerIdStr)
		req.Header.Add("X-Secret", c.secretB64)
	}

	r, err := c.httpClient.Do(req)
	handleErr(err)

	return r
}

func (c *taskTrackerClient) post(path string, x interface{}) *http.Response {

	body, err := json.Marshal(x)
	buf := bytes.NewBuffer(body)

	req, err := http.NewRequest("POST", "http://"+c.serverAddress+path, buf)
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
	fmt.Println(string(data))
	err = json.Unmarshal(data, result)
	if err != nil {
		return err
	}
	return nil
}

func (c taskTrackerClient) MakeWorker(alias string) (*Worker, error) {

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

func (c taskTrackerClient) FetchTask(projectId int) (*AssignTaskResponse, error) {

	httpResp := c.get("/task/get/" + strconv.Itoa(projectId))
	var jsonResp AssignTaskResponse
	err := unmarshalResponse(httpResp, &jsonResp)

	//TODO: Handle rate limiting here?

	return &jsonResp, err
}

func (c taskTrackerClient) ReleaseTask(req api.ReleaseTaskRequest) (*ReleaseTaskResponse, error) {

	httpResp := c.post("/task/release", req)
	var jsonResp ReleaseTaskResponse
	err := unmarshalResponse(httpResp, &jsonResp)

	return &jsonResp, err
}

func (c taskTrackerClient) SubmitTask(req api.SubmitTaskRequest) (AssignTaskResponse, error) {

	httpResp := c.post("/task/submit", req)
	var jsonResp AssignTaskResponse
	err := unmarshalResponse(httpResp, &jsonResp)

	//TODO: Handle rate limiting here?

	return jsonResp, err
}

func (c taskTrackerClient) GetProjectSecret(projectId int) (string, error) {

	httpResp := c.get("/project/secret/" + strconv.Itoa(projectId))
	var jsonResp ProjectSecretResponse
	err := unmarshalResponse(httpResp, &jsonResp)

	return jsonResp.Content.Secret, err
}

func (c taskTrackerClient) RequestAccess(req api.CreateWorkerAccessRequest) (api.JsonResponse, error) {

	httpResp := c.post("/project/request_access", req)
	var jsonResp api.JsonResponse
	err := unmarshalResponse(httpResp, &jsonResp)

	return jsonResp, err
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
