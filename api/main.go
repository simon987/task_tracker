package api

import (
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

	SetupLogger()

	api := new(WebAPI)

	api.router = &fasthttprouter.Router{}

	api.server = &fasthttp.Server{
		Handler: api.router.Handler,
		Name:    info.Name,
	}

	api.router.GET("/", LogRequest(Index))

	api.router.POST("/log/trace", LogRequest(LogTrace))
	api.router.POST("/log/info", LogRequest(LogInfo))
	api.router.POST("/log/warn", LogRequest(LogWarn))
	api.router.POST("/log/error", LogRequest(LogError))

	api.router.POST("/worker/create", LogRequest(api.WorkerCreate))
	api.router.GET("/worker/get/:id", LogRequest(api.WorkerGet))

	api.router.POST("/project/create", LogRequest(api.ProjectCreate))
	api.router.GET("/project/get/:id", LogRequest(api.ProjectGet))

	api.router.POST("/task/create", LogRequest(api.TaskCreate))
	api.router.GET("/task/get/:project", LogRequest(api.TaskGetFromProject))
	api.router.GET("/task/get", LogRequest(api.TaskGet))

	return api
}

func (api *WebAPI) Run() {

	logrus.Infof("Started web api at address %s", config.Cfg.ServerAddr)

	err := api.server.ListenAndServe(config.Cfg.ServerAddr)
	if err != nil {
		logrus.Fatalf("Error in ListenAndServe: %s", err)
	}
}
