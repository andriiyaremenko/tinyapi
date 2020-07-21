package middleware

import (
	"net/http"

	"github.com/andriiyaremenko/tinyapi/api"
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
