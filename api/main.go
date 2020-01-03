package api

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/fasthttp/session"
	"github.com/fasthttp/session/memory"
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
	SessionConfig  *session.Config
	Session        *session.Session
	Cron           *cron.Cron
	AssignLimiters sync.Map
	SubmitLimiters sync.Map
}

type RequestHandler func(*Request)

var info = Info{
	Name:    "task_tracker",
	Version: "1.1",
}

func Middleware(h RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		h(&Request{Ctx: ctx})
	})
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

	api.SessionConfig = &session.Config{
		CookieName: config.Cfg.SessionCookieName,
		Expires:    config.Cfg.SessionCookieExpiration,
		Secure:     false,
	}

	api.Session = session.New(api.SessionConfig)
	_ = api.Session.SetProvider("memory", &memory.Config{})

	api.server = &fasthttp.Server{
		Handler: api.router.Handler,
		Name:    info.Name,
	}

	api.router.GET("/", Middleware(Index))

	api.router.POST("/log/trace", Middleware(api.LogTrace))
	api.router.POST("/log/info", Middleware(api.LogInfo))
	api.router.POST("/log/warn", Middleware(api.LogWarn))
	api.router.POST("/log/error", Middleware(api.LogError))

	api.router.POST("/worker/create", Middleware(api.CreateWorker))
	api.router.POST("/worker/update", Middleware(api.UpdateWorker))
	api.router.POST("/worker/set_paused", Middleware(api.WorkerSetPaused))
	api.router.GET("/worker/get/:id", Middleware(api.GetWorker))
	api.router.GET("/worker/stats", Middleware(api.GetAllWorkerStats))

	api.router.POST("/project/create", Middleware(api.CreateProject))
	api.router.GET("/project/get/:id", Middleware(api.GetProject))
	api.router.POST("/project/update/:id", Middleware(api.UpdateProject))
	api.router.GET("/project/list", Middleware(api.GetProjectList))
	api.router.GET("/project/monitoring-between/:id", Middleware(api.GetSnapshotsWithinRange))
	api.router.GET("/project/monitoring/:id", Middleware(api.GetNSnapshots))
	api.router.GET("/project/assignees/:id", Middleware(api.GetAssigneeStatsForProject))
	api.router.GET("/project/access_list/:id", Middleware(api.GetWorkerAccessListForProject))
	api.router.POST("/project/request_access", Middleware(api.CreateWorkerAccess))
	api.router.POST("/project/accept_request/:id/:wid", Middleware(api.AcceptAccessRequest))
	api.router.POST("/project/reject_request/:id/:wid", Middleware(api.RejectAccessRequest))
	api.router.GET("/project/secret/:id", Middleware(api.GetSecret))
	api.router.POST("/project/secret/:id", Middleware(api.SetSecret))
	api.router.GET("/project/webhook_secret/:id", Middleware(api.GetWebhookSecret))
	api.router.POST("/project/webhook_secret/:id", Middleware(api.SetWebhookSecret))
	api.router.POST("/project/reset_failed_tasks/:id", Middleware(api.ResetFailedTasks))
	api.router.POST("/project/hard_reset/:id", Middleware(api.HardReset))
	api.router.POST("/project/reclaim_assigned_tasks/:id", Middleware(api.ReclaimAssignedTasks))

	api.router.POST("/task/submit", Middleware(api.SubmitTask))
	api.router.POST("/task/bulk_submit", Middleware(api.BulkSubmitTask))
	api.router.GET("/task/get/:project", Middleware(api.GetTaskFromProject))
	api.router.POST("/task/release", Middleware(api.ReleaseTask))

	api.router.POST("/git/receivehook", Middleware(api.ReceiveGitWebHook))

	api.router.POST("/logs", Middleware(api.GetLog))

	api.router.POST("/register", Middleware(api.Register))
	api.router.POST("/login", Middleware(api.Login))
	api.router.GET("/logout", Middleware(api.Logout))
	api.router.GET("/account", Middleware(api.GetAccountDetails))
	api.router.GET("/manager/list", Middleware(api.GetManagerList))
	api.router.GET("/manager/list_for_project/:id", Middleware(api.GetManagerListWithRoleOn))
	api.router.GET("/manager/promote/:id", Middleware(api.PromoteManager))
	api.router.GET("/manager/demote/:id", Middleware(api.DemoteManager))
	api.router.POST("/manager/set_role_for_project/:id", Middleware(api.SetManagerRoleOnProject))

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
