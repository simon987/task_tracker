package api

import (
	"encoding/json"
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"src/task_tracker/config"
	"src/task_tracker/storage"
	"time"
)

type RequestHandler func(*Request)

type GetLogRequest struct {
	Level storage.LogLevel `json:"level"`
	Since int64            `json:"since"`
}

type LogRequest struct {
	Scope     string `json:"scope"`
	Message   string `json:"Message"`
	TimeStamp int64  `json:"timestamp"`
	worker    *storage.Worker
}

type GetLogResponse struct {
	Ok      bool                `json:"ok"`
	Message string              `json:"message"`
	Logs    *[]storage.LogEntry `json:"logs"`
}

func (e *LogRequest) Time() time.Time {

	t := time.Unix(e.TimeStamp, 0)
	return t
}

func LogRequestMiddleware(h RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {

		logrus.WithFields(logrus.Fields{
			"path":   string(ctx.Path()),
			"header": ctx.Request.Header.String(),
		}).Trace(string(ctx.Method()))

		h(&Request{Ctx: ctx})
	})
}

func (api *WebAPI) SetupLogger() {
	logrus.SetLevel(config.Cfg.LogLevel)

	api.Database.SetupLoggerHook()
}

func (api *WebAPI) parseLogEntry(r *Request) (*LogRequest, error) {

	worker, err := api.validateSignature(r)
	if err != nil {
		return nil, err
	}

	entry := LogRequest{}

	err = json.Unmarshal(r.Ctx.Request.Body(), &entry)
	if err != nil {
		return nil, err
	}

	if len(entry.Message) == 0 {
		return nil, errors.New("invalid message")
	} else if len(entry.Scope) == 0 {
		return nil, errors.New("invalid scope")
	} else if entry.TimeStamp <= 0 {
		return nil, errors.New("invalid timestamp")
	}

	entry.worker = worker

	return &entry, nil
}

func (api *WebAPI) LogTrace(r *Request) {

	entry, err := api.parseLogEntry(r)
	if err != nil {
		r.Json(GetLogResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}

	logrus.WithFields(logrus.Fields{
		"scope":  entry.Scope,
		"worker": entry.worker.Id,
	}).WithTime(entry.Time()).Trace(entry.Message)
}

func (api *WebAPI) LogInfo(r *Request) {

	entry, err := api.parseLogEntry(r)
	if err != nil {
		r.Json(GetLogResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}

	logrus.WithFields(logrus.Fields{
		"scope":  entry.Scope,
		"worker": entry.worker.Id,
	}).WithTime(entry.Time()).Info(entry.Message)
}

func (api *WebAPI) LogWarn(r *Request) {

	entry, err := api.parseLogEntry(r)
	if err != nil {
		r.Json(GetLogResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}

	logrus.WithFields(logrus.Fields{
		"scope":  entry.Scope,
		"worker": entry.worker.Id,
	}).WithTime(entry.Time()).Warn(entry.Message)
}

func (api *WebAPI) LogError(r *Request) {

	entry, err := api.parseLogEntry(r)
	if err != nil {
		r.Json(GetLogResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}

	logrus.WithFields(logrus.Fields{
		"scope":  entry.Scope,
		"worker": entry.worker.Id,
	}).WithTime(entry.Time()).Error(entry.Message)
}

func (api *WebAPI) GetLog(r *Request) {

	req := &GetLogRequest{}
	err := json.Unmarshal(r.Ctx.Request.Body(), req)
	if err != nil {
		r.Json(GetLogResponse{
			Ok:      false,
			Message: "Could not parse request",
		}, 400)
		return
	}
	if req.isValid() {

		logs := api.Database.GetLogs(req.Since, req.Level)

		logrus.WithFields(logrus.Fields{
			"getLogRequest": req,
			"logCount":      len(*logs),
		}).Trace("Get log request")

		r.OkJson(GetLogResponse{
			Ok:   true,
			Logs: logs,
		})
	} else {
		logrus.WithFields(logrus.Fields{
			"getLogRequest": req,
		}).Warn("Invalid log request")

		r.Json(GetLogResponse{
			Ok:      false,
			Message: "Invalid log request",
		}, 400)
	}
}

func (r GetLogRequest) isValid() bool {

	if r.Since <= 0 {
		return false
	}

	return true
}
