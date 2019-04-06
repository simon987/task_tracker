package api

import (
	"golang.org/x/time/rate"
	"time"
)

func (api *WebAPI) ReserveSubmit(pid int64, count int) *rate.Reservation {

	limiter, ok := api.SubmitLimiters.Load(pid)
	if !ok {
		project := api.Database.GetProject(pid)
		if project == nil {
			return nil
		}

		limiter = rate.NewLimiter(project.SubmitRate, 1)
		api.SubmitLimiters.Store(pid, limiter)
	}

	return limiter.(*rate.Limiter).ReserveN(time.Now(), count)
}

func (api *WebAPI) ReserveAssign(pid int64) *rate.Reservation {

	limiter, ok := api.AssignLimiters.Load(pid)
	if !ok {
		project := api.Database.GetProject(pid)
		if project == nil {
			return nil
		}

		limiter = rate.NewLimiter(project.AssignRate, 1)
		api.AssignLimiters.Store(pid, limiter)
	}

	return limiter.(*rate.Limiter).ReserveN(time.Now(), 1)
}
