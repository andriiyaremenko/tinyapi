package middleware

import (
	"github.com/andriiyaremenko/tinyapi/api"
)

func AddContentType(contentType string) api.Middleware {
	return AddHeader("Content-Type", contentType)
}

func AddJSONContentType() api.Middleware {
	return AddContentType("application/json")
}
