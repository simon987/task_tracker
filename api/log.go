package api

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"src/task_tracker/config"
	"src/task_tracker/storage"
	"time"
)

type RequestHandler func(*Request)

type GetLogRequest struct {
	Level logrus.Level `json:"level"`
	Since int64        `json:"since"`
}

type LogRequest struct {
	Scope     string `json:"scope"`
	Message   string `json:"Message"`
	TimeStamp int64  `json:"timestamp"`
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

func parseLogEntry(r *Request) *LogRequest {

	entry := LogRequest{}

	if r.GetJson(&entry) {
		if len(entry.Message) == 0 {
			handleErr(errors.New("invalid message"), r)
		} else if len(entry.Scope) == 0 {
			handleErr(errors.New("invalid scope"), r)
		} else if entry.TimeStamp <= 0 {
			handleErr(errors.New("invalid timestamp"), r)
		}
	}

	return &entry
}

func LogTrace(r *Request) {

	entry := parseLogEntry(r)

	logrus.WithFields(logrus.Fields{
		"scope": entry.Scope,
	}).WithTime(entry.Time()).Trace(entry.Message)
}

func LogInfo(r *Request) {

	entry := parseLogEntry(r)

	logrus.WithFields(logrus.Fields{
		"scope": entry.Scope,
	}).WithTime(entry.Time()).Info(entry.Message)
}

func LogWarn(r *Request) {

	entry := parseLogEntry(r)

	logrus.WithFields(logrus.Fields{
		"scope": entry.Scope,
	}).WithTime(entry.Time()).Warn(entry.Message)
}

func LogError(r *Request) {

	entry := parseLogEntry(r)

	logrus.WithFields(logrus.Fields{
		"scope": entry.Scope,
	}).WithTime(entry.Time()).Error(entry.Message)
}

func (api *WebAPI) GetLog(r *Request) {

	req := &GetLogRequest{}
	if r.GetJson(req) {
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
}

func (r GetLogRequest) isValid() bool {

	if r.Since <= 0 {
		return false
	}

	return true
}
