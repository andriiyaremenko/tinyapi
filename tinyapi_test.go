package tinyapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEndpoint(t *testing.T) {
	endpoint := NewEndpoint("/", func(e Endpoint) Endpoint {
		e.Middleware(func(w http.ResponseWriter, req *http.Request) {
			t.Logf("%v: %v", req.Method, req.URL.RequestURI())
		})
		e.Get(":id", func(w http.ResponseWriter, req *http.Request, param map[string]string) {
			fmt.Fprintf(w, "got %s", param["id"])
		})
		return e
	})

	ts := httptest.NewServer(endpoint)
	defer ts.Close()
	resp, err := http.Get(fmt.Sprintf("%s/15", ts.URL))
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	defer resp.Body.Close()
	if "got 15" != string(r) {
		t.Errorf(`Endpoint.Get(:id) = %v`, resp)
	}
}

func Test404(t *testing.T) {
	endpoint := NewEndpoint("/", func(e Endpoint) Endpoint {
		e.Get("/", func(w http.ResponseWriter, req *http.Request, _ map[string]string) {
			w.WriteHeader(http.StatusOK)
		})
		return e
	})

	ts := httptest.NewServer(endpoint)
	defer ts.Close()
	resp, err := http.Get(fmt.Sprintf("%s/15", ts.URL))
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if resp.StatusCode != 404 {
		t.Errorf(`Endpoint.Get(:id) = %v`, resp)
	}
}
