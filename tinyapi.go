package tinyapi

import (
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/andriiyaremenko/tinyapi/api"
	"github.com/andriiyaremenko/tinyapi/internal"
)

const (
	ANSIReset       = "\033[0m"
	ANSIColorYellow = "\033[33m"
)

func CombineMiddleware(ms ...api.Middleware) api.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		for _, m := range ms {
			next = m(next)
		}
		return next
	}
}

func CombineEndpoints(endpoints map[string]api.Endpoint, middleware api.Middleware, notFound http.HandlerFunc) http.Handler {
	var sb strings.Builder
	var pathSegments []string

	sb.WriteString(ANSIColorYellow)
	sb.WriteString("API definition:\n")

	for base, e := range endpoints {
		for method, routeSegments := range e {
			for pathSegment := range routeSegments {
				path := internal.CombinePath(base, pathSegment)
				pathSegments = append(pathSegments, fmt.Sprintf("\t%s \t%s\n", method, path))
			}
		}
	}

	sort.Strings(pathSegments)

	for _, pathSegment := range pathSegments {
		sb.WriteString(pathSegment)
	}

	sb.WriteString(ANSIReset)
	fmt.Println(sb.String())

	api := &apiHandler{endpoints, notFound}
	if middleware == nil {
		return api
	}

	return http.HandlerFunc(middleware(api.ServeHTTP))
}

type apiHandler struct {
	Endpoints map[string]api.Endpoint
	NotFound  http.HandlerFunc
}

func (apiH *apiHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// copied from  http.ServeMux.ServeHTTP method
	if req.RequestURI == "*" {
		if req.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if apiH.NotFound == nil {
		apiH.NotFound = http.NotFound
	}

	var baseRoutes []string
	for base := range apiH.Endpoints {
		baseRoutes = append(baseRoutes, base)
	}

	sort.Strings(baseRoutes)

	for _, base := range baseRoutes {
		e := apiH.Endpoints[base]
		for method, routeSegments := range e {
			if method != api.WEBSOCKET && string(method) != req.Method {
				continue
			}

			var pathSegments []string
			for pathSegment := range routeSegments {
				pathSegments = append(pathSegments, pathSegment)
			}

			sort.Strings(pathSegments)

			for _, pathSegment := range pathSegments {
				handler := routeSegments[pathSegment]
				path := internal.CombinePath(base, pathSegment)

				if props, ok := internal.GetRouteProps(path, req.URL.Path); ok {
					q := req.URL.Query()
					for key, value := range props {
						q.Set(getParamKey(key), value)
					}
					req.URL.RawQuery = q.Encode()
					handler(w, req)
					return
				}
			}
		}
	}

	apiH.NotFound(w, req)
}

func GetRouteParams(req *http.Request, param string) (string, bool) {
	q := req.URL.Query()
	vals, ok := q[getParamKey(param)]

	if !ok || len(vals) == 0 {
		return "", false
	}

	return vals[0], true
}

func getParamKey(param string) string {
	return fmt.Sprintf("params:%s", param)
}
