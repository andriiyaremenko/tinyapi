package tinyapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/andriiyaremenko/tinyapi/utils"
)

func TestEndpoint(t *testing.T) {
	endpoint := NewEndpoint(func(e Endpoint) {
		e.Middleware(func(w http.ResponseWriter, req *http.Request) {
			t.Logf("%v: %v", req.Method, req.URL)
		})
		e.Get(`\d+`, func(w http.ResponseWriter, req *http.Request) {
			url := req.URL.RequestURI()
			param := path.Base(url)
			fmt.Fprintf(w, "got %s", param)
		})
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
		t.Errorf(`Endpoint.Get(\d+) = %v`, resp)
	}
}

func Test404(t *testing.T) {
	endpoint := NewEndpoint(func(e Endpoint) {
		e.Get("", func(w http.ResponseWriter, req *http.Request) {

			url := req.URL.RequestURI()
			param := path.Base(url)
			fmt.Fprintf(w, "got %s", param)
		})
		e.Get("/", func(w http.ResponseWriter, req *http.Request) {
			url := req.URL.RequestURI()
			param := path.Base(url)
			fmt.Fprintf(w, "got %s", param)
		})
	})

	ts := httptest.NewServer(endpoint)
	defer ts.Close()
	resp, err := http.Get(fmt.Sprintf("%s/15", ts.URL))
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	if resp.StatusCode != 404 {
		t.Errorf(`Endpoint.Get(\d+) = %v`, resp)
	}
}

func TestGetParam(t *testing.T) {
	req := httptest.NewRequest("GET", "http://example.com/foo/1", nil)
	v, ok := utils.RequestParam(req)
	if !ok {
		t.Errorf(`RequestParam(%v) = "%s", false`, req, v)
	}
	if v != "1" {
		t.Errorf(`RequestParam(%v) = "%s", true`, req, v)
	}
}
