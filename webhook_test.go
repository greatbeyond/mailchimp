// © Copyright 2016 GREAT BEYOND AB
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

package mailchimp

import (
	"context"
	"os"

	check "gopkg.in/check.v1"

	t "github.com/greatbeyond/mailchimp/testing"
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

	s.ctx = NewContextWithToken(context.Background(), os.Getenv("MAILCHIMP_TEST_TOKEN"))
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
