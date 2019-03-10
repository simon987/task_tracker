package api

import (
	"github.com/sirupsen/logrus"
	"runtime/debug"
)

func handleErr(err error, r *Request) {

	if err != nil {
		logrus.Error(err.Error())
		//debug.PrintStack()

		r.Json(JsonResponse{
			Message: err.Error(),
			Content: ErrorResponse{
				StackTrace: string(debug.Stack()),
			},
		}, 500)
	}
}
