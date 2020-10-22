package tinyapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andriiyaremenko/tinyapi/api"
	"github.com/andriiyaremenko/tinyapi/middleware"
	"github.com/stretchr/testify/assert"
)

func TestEndpoint(t *testing.T) {
	assert := assert.New(t)
	endpoint := map[string]api.Endpoint{
		"/": {
			api.GET: api.RouteSegment{
				"/": func(w http.ResponseWriter, req *http.Request) {
					fmt.Fprintf(w, "nothing")
				},
				"/:id": func(w http.ResponseWriter, req *http.Request) {
					id, _ := GetRouteParams(req, "id")
					fmt.Fprintf(w, id)
				},
				"/:id/:nothing": func(w http.ResponseWriter, req *http.Request) {
					id, _ := GetRouteParams(req, "id")
					fmt.Fprintf(w, id)
				},
			},
		},
	}

	Print(endpoint)
	t.Log(string(SprintJSON(endpoint)))

	ts := httptest.NewServer(CombineEndpoints(endpoint, nil, nil))

	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/15?test", ts.URL))
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	defer resp.Body.Close()

	assert.Equal(string(r), "15")

	resp, err = http.Get(fmt.Sprintf("%s/", ts.URL))
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	r, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	defer resp.Body.Close()

	assert.Equal(string(r), "nothing")
}

func Test404(t *testing.T) {
	assert := assert.New(t)
	endpoint := map[string]api.Endpoint{
		"/": {
			api.GET: api.RouteSegment{
				"/": func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusOK)
				},
			},
		},
	}

	Print(endpoint)
	t.Log(string(SprintJSON(endpoint)))

	ts := httptest.NewServer(CombineEndpoints(endpoint, nil, nil))

	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/15", ts.URL))
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	assert.Equal(404, resp.StatusCode)
}

func TestCombineEndpoints(t *testing.T) {
	assert := assert.New(t)
	endpoint := map[string]api.Endpoint{
		"/": {
			api.GET: api.RouteSegment{
				"/": func(w http.ResponseWriter, req *http.Request) {
					w.WriteHeader(http.StatusOK)
				},
			},
		},
	}

	Print(endpoint)
	t.Log(string(SprintJSON(endpoint)))

	addTest1Header := middleware.AddHeader("test1", "success")
	addTest2Header := middleware.AddHeader("test2", "success")
	ts := httptest.NewServer(CombineEndpoints(endpoint, CombineMiddleware(addTest1Header, addTest2Header), nil))

	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/15", ts.URL))
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	assert.Equal(404, resp.StatusCode)
	assert.Equal("success", resp.Header.Get("test1"), "addTest1Header was not called, headers were not set")
	assert.Equal("success", resp.Header.Get("test2"), "addTest2Header was not called, headers were not set")
}

func TestSameSectionEndpoint(t *testing.T) {
	assert := assert.New(t)
	endpoint := map[string]api.Endpoint{
		"/": {
			api.GET: {
				"/abc/ad": func(w http.ResponseWriter, req *http.Request) {
					fmt.Fprintf(w, "right")
				},
				"/be/de": func(w http.ResponseWriter, req *http.Request) {
					fmt.Fprintf(w, "wrong")
				},
			},
			api.CONNECT: {"/abc/ad": func(w http.ResponseWriter, req *http.Request) {}},
			api.DELETE:  {"/abc/ad": func(w http.ResponseWriter, req *http.Request) {}},
			api.HEAD:    {"/abc/ad": func(w http.ResponseWriter, req *http.Request) {}},
			api.OPTIONS: {"/abc/ad": func(w http.ResponseWriter, req *http.Request) {}},
			api.PATCH:   {"/abc/ad": func(w http.ResponseWriter, req *http.Request) {}},
			api.POST:    {"/abc/ad": func(w http.ResponseWriter, req *http.Request) {}},
			api.PUT:     {"/abc/ad": func(w http.ResponseWriter, req *http.Request) {}},
			api.TRACE:   {"/abc/ad": func(w http.ResponseWriter, req *http.Request) {}},
		},
	}

	Print(endpoint)
	t.Log(string(SprintJSON(endpoint)))

	ts := httptest.NewServer(CombineEndpoints(endpoint, nil, nil))

	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf("%s/abc/ad", ts.URL))
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	defer resp.Body.Close()

	assert.Equal(string(r), "right")

	resp, err = http.Get(fmt.Sprintf("%s/be/de", ts.URL))
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	r, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	defer resp.Body.Close()

	assert.Equal(string(r), "wrong")
}
