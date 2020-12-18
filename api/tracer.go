package api

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/andriiyaremenko/tinyapi/internal"
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
	logger := t.getLogger(req.Context())
	color := internal.ANSIColorGreen

	switch {
	case code >= 400:
		color = internal.ANSIColorRed
	case code >= 300:
		color = internal.ANSIColorYellow
	}

	logger.Printf("%s %s %s --> %d %s",
		req.RemoteAddr, internal.PaintText(color, req.Method), internal.PaintText(color, req.RequestURI),
		internal.PaintText(color, strconv.Itoa(code)), internal.PaintText(color, http.StatusText(code)),
	)
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
