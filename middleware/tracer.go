package middleware

import (
	"net/http"

	"github.com/andriiyaremenko/tinyapi/api"
)

func AddTracer(getTracer api.GetTracer) api.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			t := getTracer(w)

			next(t, req)
			t.Trace(req)
		}
	}
}

func AddDefaultTracer(getLogger api.GetLogger) api.Middleware {
	getTracer := func(w http.ResponseWriter) api.Tracer { return api.NewDefaultTracer(w, getLogger) }

	return AddTracer(getTracer)
}
