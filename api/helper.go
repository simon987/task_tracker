package api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Request struct {
	Ctx *fasthttp.RequestCtx
}

func (r *Request) OkJson(object JsonResponse) {

	resp, err := json.Marshal(object)
	handleErr(err, r)

	r.Ctx.Response.Header.Set("Content-Type", "application/json")
	_, err = r.Ctx.Write(resp)
	handleErr(err, r)
}

func (r *Request) Ok() {
	r.Ctx.Response.SetStatusCode(204)
}

func (r *Request) Json(object JsonResponse, code int) {

	resp, err := json.Marshal(object)
	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"code": code,
		}).Error("Error during json encoding of object")
		return
	}

	r.Ctx.Response.SetStatusCode(code)
	r.Ctx.Response.Header.Set("Content-Type", "application/json")
	_, err = r.Ctx.Write(resp)
	if err != nil {
		panic(err)
	}

}
