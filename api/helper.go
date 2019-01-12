package api

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

type Request struct {
	Ctx *fasthttp.RequestCtx
}


func (r *Request) OkJson(object interface{}) {

	resp, err := json.Marshal(object)
	handleErr(err, r)

	r.Ctx.Response.Header.Set("Content-Type", "application/json")
	_, err = r.Ctx.Write(resp)
	handleErr(err, r)
}

func (r *Request) Json(object interface{}, code int) {

	resp, err := json.Marshal(object)
	if err != nil {
		fmt.Fprint(r.Ctx,"Error during json encoding of error")
		logrus.Error("Error during json encoding of error")
	}

	r.Ctx.Response.SetStatusCode(code)
	r.Ctx.Response.Header.Set("Content-Type", "application/json")
	_, err = r.Ctx.Write(resp)
	if err != nil {
		panic(err) //todo handle differently
	}

}

func (r *Request) GetJson(x interface{}) bool {

	err := json.Unmarshal(r.Ctx.Request.Body(), x)
	handleErr(err, r)

	return err == nil
}
