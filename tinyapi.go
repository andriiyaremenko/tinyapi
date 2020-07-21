package tinyapi

import (
	"net/http"

	"github.com/andriiyaremenko/tinyapi/api"
	"github.com/andriiyaremenko/tinyapi/internal"
)

func CombineMiddleware(ms ...api.Middleware) api.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for _, m := range ms {
			next = m(next)
		}
		return next
	}
}

func CombineEndpoints(path string, middleware api.Middleware, endpoints ...api.Endpoint) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		api := http.NewServeMux()

		for _, e := range endpoints {
			e.PrependPath(path)
			api.Handle(e.Path(), e)
		}

		api.ServeHTTP(w, req)
	}

	if middleware == nil {
		return http.HandlerFunc(fn)
	}

	return http.HandlerFunc(middleware(fn))
}

func NewEndpoint(path string, configure api.Configure) api.Endpoint {
	if path[0] != '/' {
		path = "/" + path
	}

	if path[len(path)-1] != '/' {
		path = path + "/"
	}

	return configure(internal.DefaultEndpoint(path))
}
