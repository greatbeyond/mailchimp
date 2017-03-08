// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package mailchimp

import (
	"context"

	t "github.com/greatbeyond/mailchimp/testing"

	check "gopkg.in/check.v1"
)

var _ = check.Suite(&WebhookTestSuite{})

type WebhookTestSuite struct {
	client *Client
	server *t.MockServer
	ctx    context.Context
}

func (s *WebhookTestSuite) SetUpSuite(c *check.C) {}

func (s *WebhookTestSuite) SetUpTest(c *check.C) {
	s.server = t.NewMockServer()
	s.server.SetChecker(c)

	s.client = NewClient()
	s.client.HTTPClient = s.server.HTTPClient

	s.ctx = NewContextWithToken(context.Background(), "b12824bd84759ef84abc67fd789e7570-us13")
	// We need http to use the mock server
	s.ctx = NewContextWithURL(s.ctx, "http://us13.api.mailchimp.com/3.0/")
}

func (s *WebhookTestSuite) TearDownTest(c *check.C) {}

func (s *WebhookTestSuite) Skip_CreateWebhook(c *check.C) {
	createWebhookResponse, err := s.client.CreateWebhook(s.ctx, &CreateWebhook{
		ListID: "1",
		URL:    "http://test.url/webhook",
		Events: &WebhookEvents{
			Subscribe:   true,
			Unsubscribe: true,
		},
		Sources: &WebhookSources{
			User: true,
		},
	})
	c.Assert(err, check.IsNil)
	c.Assert(createWebhookResponse, check.NotNil)
}

func (s *WebhookTestSuite) Skip_GetWebhook(c *check.C) {
	getWebhookResponse, err := s.client.GetWebhook(s.ctx, "1", "2")
	c.Assert(err, check.IsNil)
	c.Assert(getWebhookResponse, check.NotNil)
}

func (s *WebhookTestSuite) Skip_GetWebhooks(c *check.C) {
	getWebhooksResponse, err := s.client.GetWebhooks(s.ctx, "1")
	c.Assert(err, check.IsNil)
	c.Assert(getWebhooksResponse, check.NotNil)
}

func (s *WebhookTestSuite) Skip_DeleteWebhook(c *check.C) {
	getWebhookResponse, err := s.client.GetWebhook(s.ctx, "1", "2")
	c.Assert(err, check.IsNil)
	c.Assert(getWebhookResponse, check.NotNil)

	err = getWebhookResponse.DeleteWebhook(s.ctx)
	c.Assert(err, check.IsNil)
}
