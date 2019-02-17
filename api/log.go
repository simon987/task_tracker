package api

import (
	"encoding/json"
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/simon987/task_tracker/config"
	"github.com/valyala/fasthttp"
	"time"
)

func (e *LogRequest) Time() time.Time {

	t := time.Unix(e.TimeStamp, 0)
	return t
}

func LogRequestMiddleware(h RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {

		ctx.Response.Header.Add("Access-Control-Allow-Headers", "Content-Type")
		ctx.Response.Header.Add("Access-Control-Allow-Methods", "GET, POST, OPTION")
		ctx.Response.Header.Add("Access-Control-Allow-Origin", "*")

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
		r.Json(JsonResponse{
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
		r.Json(JsonResponse{
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
		r.Json(JsonResponse{
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
		r.Json(JsonResponse{
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
		r.Json(JsonResponse{
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

		r.OkJson(JsonResponse{
			Ok: true,
			Content: GetLogResponse{
				Logs: logs,
			},
		})
	} else {
		logrus.WithFields(logrus.Fields{
			"getLogRequest": req,
		}).Warn("Invalid log request")

		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid log request",
		}, 400)
	}
}
