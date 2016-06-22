package mailchimp

import "encoding/json"

// ------------------------------------------------------------------------------
// Common Webhook definitions
// ------------------------------------------------------------------------------

// WebhooksURL is REST url for webhooks
const WebhooksURL = "/webhooks"

// Webhook defines the structure of webhook from mailchimp.
type Webhook struct {
	ID      string          `json:"id"`
	URL     string          `json:"url"`
	Events  WebhookEvents   `json:"events"`
	Sources WebhookSources  `json:"sources"`
	ListID  string          `json:"list_id"`
	Links   json.RawMessage `json:"_links"` // This field has links to schema types

	// Internal
	client *Client
}

// WebhookEvents defines all valid fields for webhook events.
type WebhookEvents struct {
	Subscribe   bool `json:"subscribe,omitempty"`
	Unsubscribe bool `json:"unsubscribe,omitempty"`
	Profile     bool `json:"profile,omitempty"`
	Cleaned     bool `json:"cleaned,omitempty"`
	UpEmail     bool `json:"upemail,omitempty"`
	Campaign    bool `json:"campaign,omitempty"`
}

// WebhookSources defines all valid fields for webhook sources.
type WebhookSources struct {
	User  bool `json:"user,omitempty"`
	Admin bool `json:"admin,omitempty"`
	API   bool `json:"api,omitempty"`
}

// ------------------------------------------------------------------------------
// Webhook request, response definitions and implementation
// ------------------------------------------------------------------------------

// CreateWebhookRequest defines the structure of a create webhook request to mailchimp.
type CreateWebhookRequest struct {
	ListID  string         `json:"-"` // json marshal ignore
	URL     string         `json:"url"`
	Events  WebhookEvents  `json:"events"`
	Sources WebhookSources `json:"sources"`
}

// CreateWebhookResponse defines the structure of a create webhook response from mailchimp.
type CreateWebhookResponse Webhook

// CreateWebhook adds a webhook to a list. Mailchimp will send events through this webhook on:
// subcribes, unsubscribes, profile updates, email address changes and campaign sending status.
// Returns webhook ID on success, otherwise error.
func (c *Client) CreateWebhook(request *CreateWebhookRequest) (*CreateWebhookResponse, error) {
	_, err := c.GetList(request.ListID)
	if err != nil {
		return nil, err
	}

	response, err := c.Post(slashJoin(ListsURL, request.ListID, WebhooksURL), nil, request)

	createWebhookResponse := CreateWebhookResponse{}
	err = json.Unmarshal(response, &createWebhookResponse)
	if err != nil {
		return nil, err
	}

	// Add internal client
	createWebhookResponse.client = c

	return &createWebhookResponse, nil
}

// DeleteWebhook removes a webhook from mailchimp.
// Returns error on failure
func (w *Webhook) DeleteWebhook() error {
	if w.client == nil {
		return ErrorNoClient
	}

	return w.client.Delete(slashJoin(ListsURL, w.ListID, WebhooksURL, w.ID))
}

// GetWebhooksResponse defines the structure of a info webhooks response from mailchimp.
type GetWebhooksResponse struct {
	Webhooks   []Webhook       `json:"webhooks"`
	ListID     string          `json:"list_id"`
	TotalItems int             `json:"total_items"`
	Links      json.RawMessage `json:"_links"` // This field has links to schema types
}

// GetWebhooks returns information from all webhooks on a list.
// Returns webhook info on success, otherwise nil and error.
func (c *Client) GetWebhooks(listID string) (*GetWebhooksResponse, error) {
	response, err := c.Get(slashJoin(ListsURL, listID, WebhooksURL), nil)
	if err != nil {
		return nil, err
	}

	getWebhooksResponse := GetWebhooksResponse{}
	err = json.Unmarshal(response, &getWebhooksResponse)
	if err != nil {
		return nil, err
	}

	// Add internal client
	for _, webhook := range getWebhooksResponse.Webhooks {
		webhook.client = c
	}

	return &getWebhooksResponse, nil
}

// GetWebhookResponse defines the structure of a info webhook response from mailchimp.
type GetWebhookResponse Webhook

// GetWebhook returns information for a single webhook.
func (c *Client) GetWebhook(listID string, webhookID string) (*GetWebhookResponse, error) {
	response, err := c.Get(slashJoin(ListsURL, listID, WebhooksURL, webhookID), nil)
	if err != nil {
		return nil, err
	}

	getWebhookResponse := GetWebhookResponse{}
	err = json.Unmarshal(response, &getWebhookResponse)
	if err != nil {
		return nil, err
	}

	// Add internal client
	getWebhookResponse.client = c

	return &getWebhookResponse, nil
}

// DeleteWebhook removes a webhook from mailchimp.
// Returns error on failure
func (w *GetWebhookResponse) DeleteWebhook() error {
	if w.client == nil {
		return ErrorNoClient
	}

	return w.client.Delete(slashJoin(ListsURL, w.ListID, WebhooksURL, w.ID))
}
