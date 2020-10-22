package api

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

func NewDefaultTracer(w http.ResponseWriter, getLogger GetLogger) *DefaultTracer {
	return &DefaultTracer{w, getLogger, http.StatusOK}
}

type DefaultTracer struct {
	http.ResponseWriter

	getLogger  GetLogger
	statusCode int
}

func (t *DefaultTracer) Trace(req *http.Request) {
	code := t.statusCode
	logger := t.getLogger(req.Context(), "DefaultTracer")

	switch {
	case code >= 400:
		logger.Errorf("%s %s %s --> %d %s",
			req.RemoteAddr, req.Method, req.RequestURI,
			code, http.StatusText(code),
		)
	case code >= 300:
		logger.Warnf("%s %s %s --> %d %s",
			req.RemoteAddr, req.Method, req.RequestURI,
			code, http.StatusText(code),
		)
	default:
		logger.Infof("%s %s %s --> %d %s",
			req.RemoteAddr, req.Method, req.RequestURI,
			code, http.StatusText(code),
		)
	}

}

func (t *DefaultTracer) WriteHeader(code int) {
	t.statusCode = code
	t.ResponseWriter.WriteHeader(code)
}

func (w *DefaultTracer) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if hj, ok := w.ResponseWriter.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, fmt.Errorf("Hijacker interface is not supported")
}

// verify Hijacker interface implementation
var _ http.Hijacker = &DefaultTracer{}
