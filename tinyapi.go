package tinyapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/andriiyaremenko/tinyapi/api"
	"github.com/andriiyaremenko/tinyapi/internal"
	"github.com/andriiyaremenko/tinyapi/utils"
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
	api := &apiHandler{endpoints, notFound}
	if middleware == nil {
		return api
	}

	return http.HandlerFunc(middleware(api.ServeHTTP))
}

func Sprint(endpoints map[string]api.Endpoint) string {
	var sb strings.Builder
	var pathSegments []string
	methods := make(map[string][]string)

	for base, e := range endpoints {
		for method, routeSegments := range e {
			for pathSegment := range routeSegments {
				path := internal.CombinePath(base, pathSegment)
				if _, ok := methods[path]; !ok {
					pathSegments = append(pathSegments, path)
				}
				methods[path] = append(methods[path], string(method))
			}
		}
	}

	sort.Strings(pathSegments)
	sb.WriteString("\nAPI definition:\n")

	for _, pathSegment := range pathSegments {
		methods := methods[pathSegment]
		sort.Strings(methods)
		sb.WriteByte('\n')
		for _, method := range methods {
			sb.WriteByte('\t')
			sb.WriteString(method)
			sb.WriteByte('\t')
			sb.WriteString(pathSegment)
			sb.WriteByte('\n')
		}
	}
	sb.WriteByte('\n')

	return sb.String()
}

func SprintJSON(endpoints map[string]api.Endpoint) []byte {
	var pathSegments []string
	methods := make(map[string][]string)
	definition := new(utils.ApiDefinition)

	for base, e := range endpoints {
		for method, routeSegments := range e {
			for pathSegment := range routeSegments {
				path := internal.CombinePath(base, pathSegment)
				if _, ok := methods[path]; !ok {
					pathSegments = append(pathSegments, path)
				}
				methods[path] = append(methods[path], string(method))
			}
		}
	}

	sort.Strings(pathSegments)

	for _, pathSegment := range pathSegments {
		methods := methods[pathSegment]
		sort.Strings(methods)
		for _, method := range methods {
			definition.Routes = append(definition.Routes,
				utils.RouteDefinition{Method: method, Path: pathSegment})
		}
	}

	b, err := json.Marshal(definition)
	if err != nil {
		panic(err)
	}

	return b
}

func Print(endpoints map[string]api.Endpoint) {
	var sb strings.Builder

	sb.WriteString(ANSIColorYellow)
	sb.WriteString(Sprint(endpoints))
	sb.WriteString(ANSIReset)
	fmt.Print(sb.String())
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
