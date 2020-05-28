package middleware

import (
	"net/http"
)

func ContentType(contentType string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", contentType)
	}
}

func JSONContentType() http.HandlerFunc {
	return ContentType("application/json")
}
