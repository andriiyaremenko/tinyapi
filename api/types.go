package api

import (
	"net/http"
)

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

type HandlerFunc func(w http.ResponseWriter, req *http.Request, params map[string]string)

type Configure func(e Endpoint) Endpoint

type Middleware func(next http.HandlerFunc) http.HandlerFunc
