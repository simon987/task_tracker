package api

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/kataras/go-sessions"
	"github.com/robfig/cron"
	"github.com/simon987/task_tracker/config"
	"github.com/simon987/task_tracker/storage"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
	"sync"
)

type WebAPI struct {
	server         *fasthttp.Server
	router         *fasthttprouter.Router
	Database       *storage.Database
	SessionConfig  sessions.Config
	Session        *sessions.Sessions
	Cron           *cron.Cron
	AssignLimiters sync.Map
	SubmitLimiters sync.Map
}

type RequestHandler func(*Request)

var info = Info{
	Name:    "task_tracker",
	Version: "1.0",
}

func Index(r *Request) {
	r.OkJson(JsonResponse{
		Ok:      true,
		Content: info,
	})
}

func (api *WebAPI) setupMonitoring() {

	api.Cron = cron.New()
	monSchedule := cron.Every(config.Cfg.MonitoringInterval)
	api.Cron.Schedule(monSchedule, cron.FuncJob(api.Database.MakeProjectSnapshots))

	timeoutSchedule := cron.Every(config.Cfg.ResetTimedOutTasksInterval)
	api.Cron.Schedule(timeoutSchedule, cron.FuncJob(api.Database.ResetTimedOutTasks))
	api.Cron.Start()

	logrus.WithFields(logrus.Fields{
		"every": config.Cfg.MonitoringInterval.String(),
	}).Info("Started monitoring")
	logrus.WithFields(logrus.Fields{
		"every": config.Cfg.ResetTimedOutTasksInterval.String(),
	}).Info("Started task cleanup cron")
}

func New() *WebAPI {

	api := new(WebAPI)
	api.Database = storage.New()
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

	api.router.POST("/worker/create", LogRequestMiddleware(api.CreateWorker))
	api.router.POST("/worker/update", LogRequestMiddleware(api.UpdateWorker))
	api.router.GET("/worker/get/:id", LogRequestMiddleware(api.GetWorker))
	api.router.GET("/worker/stats", LogRequestMiddleware(api.GetAllWorkerStats))

	api.router.POST("/project/create", LogRequestMiddleware(api.CreateProject))
	api.router.GET("/project/get/:id", LogRequestMiddleware(api.GetProject))
	api.router.POST("/project/update/:id", LogRequestMiddleware(api.UpdateProject))
	api.router.GET("/project/list", LogRequestMiddleware(api.GetProjectList))
	api.router.GET("/project/monitoring-between/:id", LogRequestMiddleware(api.GetSnapshotsWithinRange))
	api.router.GET("/project/monitoring/:id", LogRequestMiddleware(api.GetNSnapshots))
	api.router.GET("/project/assignees/:id", LogRequestMiddleware(api.GetAssigneeStatsForProject))
	api.router.GET("/project/access_list/:id", LogRequestMiddleware(api.GetWorkerAccessListForProject))
	api.router.POST("/project/request_access", LogRequestMiddleware(api.CreateWorkerAccess))
	api.router.POST("/project/accept_request/:id/:wid", LogRequestMiddleware(api.AcceptAccessRequest))
	api.router.POST("/project/reject_request/:id/:wid", LogRequestMiddleware(api.RejectAccessRequest))
	api.router.GET("/project/secret/:id", LogRequestMiddleware(api.GetSecret))
	api.router.POST("/project/secret/:id", LogRequestMiddleware(api.SetSecret))
	api.router.GET("/project/webhook_secret/:id", LogRequestMiddleware(api.GetWebhookSecret))
	api.router.POST("/project/webhook_secret/:id", LogRequestMiddleware(api.SetWebhookSecret))
	api.router.POST("/project/reset_failed_tasks/:id", LogRequestMiddleware(api.ResetFailedTasks))
	api.router.POST("/project/hard_reset/:id", LogRequestMiddleware(api.HardReset))

	api.router.POST("/task/submit", LogRequestMiddleware(api.SubmitTask))
	api.router.GET("/task/get/:project", LogRequestMiddleware(api.GetTaskFromProject))
	api.router.POST("/task/release", LogRequestMiddleware(api.ReleaseTask))

	api.router.POST("/git/receivehook", LogRequestMiddleware(api.ReceiveGitWebHook))

	api.router.POST("/logs", LogRequestMiddleware(api.GetLog))

	api.router.POST("/register", LogRequestMiddleware(api.Register))
	api.router.POST("/login", LogRequestMiddleware(api.Login))
	api.router.GET("/logout", LogRequestMiddleware(api.Logout))
	api.router.GET("/account", LogRequestMiddleware(api.GetAccountDetails))
	api.router.GET("/manager/list", LogRequestMiddleware(api.GetManagerList))
	api.router.GET("/manager/list_for_project/:id", LogRequestMiddleware(api.GetManagerListWithRoleOn))
	api.router.GET("/manager/promote/:id", LogRequestMiddleware(api.PromoteManager))
	api.router.GET("/manager/demote/:id", LogRequestMiddleware(api.DemoteManager))
	api.router.POST("/manager/set_role_for_project/:id", LogRequestMiddleware(api.SetManagerRoleOnProject))

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
