package router

import (
	"fmt"
	"github.com/bernishen/exception"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
	t      *tree
	route  string
	config *Config
}

type Context struct {
	Token   string
	Params  Parameters
	Request *http.Request
}

// Default is get this route  default instance.
func InitRouter(c *Config) *LionRouter {
	return &LionRouter{
		t:      initTree(),
		route:  optimiizeRoute(c.RoutePrefix),
		config: c,
	}
}

func addRoute(r1 string, r2 string) string {
	rt1 := optimiizeRoute(r1)
	rt2 := optimiizeRoute(r2)
	if rt2 == "/" {
		return rt1
	}
	if rt1 == "/" {
		return rt2
	}
	return rt1 + rt2
}

func optimiizeRoute(r string) string {
	rt := strings.Trim(r, " ")
	if rt == "/" {
		return rt
	}
	if rt == "" {
		rt = "/"
		return rt
	}
	if rt[len(rt)-1:] == "/" {
		rt = rt[:len(rt)-1]
	}
	if rt[:1] != "/" {
		rt = "/" + rt
	}
	return rt
}

// ListenDefault is to start listen by this route default instance.
func ListenDefault(r *LionRouter) {
	server := r.config.Address + ":" + r.config.Port
	fmt.Println("\tLinsen address:" + server)
	log.Fatal(http.ListenAndServe(server, r))
}

func (r *LionRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.findAndRun(w, req)
}

func (r *LionRouter) findAndRun(w http.ResponseWriter, req *http.Request) {
	urlRoute := req.URL.Path
	c, params, ex := r.t.find(urlRoute)
	if ex != nil {
		w.WriteHeader(http.StatusNotFound)
		msg := []byte(ex.Message)
		_, err := w.Write(msg)
		if err != nil {
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
