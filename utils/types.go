package utils

import (
	_ "encoding/json"
)

type RouteDefinition struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

type ApiDefinition struct {
	Routes []RouteDefinition `json:"routes"`
}
