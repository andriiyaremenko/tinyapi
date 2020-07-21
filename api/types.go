package api

import (
	"context"
	"net/http"
)

type Endpoint interface {
	http.Handler

	Path() string
	PrependPath(prefix string)

	NotFound(handler http.HandlerFunc)

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

type Logger interface {
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}

type GetLogger func(ctx context.Context, module string) Logger
