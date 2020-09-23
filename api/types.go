package api

import (
	"context"
	"net/http"
)

type verb string

const (
	CONNECT   = http.MethodConnect
	DELETE    = http.MethodDelete
	GET       = http.MethodGet
	HEAD      = http.MethodHead
	OPTIONS   = http.MethodOptions
	PATCH     = http.MethodPatch
	POST      = http.MethodPost
	PUT       = http.MethodPut
	TRACE     = http.MethodTrace
	WEBSOCKET = "WEBSOCKET"
)

type RouteSegment = map[string]func(map[string]string) http.HandlerFunc
type Endpoint = map[verb]RouteSegment
type Middleware func(next http.HandlerFunc) http.HandlerFunc

type Logger interface {
	Infof(format string, v ...interface{})
	Warnf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
}

type GetLogger func(ctx context.Context, module string) Logger
