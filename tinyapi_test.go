package tinyapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andriiyaremenko/tinyapi/api"
	"github.com/andriiyaremenko/tinyapi/middleware"
)

func TestEndpoint(t *testing.T) {
	endpoint := NewEndpoint("/", func(e api.Endpoint) api.Endpoint {
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
	endpoint := NewEndpoint("/", func(e api.Endpoint) api.Endpoint {
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

func TestCombineEndpoints(t *testing.T) {
	endpoint := NewEndpoint("/", func(e api.Endpoint) api.Endpoint {
		e.Get("/", func(w http.ResponseWriter, req *http.Request, _ map[string]string) {
			w.WriteHeader(http.StatusOK)
		})
		return e
	})

	addTest1Header := middleware.AddHeader("test1", "success")
	addTest2Header := middleware.AddHeader("test2", "success")

	ts := httptest.NewServer(CombineEndpoints("/", CombineMiddleware(addTest1Header, addTest2Header), endpoint))
	defer ts.Close()
	resp, err := http.Get(fmt.Sprintf("%s/15", ts.URL))

	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if resp.Header.Get("test1") != "success" {
		t.Error("addTest1Header was not called, headers were not set")
	}

	if resp.Header.Get("test2") != "success" {
		t.Error("addTest2Header was not called, headers were not set")
	}
}
