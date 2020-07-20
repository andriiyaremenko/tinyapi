package utils

import (
	"fmt"
	"net/http"
)

func NotFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprint(w, "not found")
}

func InternalServerError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, err)
}
