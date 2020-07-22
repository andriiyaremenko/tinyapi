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

func CombineEndpoints(path string, notFound http.HandlerFunc, middleware api.Middleware, endpoints ...api.Endpoint) http.Handler {
	api := http.NewServeMux()
	if notFound == nil {
		notFound = http.NotFound
	}

	for _, e := range endpoints {
		e.PrependPath(path)
		e.NotFound(notFound)
		api.Handle(e.Path(), e)
	}

	fn := func(w http.ResponseWriter, req *http.Request) {
		// we would just call api.ServeHTTP(w, req), but we want our own NotFound handler
		// and there is no evident way to replace default NotFound handler with http.ServeMux
		// so we replicate http.ServeMux.ServeHTTP method here and hijack NotFound and handle it by ourselfs
		if req.RequestURI == "*" {
			if req.ProtoAtLeast(1, 1) {
				w.Header().Set("Connection", "close")
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		h, pattern := api.Handler(req)

		if pattern == "" {
			notFound(w, req)
			return
		}

		h.ServeHTTP(w, req)
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
