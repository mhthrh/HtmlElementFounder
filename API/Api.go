package API

import (
	"GitHub.com/mhthrh/HtmlElementsFinder/Controler"
	"GitHub.com/mhthrh/HtmlElementsFinder/Helper/JsonUtil"
	"GitHub.com/mhthrh/HtmlElementsFinder/Helper/LogUtil"
	"GitHub.com/mhthrh/HtmlElementsFinder/Helper/Validation"
	"GitHub.com/mhthrh/HtmlElementsFinder/Model"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

var (
	mux map[string]func(http.ResponseWriter, *http.Request)
	l   *logrus.Entry
	v   *Validation.Validation
)

func init() {
	l = LogUtil.NewLogger()
	v = Validation.NewValidation()
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/page"] = func(writer http.ResponseWriter, request *http.Request) {
		path := request.URL.Path

		if path == "/page" {
			path = "./Site/index.html"
		}
		http.ServeFile(writer, request, path)
		return
	}
	mux["/find"] = func(writer http.ResponseWriter, request *http.Request) {
		l.Println("newRequest")
		Inp := Model.Request{}
		err := json.NewDecoder(request.Body).Decode(&Inp)

		if err != nil {
			Model.NewResult(1002, http.StatusInternalServerError, err).SendResponse(writer)
			return
		}
		errs := v.Validate(Inp)
		if len(errs) != 0 {
			Model.NewResult(1003, http.StatusUnavailableForLegalReasons, errs).SendResponse(writer)
			return
		}
		result, err := Controler.New(l).HtmlElementInfo(Inp)

		Model.NewResult(1000, http.StatusOK, JsonUtil.New(nil, nil).Struct2Json(result)).SendResponse(writer)

	}
}

type RequestHandler struct{}

func (*RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if method, ok := mux[r.URL.Path]; ok {
		method(w, r)
		return
	}
	Model.NewResult(1004, http.StatusNotFound, "Page Not found!").SendResponse(w)
}
