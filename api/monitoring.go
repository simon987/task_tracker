package api

import (
	"math"
	"strconv"
)

func (api *WebAPI) GetSnapshotsWithinRange(r *Request) {

	id, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	from := r.Ctx.Request.URI().QueryArgs().GetUintOrZero("from")
	to := r.Ctx.Request.URI().QueryArgs().GetUintOrZero("to")
	if err != nil || id <= 0 || from <= 0 || to <= 0 || from >= math.MaxInt32 || to >= math.MaxInt32 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid request",
		}, 400)
		return
	}

	snapshots := api.Database.GetMonitoringSnapshotsBetween(id, from, to)
	r.OkJson(JsonResponse{
		Ok: true,
		Content: GetSnapshotsResponse{
			Snapshots: snapshots,
		},
	})
}

func (api *WebAPI) GetNSnapshots(r *Request) {

	id, err := strconv.ParseInt(r.Ctx.UserValue("id").(string), 10, 64)
	count := r.Ctx.Request.URI().QueryArgs().GetUintOrZero("count")
	if err != nil || id <= 0 || count <= 0 || count >= 1000 {
		r.Json(JsonResponse{
			Ok:      false,
			Message: "Invalid request",
		}, 400)
		return
	}

	snapshots := api.Database.GetNMonitoringSnapshots(id, count)
	r.OkJson(JsonResponse{
		Ok: true,
		Content: GetSnapshotsResponse{
			Snapshots: snapshots,
		},
	})
}
