// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package mailchimp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
)

// Client handles communication with mailchimp servers
// it fulfills the MailchimpClient interface
type Client struct {
	token      string
	APIURL     string
	HTTPClient *http.Client
	Batch      *BatchQueue
	debug      bool

	log *logrus.Logger
}

// increments on each request made to Do()
var _requestCount int

// NewClient returns a new Mailchimp client with your token
func NewClient(token string) *Client {

	split := strings.Split(token, "-")
	if len(split) != 2 {
		logrus.WithFields(logrus.Fields{
			"token": token,
		}).Debug("malformed token", caller())
		return nil
	}
	apiurl := "https://" + split[1] + ".api.mailchimp.com/3.0/"

	httpclient := &http.Client{}
	return &Client{
		token:      token,
		APIURL:     apiurl,
		HTTPClient: httpclient,
		debug:      false,
		log:        logrus.New(),
	}
}

// Debug will print some request debug information to console.
// toggle with set parameter
func (c *Client) Debug(set ...bool) bool {
	if len(set) > 0 {
		c.debug = set[0]
	}
	if c.debug {
		c.log.Level = logrus.DebugLevel
	} else {
		c.log.Level = logrus.ErrorLevel
	}
	return c.debug
}

// Log returns a logrus interface
func (c *Client) Log() *logrus.Logger {
	return c.log
}

// Clone returns a client with the same preferences.
// http client (and optional batch operation) is ignored
func (c *Client) Clone() *Client {
	return &Client{
		token:      c.token,
		APIURL:     c.APIURL,
		HTTPClient: &http.Client{},
		debug:      c.debug,
		log:        c.log,
	}
}

// Parameters is an alias for Request parameters map string interface
type Parameters map[string]interface{}

// ----------------------------
// internal methods

// Get prepares a GET request to a resource with parameters.
// It returns the body as []byte
func (c *Client) Get(resource string, parameters map[string]interface{}) ([]byte, error) {

	req, err := http.NewRequest("GET", singleJoiningSlash(c.APIURL, resource), nil)
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Debug("malformed request", caller())
		return nil, err
	}

	// add parameters
	c.addParameters(req, parameters)

	return c.Do(req)
}

// Post prepares a POST request to a resource with parameters and JSON body marshalled from
// the data object provided. It returns the body as []byte
func (c *Client) Post(resource string, parameters map[string]interface{}, data interface{}) ([]byte, error) {

	js, err := json.Marshal(data)
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Debug("json error", caller())
		return nil, err
	}

	body := bytes.NewBuffer(js)
	req, err := http.NewRequest("POST", singleJoiningSlash(c.APIURL, resource), body)
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Debug("malformed request", caller())
		return nil, err
	}

	// add parameters
	c.addParameters(req, parameters)

	return c.Do(req)
}

// Patch prepares a PATCH request to a resource with parameters and JSON body marshalled from
// the data object provided. It returns the body as []byte
func (c *Client) Patch(resource string, parameters map[string]interface{}, data interface{}) ([]byte, error) {

	js, err := json.Marshal(data)
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Debug("json error", caller())
		return nil, err
	}

	body := bytes.NewBuffer(js)
	req, err := http.NewRequest("PATCH", singleJoiningSlash(c.APIURL, resource), body)
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Debug("malformed request", caller())
		return nil, err
	}

	// add parameters
	c.addParameters(req, parameters)

	return c.Do(req)
}

// Put prepares a PUT request to a resource with parameters and JSON body marshalled from
// the data object provided. It returns the body as []byte
// Compared to PATCH, PUT will succeed even when the object has previously been deleted.
func (c *Client) Put(resource string, parameters map[string]interface{}, data interface{}) ([]byte, error) {

	js, err := json.Marshal(data)
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Debug("json error", caller())
		return nil, err
	}

	body := bytes.NewBuffer(js)
	req, err := http.NewRequest("PUT", singleJoiningSlash(c.APIURL, resource), body)
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Debug("malformed request", caller())
		return nil, err
	}

	// add parameters
	c.addParameters(req, parameters)

	return c.Do(req)
}

// Delete prepares a DELETE request to a resource
func (c *Client) Delete(resource string) error {

	req, err := http.NewRequest("DELETE", singleJoiningSlash(c.APIURL, resource), nil)
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Debug("malformed request", caller())
		return err
	}
	_, err = c.Do(req)
	return err
}

// Do adds auth token and performs a request to the api. It returns the body as []byte or an error
// that is castable to a mailchimp.Error type for more information about the request.
func (c *Client) Do(request *http.Request) ([]byte, error) {

	if request == nil {
		return nil, fmt.Errorf("can't send nil request")
	}

	_requestCount++
	c.log.WithFields(logrus.Fields{
		"request_count": _requestCount,
		"method":        request.Method,
		"url":           request.URL,
	}).Debug(request.Method, "request")

	// Do we have a batch operation running currently?
	if c.Batch != nil {
		return c.Batch.Do(request)
	}

	if c.token != "" {
		request.SetBasicAuth("OAuthToken", c.token)
	}

	resp, err := c.HTTPClient.Do(request)
	if err != nil {
		c.log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Debug("request error", caller())
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.log.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Debug("response error", caller())
			return nil, err
		}
		return body, nil

	// normal repsonse for all DELETE requests and some POST request.
	case http.StatusNoContent:
		return []byte{}, nil

	default:
		c.log.WithFields(logrus.Fields{
			"code":   resp.StatusCode,
			"method": request.Method,
			"url":    request.URL.String(),
		}).Debug("non success response code", caller())

		err := c.handleError(resp)
		return nil, err
	}
}

// addParameters adds parameters from a map to a request
func (c *Client) addParameters(request *http.Request, params map[string]interface{}) {
	// add parameters
	values := request.URL.Query()
	for key, value := range params {
		switch v := value.(type) {
		case string:
			values.Add(key, v)
		case int:
			values.Add(key, fmt.Sprintf("%d", v))
		}
	}
	request.URL.RawQuery = values.Encode()
}

// handleError translates errors provided by the API to a Error struct
func (c *Client) handleError(response *http.Response) Error {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return Error{
			Title:  "Unknown error",
			Detail: err.Error(),
		}
	}

	var e Error
	err = json.Unmarshal(body, &e)
	if err != nil {
		c.log.Debug(string(body))
		return Error{
			Title:  "Response error",
			Detail: err.Error(),
		}
	}

	return e
}
