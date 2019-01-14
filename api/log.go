package api

import (
	"errors"
	"github.com/Sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"src/task_tracker/config"
	"time"
)

type RequestHandler func(*Request)

type LogEntry struct {
	Scope     string `json:"scope"`
	Message   string `json:"Message"`
	TimeStamp int64  `json:"timestamp"`
}

func (e *LogEntry) Time() time.Time {

	t := time.Unix(e.TimeStamp, 0)
	return t
}

func LogRequest(h RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {

		logrus.WithFields(logrus.Fields{
			"path": string(ctx.Path()),
		}).Info(string(ctx.Method()))

		h(&Request{Ctx: ctx})
	})
}

func SetupLogger() {
	logrus.SetLevel(config.Cfg.LogLevel)
}

func parseLogEntry(r *Request) *LogEntry {

	entry := LogEntry{}

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
