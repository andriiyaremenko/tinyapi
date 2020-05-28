package tinyapi

import (
	"net/http"
	"regexp"

	"github.com/andriiyaremenko/tinyapi/utils"
)

type Endpoint interface {
	http.Handler
	NotFound(handler http.HandlerFunc)
	Handle(method string, param string, handler http.HandlerFunc)
	Get(param string, handler http.HandlerFunc)
	Post(param string, handler http.HandlerFunc)
	Put(param string, handler http.HandlerFunc)
	Patch(param string, handler http.HandlerFunc)
	Delete(param string, handler http.HandlerFunc)
	Middleware(handlers ...http.HandlerFunc)
}

type endpoint struct {
	notFound   http.HandlerFunc
	routes     map[string]map[*regexp.Regexp]http.HandlerFunc
	middleware []http.HandlerFunc
}

func (e *endpoint) Handle(method string, param string, handler http.HandlerFunc) {
	if param == "" {
		param = "/"
	}
	_, ok := e.routes[method]
	if !ok {
		e.routes[method] = make(map[*regexp.Regexp]http.HandlerFunc)
	}
	e.routes[method][regexp.MustCompile(param)] = handler
}

func (e *endpoint) Get(param string, handler http.HandlerFunc) {
	e.Handle(http.MethodGet, param, handler)
}

func (e *endpoint) Post(param string, handler http.HandlerFunc) {
	e.Handle(http.MethodPost, param, handler)
}

func (e *endpoint) Put(param string, handler http.HandlerFunc) {
	e.Handle(http.MethodPut, param, handler)
}

func (e *endpoint) Patch(param string, handler http.HandlerFunc) {
	e.Handle(http.MethodPatch, param, handler)
}

func (e *endpoint) Delete(param string, handler http.HandlerFunc) {
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
	param, ok := utils.RequestParam(req)

	if !ok {
		param = "/"
	}

	for k, v := range r {
		if k.MatchString(param) {
			v(w, req)
			return
		}
	}
	e.notFound(w, req)
}

func NewNilEndpoint() Endpoint {
	return &endpoint{routes: make(map[string]map[*regexp.Regexp]http.HandlerFunc), notFound: utils.NotFound}
}

func NewEndpoint(configure func(e Endpoint)) Endpoint {
	endpoint := NewNilEndpoint()
	configure(endpoint)
	return endpoint
}
