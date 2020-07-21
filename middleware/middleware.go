package middleware

import (
	"context"
	"net/http"

	"github.com/andriiyaremenko/tinyapi/api"
	"github.com/andriiyaremenko/tinyapi/internal"
)

func AddHeader(key, value string) api.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set(key, value)
			next(w, req)
		}
	}
}

func AddContentType(contentType string) api.Middleware {
	return AddHeader("Content-Type", contentType)
}

func AddJSONContentType() api.Middleware {
	return AddContentType("application/json")
}

func AddTracer(getLogger func(ctx context.Context, module string) api.Logger) api.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			t := internal.NewTracer(w)
			next(t, req)

			l := getLogger(req.Context(), "Tracer")
			code := t.StatusCode

			switch {
			case code >= 400:
				l.Errorf("%s %s %s --> %d %s",
					req.RemoteAddr, req.Method, req.RequestURI,
					code, http.StatusText(code),
				)
			case code >= 300:
				l.Warnf("%s %s %s --> %d %s",
					req.RemoteAddr, req.Method, req.RequestURI,
					code, http.StatusText(code),
				)
			default:
				l.Infof("%s %s %s --> %d %s",
					req.RemoteAddr, req.Method, req.RequestURI,
					code, http.StatusText(code),
				)
			}
		}
	}

}
