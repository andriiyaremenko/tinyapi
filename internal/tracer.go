package internal

import "net/http"

type tracer struct {
	http.ResponseWriter
	StatusCode int
}

func NewTracer(w http.ResponseWriter) *tracer {
	return &tracer{w, http.StatusOK}
}

func (t *tracer) WriteHeader(code int) {
	t.StatusCode = code
	t.ResponseWriter.WriteHeader(code)
}
