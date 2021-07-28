package router

import (
	"fmt"
	"github.com/Berni-Shen/lion-go/utils/exception"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	HttpRequest_1 = iota
	WebSocket
)

type ILionController interface {
	FuncRun(method string, params *Parameters, req *http.Request) (*[]byte, *exception.Exception)
	Get(function interface{}, paramsorder ...string) *BaseController
	Post(function interface{}, paramsorder ...string) *BaseController
	Put(function interface{}, paramsorder ...string) *BaseController
	Delete(function interface{}, paramsorder ...string) *BaseController
	WebSocket(w http.ResponseWriter, req *http.Request, params *map[string]string)
}

type LionRouter struct {
	t *tree
}

type Context struct {
	Token   string
	Params  Parameters
	Request *http.Request
}

var defaultInstance = &LionRouter{t: initTree()}

// Default is get this router  default instance.
func Default() *LionRouter {
	return defaultInstance
}

// ListenDefault is to start listen by this router default instance.
func ListenDefault(config Config) {
	server := config.Address + ":" + config.Port
	fmt.Println("\tLinsen address:" + config.Address)
	log.Fatal(http.ListenAndServe(server, defaultInstance))
}

func (*LionRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	findAndRun(w, req)
}

func findAndRun(w http.ResponseWriter, req *http.Request) {
	urlRoute := req.URL.Path
	c, params, ex := defaultInstance.t.find(urlRoute)
	if ex != nil {
		msg := []byte(ex.Message)
		_, err := w.Write(msg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
		}
		return
	}

	findParamsFromRawQuery(req.URL.RawQuery, params)

	body, err := ioutil.ReadAll(req.Body)
	if err == nil && len(body) > 0 {
		str := string(body)
		params.AddParamData(BodyTag, str)
	}

	ret, ex := (*c).FuncRun(req.Method, params, req)
	if ex != nil {
		msg := []byte(ex.Message)
		_, err = w.Write(msg)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
		}
		return
	}
	_, err = w.Write(*ret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
	}
}

func findParamsFromRawQuery(rawQuery string, params *Parameters) {
	if strings.Trim(rawQuery, " ") == "" {
		return
	}
	queries := strings.Split(rawQuery, "&")
	for i := 0; i < len(queries); i++ {
		query := strings.Trim(queries[i], " ")
		if query == "" {
			continue
		}
		item := strings.Split(query, "=")
		if len(item) != 2 {
			continue
		}
		params.AddParamData(item[0], item[1])
	}
}
