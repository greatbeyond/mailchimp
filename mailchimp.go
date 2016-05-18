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

	log "github.com/Sirupsen/logrus"
)

// Client handles communication with mailchimp servers
type Client struct {
	token      string
	ApiURL     string
	HttpClient *http.Client
}

// NewClient returns a new Mailchimp client with your token
func NewClient(token string) *Client {

	split := strings.Split(token, "-")
	if len(split) != 2 {
		log.WithFields(log.Fields{
			"token": token,
		}).Warn("malformed token", caller())
		return nil
	}
	apiurl := "https://" + split[1] + ".api.mailchimp.com/3.0/"

	httpclient := &http.Client{}
	return &Client{
		token:      token,
		ApiURL:     apiurl,
		HttpClient: httpclient,
	}
}

// Parameters is an alias for Request parameters map strin interface
type Parameters map[string]interface{}

// ----------------------------
// internal methods

// get prepares a GET request to a resource with parameters.
// It returns the body as []byte
func (c *Client) get(resource string, parameters map[string]interface{}) ([]byte, error) {

	req, err := http.NewRequest("GET", singleJoiningSlash(c.ApiURL, resource), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("malformed request", caller())
		return nil, err
	}

	// add parameters
	c.addParameters(req, parameters)

	return c.do(req)
}

// post prepares a POST request to a resource with parameters and JSON body marshalled from
// the data object provided. It returns the body as []byte
func (c *Client) post(resource string, parameters map[string]interface{}, data interface{}) ([]byte, error) {

	js, err := json.Marshal(data)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("json error", caller())
		return nil, err
	}

	body := bytes.NewBuffer(js)
	req, err := http.NewRequest("POST", singleJoiningSlash(c.ApiURL, resource), body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("malformed request", caller())
		return nil, err
	}

	// add parameters
	c.addParameters(req, parameters)

	return c.do(req)
}

// patch prepares a PATCH request to a resource with parameters and JSON body marshalled from
// the data object provided. It returns the body as []byte
func (c *Client) patch(resource string, parameters map[string]interface{}, data interface{}) ([]byte, error) {

	js, err := json.Marshal(data)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("json error", caller())
		return nil, err
	}

	body := bytes.NewBuffer(js)
	req, err := http.NewRequest("PATCH", singleJoiningSlash(c.ApiURL, resource), body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("malformed request", caller())
		return nil, err
	}

	// add parameters
	c.addParameters(req, parameters)

	return c.do(req)
}

// put prepares a PUT request to a resource with parameters and JSON body marshalled from
// the data object provided. It returns the body as []byte
// Compared to PATCH, PUT will succeed even when the object has previously been deleted.
func (c *Client) put(resource string, parameters map[string]interface{}, data interface{}) ([]byte, error) {

	js, err := json.Marshal(data)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("json error", caller())
		return nil, err
	}

	body := bytes.NewBuffer(js)
	req, err := http.NewRequest("PUT", singleJoiningSlash(c.ApiURL, resource), body)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("malformed request", caller())
		return nil, err
	}

	// add parameters
	c.addParameters(req, parameters)

	return c.do(req)
}

// delete prepares a DELETE request to a resource
func (c *Client) delete(resource string) error {

	req, err := http.NewRequest("DELETE", singleJoiningSlash(c.ApiURL, resource), nil)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("malformed request", caller())
		return err
	}
	_, err = c.do(req)
	return err
}

// do adds auth token and performs a request to the api. It returns the body as []byte or an error
// that is castable to a mailchimp.Error type for more information about the request.
func (c *Client) do(request *http.Request) ([]byte, error) {

	if request == nil {
		return nil, fmt.Errorf("can't send nil request")
	}

	if c.token != "" {
		request.SetBasicAuth("OAuthToken", c.token)
	}

	resp, err := c.HttpClient.Do(request)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("request error", caller())
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Warn("response error", caller())
			return nil, err
		}
		return body, nil

		// normal repsonse for DELETE requests
	case http.StatusNoContent:
		return []byte{}, nil

	default:
		log.WithFields(log.Fields{
			"code":   resp.StatusCode,
			"method": request.Method,
			"url":    request.URL.String(),
		}).Warn("non success response code", caller())

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
		log.Warn(string(body))
		return Error{
			Title:  "Response error",
			Detail: err.Error(),
		}
	}

	return e
}
