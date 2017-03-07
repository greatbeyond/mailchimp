// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by Sonny Vidfamn <sonny.vidfamn@greatbeyond.se>, 2016

package mailchimp

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/schema"
)

// ------------------------------------------------------------------------------
// Common Webhook definitions
// ------------------------------------------------------------------------------

// WebhooksURL is REST url for webhooks
const WebhooksURL = "/webhooks"

// Webhook defines the structure of webhook from mailchimp.
type Webhook struct {
	// An string that uniquely identifies this webhook.
	ID string `json:"id"`

	// A valid URL for the Webhook.
	URL string `json:"url"`

	// The events that can trigger the webhook and whether they are enabled.
	Events *WebhookEvents `json:"events"`

	// The possible sources of any events that can trigger the webhook and whether they are enabled.
	Sources *WebhookSources `json:"sources"`

	// The unique id for the list.
	ListID string `json:"list_id"`

	// This field has links to schema types
	Links json.RawMessage `json:"_links"`

	// Internal
	client MailchimpClient
}

// SetClient fulfills ClientType
func (w *Webhook) SetClient(c MailchimpClient) { w.client = c }

// WebhookEvents defines all valid fields for webhook events.
type WebhookEvents struct {
	// Whether the webhook is triggered when a list subscriber is added.
	Subscribe bool `json:"subscribe,omitempty"`

	// Whether the webhook is triggered when a list member unsubscribes.
	Unsubscribe bool `json:"unsubscribe,omitempty"`

	// Whether the webhook is triggered when a subscriber’s profile is updated.
	Profile bool `json:"profile,omitempty"`

	// Whether the webhook is triggered when a subscriber’s email address is cleaned from
	Cleaned bool `json:"cleaned,omitempty"`

	// Whether the webhook is triggered when a subscriber’s email address is changed.
	UpEmail bool `json:"upemail,omitempty"`

	// Whether the webhook is triggered when a campaign is sent or cancelled.
	Campaign bool `json:"campaign,omitempty"`
}

// WebhookSources defines all valid fields for webhook sources.
type WebhookSources struct {
	//Whether the webhook is triggered by subscriber-initiated actions.
	User bool `json:"user,omitempty"`
	//Whether the webhook is triggered by admin-initiated actions in the web interface.
	Admin bool `json:"admin,omitempty"`
	//Whether the webhook is triggered by actions initiated via the API.
	API bool `json:"api,omitempty"`
}

// ------------------------------------------------------------------------------
// Webhook request, response definitions and implementation
// ------------------------------------------------------------------------------

// CreateWebhook defines the structure of a create webhook request to mailchimp.
type CreateWebhook struct {
	ListID  string          `json:"-"` // json marshal ignore
	URL     string          `json:"url"`
	Events  *WebhookEvents  `json:"events"`
	Sources *WebhookSources `json:"sources"`
}

// CreateWebhook adds a webhook to a list. Mailchimp will send events through this webhook on:
// subcribes, unsubscribes, profile updates, email address changes and campaign sending status.
// Returns webhook ID on success, otherwise error.
func (c *Client) CreateWebhook(ctx context.Context, request *CreateWebhook) (*Webhook, error) {
	_, err := c.GetList(ctx, request.ListID)
	if err != nil {
		return nil, err
	}

	response, err := c.Post(ctx, slashJoin(ListsURL, request.ListID, WebhooksURL), nil, request)

	var webhook *Webhook
	err = json.Unmarshal(response, &webhook)
	if err != nil {
		return nil, err
	}

	// Add internal client
	webhook.client = c

	return webhook, nil
}

// getWebhooksResponse defines the structure of a info webhooks response from mailchimp.
type getWebhooksResponse struct {
	Webhooks   []*Webhook      `json:"webhooks"`
	ListID     string          `json:"list_id"`
	TotalItems int             `json:"total_items"`
	Links      json.RawMessage `json:"_links"` // This field has links to schema types
}

// GetWebhooks returns information from all webhooks on a list.
// Returns webhook info on success, otherwise nil and error.
func (c *Client) GetWebhooks(ctx context.Context, listID string) ([]*Webhook, error) {
	response, err := c.Get(ctx, slashJoin(ListsURL, listID, WebhooksURL), nil)
	if err != nil {
		Log.Error(err.Error(), caller())
		return nil, err
	}

	var webhooksResponse *getWebhooksResponse
	err = json.Unmarshal(response, &webhooksResponse)
	if err != nil {
		Log.Error(err.Error(), caller())
		return nil, err
	}

	// Add internal client
	webhooks := []*Webhook{}
	for _, webhook := range webhooksResponse.Webhooks {
		webhook.client = c
		webhooks = append(webhooks, webhook)
	}

	return webhooks, nil
}

// GetWebhook returns information for a single webhook.
func (c *Client) GetWebhook(ctx context.Context, listID string, webhookID string) (*Webhook, error) {
	response, err := c.Get(ctx, slashJoin(ListsURL, listID, WebhooksURL, webhookID), nil)
	if err != nil {
		Log.Error(err.Error(), caller())
		return nil, err
	}

	var webhook *Webhook
	err = json.Unmarshal(response, &webhook)
	if err != nil {
		Log.Error(err.Error(), caller())
		return nil, err
	}

	// Add internal client
	webhook.client = c

	return webhook, nil
}

// DeleteWebhook removes a webhook from mailchimp.
// Returns error on failure
func (w *Webhook) DeleteWebhook(ctx context.Context) error {
	if w.client == nil {
		return ErrorNoClient
	}
	return w.client.Delete(ctx, slashJoin(ListsURL, w.ListID, WebhooksURL, w.ID))
}

// ------------------------------------------------------------------------------
// Mailchimp webhook events structs and parse
// ------------------------------------------------------------------------------
const (
	// WebhookEventTypeSubscribe the type of subscribe event
	WebhookEventTypeSubscribe string = "subscribe"
	// WebhookEventTypeUnsubscribe the type of unsubscribe event
	WebhookEventTypeUnsubscribe string = "unsubscribe"
	// WebhookEventTypeProfileUpdates the type of profile event
	WebhookEventTypeProfileUpdates string = "profile"
	// WebhookEventTypeEmailChanged the type of upemail event
	WebhookEventTypeEmailChanged string = "upemail"
	// WebhookEventTypeEmailCleaned the type of cleaned event
	WebhookEventTypeEmailCleaned string = "cleaned"
	// WebhookEventTypeCampaignStatus the type of campaign event
	WebhookEventTypeCampaignStatus string = "campaign"
)

// WebhookEvent is a canonical struct with fields for all supported events. See API reference
// for which fields is sent by which event: https://apidocs.mailchimp.com/webhooks/.
type WebhookEvent struct {
	Type       string `schema:"type"`
	FiredAt    string `schema:"fired_at"`
	Action     string `schema:"data[action]"`
	Reason     string `schema:"data[reason]"`
	ID         string `schema:"data[id]"`
	NewID      string `schema:"data[new_id]"`
	ListID     string `schema:"data[list_id]"`
	Email      string `schema:"data[email]"`
	EmailType  string `schema:"data[email_type]"`
	NewEmail   string `schema:"data[new_email]"`
	OldEmail   string `schema:"data[old_email]"`
	IPOpt      string `schema:"data[ip_opt]"`
	IPSignup   string `schema:"data[ip_signup]"`
	CampaignID string `schema:"data[campaign_id]"`
	Subject    string `schema:"data[subject]"`
	Status     string `schema:"data[status]"`

	mergeFields map[string]string `schema:"-"`
}

// GetMergesField returns the value of the merge field if it exists and was defined in WebhookParseEvent.
func (e *WebhookEvent) GetMergesField(field string) string {
	if e.mergeFields == nil {
		return ""
	}

	return e.mergeFields[field]
}

// WebhookParseEvent will parse the webhook form event and return a WebhookEvent on success,
// otherwise error.
func WebhookParseEvent(r *http.Request, mergesFields ...string) (*WebhookEvent, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	event := &WebhookEvent{}
	err = decoder.Decode(event, r.Form)
	if err != nil {
		return nil, err
	}

	// Save defined merges fields
	if event.mergeFields == nil {
		event.mergeFields = map[string]string{}
	}

	for _, v := range mergesFields {
		event.mergeFields[v] = r.Form.Get(v)
	}

	return event, nil
}
