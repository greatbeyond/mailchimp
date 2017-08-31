// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package mailchimp

import (
	"context"
	"net/http"
	"testing"

	t "github.com/greatbeyond/mailchimp/testing"
	"github.com/sirupsen/logrus"
	check "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test_Mailchimp(t *testing.T) { check.TestingT(t) }

var _ = check.Suite(&MailchimpTestSuite{})

type MailchimpTestSuite struct {
	Client *Client
	server *t.MockServer
	ctx    context.Context
}

func (s *MailchimpTestSuite) SetUpSuite(c *check.C) {
	Log.Level = logrus.DebugLevel
}

func (s *MailchimpTestSuite) SetUpTest(c *check.C) {
	s.server = t.NewMockServer()
	s.server.SetChecker(c)

	s.Client = NewClient()
	s.Client.HTTPClient = s.server.HTTPClient

	s.ctx = NewContextWithToken(context.Background(), "b12824bd84759ef84abc67fd789e7570-us13")
	// We need http to use the mock server
	s.ctx = NewContextWithURL(s.ctx, "http://us13.api.mailchimp.com/3.0/")

}

func (s *MailchimpTestSuite) TearDownTest(c *check.C) {}

// -------------------------------------------------------------------
// GET Requests

func (s *MailchimpTestSuite) Test_Get_Normal(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   "{}",
		CheckFn: func(r *http.Request, body string) {
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/test?param=value")
		},
	})
	resp, err := s.Client.Get(s.ctx, "test", map[string]interface{}{
		"param": "value",
	})
	c.Assert(err, check.IsNil)
	c.Assert(string(resp), check.Equals, "{}\n")
}

func (s *MailchimpTestSuite) Test_Get_Malformed(c *check.C) {

	resp, err := s.Client.Get(s.ctx, "%hse%fa%%", map[string]interface{}{
		"param": "value",
	})
	c.Assert(err, check.Not(check.IsNil))
	c.Assert(string(resp), check.Equals, "")
}

// -------------------------------------------------------------------
// POST Requests

func (s *MailchimpTestSuite) Test_Post_Normal(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "POST",
		Code:   200,
		Body:   "{}",
		CheckFn: func(r *http.Request, body string) {
			c.Assert(body, check.Equals, `"payload"`)
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/test?param=value")
		},
	})
	resp, err := s.Client.Post(s.ctx, "test", map[string]interface{}{
		"param": "value",
	}, "payload")
	c.Assert(err, check.IsNil)
	c.Assert(string(resp), check.Equals, "{}\n")
}

func (s *MailchimpTestSuite) Test_Post_MalformedData(c *check.C) {
	baddata := map[int]string{
		3: "three",
	}
	resp, err := s.Client.Post(s.ctx, "test", nil, baddata)
	c.Assert(resp, check.IsNil)
	c.Assert(err, check.Not(check.IsNil))
}

func (s *MailchimpTestSuite) Test_Post_MalformedURL(c *check.C) {
	resp, err := s.Client.Post(s.ctx, "%hse%fa%%", map[string]interface{}{
		"param": "value",
	}, "payload")
	c.Assert(err, check.Not(check.IsNil))
	c.Assert(string(resp), check.Equals, "")
}

// -------------------------------------------------------------------
// PATCH Requests

func (s *MailchimpTestSuite) Test_Patch_Normal(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "PATCH",
		Code:   200,
		Body:   "{}",
		CheckFn: func(r *http.Request, body string) {
			c.Assert(body, check.Equals, `"payload"`)
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/test?param=value")
		},
	})
	resp, err := s.Client.Patch(s.ctx, "test", map[string]interface{}{
		"param": "value",
	}, "payload")
	c.Assert(err, check.IsNil)
	c.Assert(string(resp), check.Equals, "{}\n")
}

func (s *MailchimpTestSuite) Test_Patch_MalformedData(c *check.C) {
	baddata := map[int]string{
		3: "three",
	}
	resp, err := s.Client.Patch(s.ctx, "test", nil, baddata)
	c.Assert(resp, check.IsNil)
	c.Assert(err, check.Not(check.IsNil))
}

func (s *MailchimpTestSuite) Test_Patch_MalformedURL(c *check.C) {
	resp, err := s.Client.Patch(s.ctx, "%hse%fa%%", map[string]interface{}{
		"param": "value",
	}, "payload")
	c.Assert(err, check.Not(check.IsNil))
	c.Assert(string(resp), check.Equals, "")
}

// -------------------------------------------------------------------
// PUT Requests

func (s *MailchimpTestSuite) Test_Put_Normal(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "PUT",
		Code:   200,
		Body:   "{}",
		CheckFn: func(r *http.Request, body string) {
			c.Assert(body, check.Equals, `"payload"`)
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/test?param=value")
		},
	})
	resp, err := s.Client.Put(s.ctx, "test", map[string]interface{}{
		"param": "value",
	}, "payload")
	c.Assert(err, check.IsNil)
	c.Assert(string(resp), check.Equals, "{}\n")
}

func (s *MailchimpTestSuite) Test_Put_MalformedData(c *check.C) {
	baddata := map[int]string{
		3: "three",
	}
	resp, err := s.Client.Put(s.ctx, "test", nil, baddata)
	c.Assert(resp, check.IsNil)
	c.Assert(err, check.Not(check.IsNil))
}

func (s *MailchimpTestSuite) Test_Put_MalformedURL(c *check.C) {
	resp, err := s.Client.Put(s.ctx, "%hse%fa%%", map[string]interface{}{
		"param": "value",
	}, "payload")
	c.Assert(err, check.Not(check.IsNil))
	c.Assert(string(resp), check.Equals, "")
}

// -------------------------------------------------------------------
// DELETE Requests

func (s *MailchimpTestSuite) Test_Delete_Normal(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "DELETE",
		Code:   http.StatusNoContent,
		Body:   "",
		CheckFn: func(r *http.Request, body string) {
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/test")
		},
	})
	err := s.Client.Delete(s.ctx, "test")
	c.Assert(err, check.IsNil)

}

func (s *MailchimpTestSuite) Test_Delete_Malformed(c *check.C) {
	err := s.Client.Delete(s.ctx, "%hse%fa%%")
	c.Assert(err, check.Not(check.IsNil))
}

// -------------------------------------------------------------------
// Do request

func (s *MailchimpTestSuite) Test_Do_Normal(c *check.C) {
	req, _ := http.NewRequest("GET", "http://us13.api.mailchimp.com/3.0/test", nil)
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   "{}",
		CheckFn: func(r *http.Request, body string) {
			// base64 encoded username:password (OAuthToken:[token])
			c.Assert(r.Header.Get("Authorization"), check.Equals, "Basic T0F1dGhUb2tlbjpiMTI4MjRiZDg0NzU5ZWY4NGFiYzY3ZmQ3ODllNzU3MC11czEz")
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/test")
		},
	})
	resp, err := s.Client.Do(req.WithContext(s.ctx))
	c.Assert(err, check.IsNil)
	c.Assert(string(resp), check.Equals, "{}\n")
}

func (s *MailchimpTestSuite) Test_Do_NonSuccessResponse(c *check.C) {

	req, _ := http.NewRequest("GET", "http://us13.api.mailchimp.com/3.0/test", nil)
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   501,
		Body:   `{"type":"internal error","title":"Internal Error","status":501}`,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/test")
		},
	})
	resp, err := s.Client.Do(req.WithContext(s.ctx))
	c.Assert(err, check.DeepEquals, Error{
		Status: 501,
		Type:   "internal error",
		Title:  "Internal Error",
	})
	c.Assert(resp, check.IsNil)
}

func (s *MailchimpTestSuite) Test_Do_BadRequest(c *check.C) {
	req, _ := http.NewRequest("GET", "http://example.net", nil)
	req.URL = nil
	_, err := s.Client.Do(req.WithContext(s.ctx))
	c.Assert(err, check.ErrorMatches, "http: nil Request.URL")
}

func (s *MailchimpTestSuite) Test_Do_NilRequest(c *check.C) {
	_, err := s.Client.Do(nil)
	c.Assert(err, check.ErrorMatches, "can't send nil request")
}

// -------------------------------------------------------------------
// parameters

func (s *MailchimpTestSuite) Test_addParameters_String(c *check.C) {
	req, _ := http.NewRequest("GET", "http://example.net", nil)
	s.Client.addParameters(req, map[string]interface{}{
		"test": "value",
	})
	c.Assert(req.URL.RequestURI(), check.Equals, "/?test=value")
}

func (s *MailchimpTestSuite) Test_addParameters_Int(c *check.C) {
	req, _ := http.NewRequest("GET", "http://example.net", nil)
	s.Client.addParameters(req, map[string]interface{}{
		"test": 2,
	})
	c.Assert(req.URL.RequestURI(), check.Equals, "/?test=2")
}
