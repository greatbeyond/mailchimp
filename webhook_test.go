// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package mailchimp

import check "gopkg.in/check.v1"

var _ = check.Suite(&WebhookTestSuite{})

type WebhookTestSuite struct {
}

func (suite *WebhookTestSuite) NewClient() *Client {
	client := NewClient("arandomtoken-us0")
	client.NewBatch()

	return client
}

func (suite *WebhookTestSuite) SetUpSuite(c *check.C) {}

func (suite *WebhookTestSuite) SetUpTest(c *check.C) {
}

func (suite *WebhookTestSuite) TearDownTest(c *check.C) {}

func (suite *WebhookTestSuite) Test_CreateWebhook(c *check.C) {
	client := suite.NewClient()

	createWebhookResponse, err := client.CreateWebhook(&CreateWebhook{
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

func (suite *WebhookTestSuite) Test_GetWebhook(c *check.C) {
	client := suite.NewClient()

	getWebhookResponse, err := client.GetWebhook("1", "2")
	c.Assert(err, check.IsNil)
	c.Assert(getWebhookResponse, check.NotNil)
}

func (suite *WebhookTestSuite) Test_GetWebhooks(c *check.C) {
	client := suite.NewClient()

	getWebhooksResponse, err := client.GetWebhooks("1")
	c.Assert(err, check.IsNil)
	c.Assert(getWebhooksResponse, check.NotNil)
}

func (suite *WebhookTestSuite) Test_DeleteWebhook(c *check.C) {
	client := suite.NewClient()

	getWebhookResponse, err := client.GetWebhook("1", "2")
	c.Assert(err, check.IsNil)
	c.Assert(getWebhookResponse, check.NotNil)

	err = getWebhookResponse.DeleteWebhook()
	c.Assert(err, check.IsNil)
}
