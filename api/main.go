package api

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/buaazp/fasthttprouter"
	"github.com/kataras/go-sessions"
	"github.com/robfig/cron"
	"github.com/simon987/task_tracker/config"
	"github.com/simon987/task_tracker/storage"
	"github.com/valyala/fasthttp"
)

type WebAPI struct {
	server        *fasthttp.Server
	router        *fasthttprouter.Router
	Database      *storage.Database
	SessionConfig sessions.Config
	Session       *sessions.Sessions
	Cron          *cron.Cron
}

type Info struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

var info = Info{
	Name:    "task_tracker",
	Version: "1.0",
}

func Index(r *Request) {
	r.OkJson(info)
}

func (api *WebAPI) setupMonitoring() {

	api.Cron = cron.New()
	schedule := cron.Every(config.Cfg.MonitoringInterval)
	api.Cron.Schedule(schedule, cron.FuncJob(api.Database.MakeProjectSnapshots))
	api.Cron.Start()

	logrus.WithFields(logrus.Fields{
		"every": config.Cfg.MonitoringInterval.String(),
	}).Info("Started monitoring")
}

func New() *WebAPI {

	api := new(WebAPI)
	api.Database = &storage.Database{}
	api.setupMonitoring()

	api.router = &fasthttprouter.Router{}

	api.SessionConfig = sessions.Config{
		Cookie:                      config.Cfg.SessionCookieName,
		Expires:                     config.Cfg.SessionCookieExpiration,
		CookieSecureTLS:             false,
		DisableSubdomainPersistence: false,
	}

	api.Session = sessions.New(api.SessionConfig)

	api.server = &fasthttp.Server{
		Handler: api.router.Handler,
		Name:    info.Name,
	}

	api.router.GET("/", LogRequestMiddleware(Index))

	api.router.POST("/log/trace", LogRequestMiddleware(api.LogTrace))
	api.router.POST("/log/info", LogRequestMiddleware(api.LogInfo))
	api.router.POST("/log/warn", LogRequestMiddleware(api.LogWarn))
	api.router.POST("/log/error", LogRequestMiddleware(api.LogError))

	api.router.POST("/worker/create", LogRequestMiddleware(api.WorkerCreate))
	api.router.POST("/worker/update", LogRequestMiddleware(api.WorkerUpdate))
	api.router.GET("/worker/get/:id", LogRequestMiddleware(api.WorkerGet))
	api.router.GET("/worker/stats", LogRequestMiddleware(api.GetAllWorkerStats))

	api.router.POST("/access/grant", LogRequestMiddleware(api.WorkerGrantAccess))
	api.router.POST("/access/remove", LogRequestMiddleware(api.WorkerRemoveAccess))

	api.router.POST("/project/create", LogRequestMiddleware(api.ProjectCreate))
	api.router.GET("/project/get/:id", LogRequestMiddleware(api.ProjectGet))
	api.router.POST("/project/update/:id", LogRequestMiddleware(api.ProjectUpdate))
	api.router.GET("/project/list", LogRequestMiddleware(api.ProjectGetAllProjects))
	api.router.GET("/project/monitoring-between/:id", LogRequestMiddleware(api.GetSnapshotsBetween))
	api.router.GET("/project/monitoring/:id", LogRequestMiddleware(api.GetNSnapshots))
	api.router.GET("/project/assignees/:id", LogRequestMiddleware(api.ProjectGetAssigneeStats))
	api.router.GET("/project/requests/:id", LogRequestMiddleware(api.ProjectGetAccessRequests))
	api.router.GET("/project/request_access/:id", LogRequestMiddleware(api.WorkerRequestAccess))

	api.router.POST("/task/create", LogRequestMiddleware(api.TaskCreate))
	api.router.GET("/task/get/:project", LogRequestMiddleware(api.TaskGetFromProject))
	api.router.GET("/task/get", LogRequestMiddleware(api.TaskGet))
	api.router.POST("/task/release", LogRequestMiddleware(api.TaskRelease))

	api.router.POST("/git/receivehook", LogRequestMiddleware(api.ReceiveGitWebHook))

	api.router.POST("/logs", LogRequestMiddleware(api.GetLog))

	api.router.POST("/register", LogRequestMiddleware(api.Register))
	api.router.POST("/login", LogRequestMiddleware(api.Login))
	api.router.GET("/logout", LogRequestMiddleware(api.Logout))
	api.router.GET("/account", LogRequestMiddleware(api.AccountDetails))

	api.router.NotFound = func(ctx *fasthttp.RequestCtx) {

		if ctx.Request.Header.IsOptions() {
			ctx.Response.Header.Add("Access-Control-Allow-Headers", "Content-Type")
			ctx.Response.Header.Add("Access-Control-Allow-Methods", "GET, POST, OPTION")
			ctx.Response.Header.Add("Access-Control-Allow-Origin", "*")
		} else {
			ctx.SetStatusCode(404)
			_, _ = fmt.Fprintf(ctx, "Not found")
		}
	}

	return api
}

func (api *WebAPI) Run() {

	logrus.Infof("Started web api at address %s", config.Cfg.ServerAddr)

	err := api.server.ListenAndServe(config.Cfg.ServerAddr)
	if err != nil {
		logrus.Fatalf("Error in ListenAndServe: %s", err)
	}
}
