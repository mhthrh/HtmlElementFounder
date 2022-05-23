package Model

import (
	"GitHub.com/mhthrh/HtmlElementsFinder/Helper/JsonUtil"
	"net/http"
	"time"
)

type Request struct {
	ThreadCount int    `json:"threadCount"validate:"required,numeric,gt=0"`
	Address     string `json:"address"validate:"required,url"`
}

type Response struct {
	Header struct {
		HeaderStatus int `json:"header-status"`
	} `json:"header"`

	Body struct {
		BodyStatus int         `json:"body-status"`
		Time       time.Time   `json:"time"`
		Result     interface{} `json:"result"`
	} `json:"body"`
}

func NewResult(BStatus, HStatus int, result interface{}) *Response {
	r := new(Response)
	r.Body.Time = time.Now()
	r.Header.HeaderStatus = HStatus
	r.Body.BodyStatus = BStatus
	r.Body.Result = result
	return r
}
func (r *Response) SendResponse(w http.ResponseWriter) {
	w.Write([]byte(JsonUtil.New(nil, nil).Struct2Json(r.Body)))
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(r.Header.HeaderStatus)
}
