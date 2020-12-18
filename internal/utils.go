package internal

import (
	"fmt"
	"strings"
)

const (
	ANSIReset       string = "\033[0m"
	ANSIColorRed    string = "\033[31m"
	ANSIColorGreen  string = "\033[32m"
	ANSIColorYellow string = "\033[33m"
)

// Returns text prepended by ANSI color code and appended by ANSI color reset code
func PaintText(color string, text string) string {
	return color + text + ANSIReset
}

func CombinePath(segments ...string) string {
	var sb strings.Builder

	sb.WriteRune('/')

	for _, s := range segments {
		if s == "/" {
			continue
		}

		if len(s) > 1 && s[0] == '/' {
			s = s[1:]
		}

		if len(s) > 1 && s[len(s)-1] == '/' {
			s = s[:len(s)-2]
		}

		sb.WriteString(s)
		sb.WriteRune('/')
	}

	path := sb.String()
	if path == "" {
		return "/"
	}

	return path
}

func GetRouteProps(apiPath, actualRoute string) (map[string]string, bool) {
	path := strings.Trim(apiPath, "/")
	route := strings.Trim(actualRoute, "/")

	pathSegments := strings.Split(path, "/")
	routeSegments := strings.Split(route, "/")

	props := make(map[string]string)

	if path == "" && route == "" {
		return props, true
	}

	if path == "" || route == "" {
		return nil, false
	}

	if len(pathSegments) != len(routeSegments) {
		return nil, false
	}

	for i, ps := range pathSegments {
		if ps[0] != ':' && ps != routeSegments[i] {
			return nil, false
		}

		if ps[0] != ':' {
			continue
		}

		key := ps[1:]
		props = addProp(props, key, routeSegments[i], 0)
	}

	return props, true
}

func addProp(props map[string]string, key, prop string, n int) map[string]string {
	if n != 0 {
		key = fmt.Sprintf("%s:%d", key, n)
	}

	if _, ok := props[key]; ok {
		return addProp(props, key, prop, n+1)
	}

	props[key] = prop
	return props
}
