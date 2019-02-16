package test

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"encoding/hex"
	"encoding/json"
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

func Post(path string, x interface{}, worker *storage.Worker, s *http.Client) *http.Response {

	if s == nil {
		s = &http.Client{}
	}

	body, err := json.Marshal(x)
	buf := bytes.NewBuffer(body)

	req, err := http.NewRequest("POST", "http://"+config.Cfg.ServerAddr+path, buf)
	handleErr(err)

	if worker != nil {
		mac := hmac.New(crypto.SHA256.New, worker.Secret)
		mac.Write(body)
		sig := hex.EncodeToString(mac.Sum(nil))

		req.Header.Add("X-Worker-Id", strconv.FormatInt(worker.Id, 10))
		req.Header.Add("X-Signature", sig)
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

		mac := hmac.New(crypto.SHA256.New, worker.Secret)
		mac.Write([]byte(path))
		sig := hex.EncodeToString(mac.Sum(nil))

		req.Header.Add("X-Worker-Id", strconv.FormatInt(worker.Id, 10))
		req.Header.Add("X-Signature", sig)
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
