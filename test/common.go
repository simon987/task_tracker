package test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"src/task_tracker/config"
)

func Post(path string, x interface{}) *http.Response {

	body, err := json.Marshal(x)
	buf := bytes.NewBuffer(body)

	r, err := http.Post("http://"+config.Cfg.ServerAddr+path, "application/json", buf)
	handleErr(err)

	return r
}

func Get(path string) *http.Response {
	r, err := http.Get("http://" + config.Cfg.ServerAddr + path)
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
