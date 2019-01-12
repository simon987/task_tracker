package api

import (
	"github.com/Sirupsen/logrus"
	"runtime/debug"
)

type ErrorResponse struct {
	Message    string `json:"message"`
	StackTrace string `json:"stack_trace"`

}

func handleErr(err error, r *Request) {

	if err != nil {
		logrus.Error(err.Error())
		//debug.PrintStack()

		r.Json(ErrorResponse{
			Message: err.Error(),
			StackTrace: string(debug.Stack()),
		}, 500)
	}
}
