package api

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"src/task_tracker/config"
	"src/task_tracker/storage"
)

type WebAPI struct {
	server   *fasthttp.Server
	router   *fasthttprouter.Router
	Database *storage.Database
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

func New() *WebAPI {

	api := new(WebAPI)
	api.Database = &storage.Database{}

	api.router = &fasthttprouter.Router{}

	api.server = &fasthttp.Server{
		Handler: api.router.Handler,
		Name:    info.Name,
	}

	api.router.GET("/", LogRequestMiddleware(Index))

	api.router.POST("/log/trace", LogRequestMiddleware(LogTrace))
	api.router.POST("/log/info", LogRequestMiddleware(LogInfo))
	api.router.POST("/log/warn", LogRequestMiddleware(LogWarn))
	api.router.POST("/log/error", LogRequestMiddleware(LogError))

	api.router.POST("/worker/create", LogRequestMiddleware(api.WorkerCreate))
	api.router.POST("/worker/update", LogRequestMiddleware(api.WorkerUpdate))
	api.router.GET("/worker/get/:id", LogRequestMiddleware(api.WorkerGet))

	api.router.POST("/access/grant", LogRequestMiddleware(api.WorkerGrantAccess))
	api.router.POST("/access/remove", LogRequestMiddleware(api.WorkerRemoveAccess))

	api.router.POST("/project/create", LogRequestMiddleware(api.ProjectCreate))
	api.router.GET("/project/get/:id", LogRequestMiddleware(api.ProjectGet))
	api.router.POST("/project/update/:id", LogRequestMiddleware(api.ProjectUpdate))
	api.router.GET("/project/stats/:id", LogRequestMiddleware(api.ProjectGetStats))
	api.router.GET("/project/stats", LogRequestMiddleware(api.ProjectGetAllStats))

	api.router.POST("/task/create", LogRequestMiddleware(api.TaskCreate))
	api.router.GET("/task/get/:project", LogRequestMiddleware(api.TaskGetFromProject))
	api.router.GET("/task/get", LogRequestMiddleware(api.TaskGet))
	api.router.POST("/task/release", LogRequestMiddleware(api.TaskRelease))

	api.router.POST("/git/receivehook", LogRequestMiddleware(api.ReceiveGitWebHook))

	api.router.POST("/logs", LogRequestMiddleware(api.GetLog))

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
