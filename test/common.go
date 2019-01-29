package test

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"src/task_tracker/config"
	"src/task_tracker/storage"
)

func Post(path string, x interface{}, worker *storage.Worker) *http.Response {

	body, err := json.Marshal(x)
	buf := bytes.NewBuffer(body)

	req, err := http.NewRequest("POST", "http://"+config.Cfg.ServerAddr+path, buf)
	handleErr(err)

	if worker != nil {
		mac := hmac.New(crypto.SHA256.New, worker.Secret)
		mac.Write(body)
		sig := hex.EncodeToString(mac.Sum(nil))

		req.Header.Add("X-Worker-Id", worker.Id.String())
		req.Header.Add("X-Signature", sig)
	}

	client := http.Client{}
	r, err := client.Do(req)
	handleErr(err)

	return r
}

func Get(path string, worker *storage.Worker) *http.Response {

	url := "http://" + config.Cfg.ServerAddr + path
	req, err := http.NewRequest("GET", url, nil)
	handleErr(err)

	if worker != nil {

		fmt.Println(worker.Secret)
		mac := hmac.New(crypto.SHA256.New, worker.Secret)
		mac.Write([]byte(path))
		sig := hex.EncodeToString(mac.Sum(nil))

		req.Header.Add("X-Worker-Id", worker.Id.String())
		req.Header.Add("X-Signature", sig)
	}

	client := http.Client{}
	r, err := client.Do(req)
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
