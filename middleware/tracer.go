package middleware

import (
	"net/http"

	"github.com/andriiyaremenko/tinyapi/api"
	"github.com/andriiyaremenko/tinyapi/internal"
)

func AddTracer(getLogger api.GetLogger) api.Middleware {
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
