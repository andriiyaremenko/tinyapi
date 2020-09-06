package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/andriiyaremenko/tinyapi/api"
)

type CORSOptions struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	MaxAge           int
}

func AddCORS(opts CORSOptions) api.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			if req.Method == http.MethodOptions && req.Header.Get("Access-Control-Request-Method") != "" {
				handlePreflight(opts, w, req)
				w.WriteHeader(http.StatusNoContent)
				return
			}
			handleReq(opts, w, req)
			next(w, req)
		}
	}
}

func handlePreflight(opts CORSOptions, w http.ResponseWriter, req *http.Request) {
	headers := w.Header()
	origin := req.Header.Get("Origin")

	headers.Add("Vary", "Origin")
	headers.Add("Vary", "Access-Control-Request-Method")
	headers.Add("Vary", "Access-Control-Request-Headers")

	if origin == "" {
		return
	}

	if !isOriginAllowed(origin, opts) {
		return
	}

	method := req.Header.Get("Access-Control-Request-Method")
	if !isMethodAllowed(method, opts) {
		return
	}

	reqHeaders := strings.Split(req.Header.Get("Access-Control-Request-Headers"), ",")
	if !areHeadersAllowed(reqHeaders, opts) {
		return
	}

	headers.Set("Access-Control-Allow-Origin", origin)
	headers.Set("Access-Control-Allow-Methods", strings.ToUpper(method))

	if len(reqHeaders) > 0 {
		headers.Set("Access-Control-Allow-Headers", strings.Join(reqHeaders, ", "))
	}

	if opts.AllowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}

	if opts.MaxAge > 0 {
		headers.Set("Access-Control-Max-Age", strconv.Itoa(opts.MaxAge))
	}
}

func handleReq(opts CORSOptions, w http.ResponseWriter, req *http.Request) {
	headers := w.Header()

	headers.Add("Vary", "Origin")

	origin := req.Header.Get("Origin")

	if origin == "" {
		return
	}

	if !isOriginAllowed(origin, opts) {
		return
	}

	if !isMethodAllowed(req.Method, opts) {
		return
	}

	headers.Set("Access-Control-Allow-Origin", origin)

	if opts.AllowCredentials {
		headers.Set("Access-Control-Allow-Credentials", "true")
	}
}

func isOriginAllowed(origin string, opts CORSOptions) bool {
	for _, o := range opts.AllowOrigins {
		if o == "*" || strings.ToLower(o) == strings.ToLower(origin) {
			return true
		}
	}

	return false
}

func isMethodAllowed(method string, opts CORSOptions) bool {
	if len(opts.AllowMethods) == 0 {
		return false
	}

	method = strings.ToUpper(method)
	if method == http.MethodOptions {
		return true
	}

	for _, m := range opts.AllowMethods {
		if strings.ToUpper(m) == method {
			return true
		}
	}

	return false
}

func areHeadersAllowed(headers []string, opts CORSOptions) bool {
	if len(opts.AllowHeaders) == 0 {
		return true
	}

loop:
	for _, h := range headers {
		header := http.CanonicalHeaderKey(h)

		for _, allowH := range opts.AllowHeaders {
			if header == allowH {
				continue loop
			}
		}

		return false
	}

	return true
}
