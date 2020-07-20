package tinyapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/andriiyaremenko/tinyapi/utils"
)

type HandlerFunc func(w http.ResponseWriter, req *http.Request, params map[string]string)

type Endpoint interface {
	http.Handler

	Path() string
	PrependPath(prefix string)

	NotFound(handler http.HandlerFunc)
	Middleware(handlers ...http.HandlerFunc)

	Handle(method string, param string, handler HandlerFunc)
	Get(param string, handler HandlerFunc)
	Post(param string, handler HandlerFunc)
	Put(param string, handler HandlerFunc)
	Patch(param string, handler HandlerFunc)
	Delete(param string, handler HandlerFunc)
}

func CombineEndpoints(path string, endpoints ...Endpoint) http.Handler {
	api := http.NewServeMux()

	for _, e := range endpoints {
		e.PrependPath(path)
		api.Handle(e.Path(), e)
	}

	return api
}

type Configure func(e Endpoint) Endpoint

func NewEndpoint(path string, configure Configure) Endpoint {
	if path[0] != '/' {
		path = "/" + path
	}

	if path[len(path)-1] != '/' {
		path = path + "/"
	}

	return configure(
		&endpoint{
			path:     path,
			routes:   make(map[string]map[string]HandlerFunc),
			notFound: utils.NotFound,
		},
	)
}

type endpoint struct {
	path       string
	routes     map[string]map[string]HandlerFunc
	notFound   http.HandlerFunc
	middleware []http.HandlerFunc
}

func (e *endpoint) Path() string {
	return e.path
}

func (e *endpoint) PrependPath(prefix string) {
	if prefix == "/" {
		return
	}

	if prefix[0] != '/' {
		prefix = "/" + prefix
	}

	e.path = prefix + e.path
}

func (e *endpoint) Handle(method string, param string, handler HandlerFunc) {
	if param == "" {
		param = "/"
	}

	if len(param) > 1 && param[0] == '/' {
		param = param[1:]
	}

	if len(param) > 1 && param[len(param)-1] == '/' {
		param = param[:len(param)-2]
	}

	_, ok := e.routes[method]
	if !ok {
		e.routes[method] = make(map[string]HandlerFunc)
	}

	if _, ok := e.routes[method][param]; ok {
		panic(fmt.Errorf("handler for %s %s already exists", method, param))
	}

	e.routes[method][param] = handler
}

func (e *endpoint) Get(param string, handler HandlerFunc) {
	e.Handle(http.MethodGet, param, handler)
}

func (e *endpoint) Post(param string, handler HandlerFunc) {
	e.Handle(http.MethodPost, param, handler)
}

func (e *endpoint) Put(param string, handler HandlerFunc) {
	e.Handle(http.MethodPut, param, handler)
}

func (e *endpoint) Patch(param string, handler HandlerFunc) {
	e.Handle(http.MethodPatch, param, handler)
}

func (e *endpoint) Delete(param string, handler HandlerFunc) {
	e.Handle(http.MethodDelete, param, handler)
}

func (e *endpoint) Middleware(handlers ...http.HandlerFunc) {
	e.middleware = append(e.middleware, handlers...)
}

func (e *endpoint) NotFound(handler http.HandlerFunc) {
	e.notFound = handler
}

func (e *endpoint) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, h := range e.middleware {
		h(w, req)
	}
	r, ok := e.routes[req.Method]
	if !ok {
		e.notFound(w, req)
		return
	}
	reqParamStr := e.requestParam(req)

	reqParams := strings.Split(reqParamStr, "/")
	reqParamsLen := len(reqParams)

	for k, v := range r {
		routeParams := strings.Split(k, "/")
		if len(routeParams) != reqParamsLen {
			continue
		}

		param, ok := getParams(routeParams, reqParams)

		if !ok {
			continue
		}

		v(w, req, param)
		return
	}
	e.notFound(w, req)
}

func (e *endpoint) requestParam(req *http.Request) string {
	result := req.URL.RequestURI()[len(req.URL.Host+e.path):]

	if result == "" {
		return "/"
	}

	return result
}

func getParams(routeParams, reqParams []string) (map[string]string, bool) {
	params := make(map[string]string)
	for i, rp := range routeParams {
		reqP := reqParams[i]
		if rp[0] != ':' && rp != reqP {
			return nil, false
		}
		params[rp[1:]] = reqP
	}
	return params, true
}
