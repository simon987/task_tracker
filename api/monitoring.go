package api

import (
	"github.com/simon987/task_tracker/storage"
	"math"
	"strconv"
)

type MonitoringSnapshotResponse struct {
	Ok        bool                                 `json:"ok"`
	Message   string                               `json:"message,omitempty"`
	Snapshots *[]storage.ProjectMonitoringSnapshot `json:"snapshots,omitempty"`
}

func (api *WebAPI) GetSnapshotsBetween(r *Request) {

	id, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	from := r.Ctx.Request.URI().QueryArgs().GetUintOrZero("from")
	to := r.Ctx.Request.URI().QueryArgs().GetUintOrZero("to")
	if err != nil || id <= 0 || from <= 0 || to <= 0 || from >= math.MaxInt32 || to >= math.MaxInt32 {
		r.Json(MonitoringSnapshotResponse{
			Ok:      false,
			Message: "Invalid request",
		}, 400)
		return
	}

	snapshots := api.Database.GetMonitoringSnapshotsBetween(id, from, to)
	r.OkJson(MonitoringSnapshotResponse{
		Ok:        true,
		Snapshots: snapshots,
	})
}

func (api *WebAPI) GetNSnapshots(r *Request) {

	id, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	count := r.Ctx.Request.URI().QueryArgs().GetUintOrZero("count")
	if err != nil || id <= 0 || count <= 0 || count >= 1000 {
		r.Json(MonitoringSnapshotResponse{
			Ok:      false,
			Message: "Invalid request",
		}, 400)
		return
	}

	snapshots := api.Database.GetNMonitoringSnapshots(id, count)
	r.OkJson(MonitoringSnapshotResponse{
		Ok:        true,
		Snapshots: snapshots,
	})
}
