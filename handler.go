package tinyapi

import (
	"net/http"
	"sort"

	"github.com/andriiyaremenko/tinyapi/api"
	"github.com/andriiyaremenko/tinyapi/internal"
)

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
			if string(method) != req.Method {
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
