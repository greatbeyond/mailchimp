// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Written by David Högborg <d@greatbeyond.se>, 2016
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testing

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/kr/pretty"
	. "gopkg.in/check.v1"
)

// MockServer holds queries of mock responses and stores the requests made
// to it.
type MockServer struct {
	// Checker are invalidated on every new function call. Update before every usage.
	Checker *C

	Responses []*MockResponse
	Requests  []*http.Request

	// The BaseURL of the server is unique with every server start.
	// Match urls to this by concactenating: s.server.BaseURL+"/resource"
	BaseURL    string
	Server     *httptest.Server
	HTTPClient *http.Client
}

// MockResponse defines a response to a matching request. Requests are matched based on
// Method and order of when they were added to the response queue.
type MockResponse struct {
	// Method matches against a incomming request
	Method string
	// use http.Status<something> to reponde to a request.
	Code int
	// the response body to send back.
	Body string
	// CheckFn is called on each request match.
	// Assert that the URL in the request matches what you expect, like so:
	// c.Assert(r.RequestURI, Equals, s.server.BaseURL+"/resource")
	CheckFn func(*http.Request, string)
	// Persistant controls if the response can remain and be used again.
	Persistant bool

	RequestBody string
	// Will hold a reference to the request after response is matched to a request
	Request *http.Request
	// Hits increments each time the response is sent back
	Hits      int
	satisfied bool
}

// NewMockServer returns a new mocking server
func NewMockServer() *MockServer {

	m := &MockServer{
		Responses: []*MockResponse{},
		Requests:  []*http.Request{},
	}

	server := httptest.NewServer(http.HandlerFunc(m.HandleRequest))

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: func(req *http.Request) (*url.URL, error) {
				return url.Parse(server.URL)
			},
		},
	}

	m.BaseURL = server.URL
	m.Server = server
	m.HTTPClient = client

	return m
}

// AddResponse adds a mock response that HandleRequest will look foor
func (m *MockServer) AddResponse(r *MockResponse) *MockResponse {
	m.Responses = append(m.Responses, r)
	return r
}

// VerifyNoMoreRequests checks that no requests are unmet
func (m *MockServer) VerifyNoMoreRequests(c *C) {
	unsatisified := []*MockResponse{}
	for _, r := range m.Responses {
		if !r.satisfied && !r.Persistant {
			pretty.Println("Unsatisfied response:", r)
			unsatisified = append(unsatisified, r)
		}
	}

	if len(unsatisified) > 0 {
		c.Fatal("server has unsatisfied responses")
		c.Fail()
	}
}

// SetChecker sets the checker for the next request handler
func (m *MockServer) SetChecker(c *C) { m.Checker = c }

// Close shuts down the server
func (m *MockServer) Close() {
	m.Server.Close()
}

// HandleRequest is a HTTP handler that matches the request to the mock responses
// If a response with a matching url is found, the response body is written
// to the writer and the CheckFn is called.
func (m *MockServer) HandleRequest(w http.ResponseWriter, r *http.Request) {

	var response *MockResponse
	for _, resp := range m.Responses {
		if !resp.satisfied && resp.Method == r.Method {
			response = resp
			break
		}
	}

	if response == nil {

		if m.Checker != nil {
			errstr := fmt.Sprintf("Mock server: no matching response to request for %s:%s\n", r.Method, r.RequestURI)
			m.Checker.Fatal(errstr)
		}

		w.WriteHeader(http.StatusTeapot)
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintf(w, "no matching response to request for %s:%s\n", r.Method, r.RequestURI)

		return
	}

	body, _ := ioutil.ReadAll(r.Body)

	response.Hits++
	if !response.Persistant {
		response.satisfied = true
	}

	response.Request = r
	response.RequestBody = string(body)

	if response.CheckFn != nil {
		response.CheckFn(r, response.RequestBody)
	}

	m.Requests = append(m.Requests, r)

	w.WriteHeader(response.Code)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, response.Body)
}
