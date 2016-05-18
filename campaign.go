// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package mailchimp

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	"github.com/antonholmquist/jason"
)

// CAMPAIGNS_URL is the url endpoint for campaign on mailchimp v3
const CAMPAIGNS_URL = "/campaigns"

// Campaign defines a campaign on mailchimp
type Campaign struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	CreateTime  string `json:"create_time"`
	ArchiveURL  string `json:"archive_url"`
	Status      string `json:"status"`
	EmailsSent  int    `json:"emails_sent"`
	SendTime    string `json:"send_time"`
	ContentType string `json:"content_type"`
	Recipients  string `json:"recipients"`
	ListID      string `json:"list_id"`
	SegmentText string `json:"segment_text"`

	Settings       *CampaignSettings       `json:"settings"`
	Tracking       *CampaignTracking       `json:"tracking"`
	DeliveryStatus *CampaignDeliveryStatus `json:"delivery_status"`

	// Internal
	client *Client
}

// CampaignSettings defines settings for a campaign
type CampaignSettings struct {
	SubjectLine     string `json:"subject_line"`
	Title           string `json:"title"`
	FromName        string `json:"from_name"`
	ReplyTo         string `json:"reply_to"`
	UseConversation bool   `json:"use_conversation"`
	ToName          string `json:"to_name"`
	FolderID        string `json:"folder_id"`
	Authenticate    bool   `json:"authenticate"`
	AutoFooter      bool   `json:"auto_footer"`
	InlineCSS       bool   `json:"inline_css"`
	AutoTweet       bool   `json:"auto_tweet"`
	FbComments      bool   `json:"fb_comments"`
	Timewarp        bool   `json:"timewarp"`
	TemplateID      int    `json:"template_id"`
	DragAndDrop     bool   `json:"drag_and_drop"`
}

type CampaignTracking struct {
	Opens           bool   `json:"opens"`
	HTMLClicks      bool   `json:"html_clicks"`
	TextClicks      bool   `json:"text_clicks"`
	GoalTracking    bool   `json:"goal_tracking"`
	Ecomm360        bool   `json:"ecomm360"`
	GoogleAnalytics string `json:"google_analytics"`
	Clicktale       string `json:"clicktale"`
}

type CampaignDeliveryStatus struct {
	Enabled bool `json:"enabled"`
}

// CreateCampaign reference:
// http://developer.mailchimp.com/documentation/mailchimp/reference/campaigns/#
type CreateCampaign struct {
	// There are four types of campaigns you can create
	// in MailChimp. A/B Split campaigns have been
	// deprecated and variate campaigns should be used instead.
	Type string `                             json:"type"`

	// List settings for the campaign.
	Recipients []*CampaignRecipient `         json:"recipients,omitempty"`

	// The settings for your campaign, including subject,
	// from name, reply-to address, and more.
	// If you only need the required fields, use CampaignCreateSettings.
	// If you need to include more fields, create your own struct and
	// put it here
	Settings interface{} `                    json:"settings,omitempty"`

	// The settings specific to A/B test campaigns.
	VariateSettings interface{} `             json:"variate_settings,omitempty"`

	// The tracking options for a campaign.
	Tracking interface{} `                    json:"tracking,omitempty"`

	// RSS options for a campaign.
	RssOpts interface{} `                     json:"rss_opts,omitempty"`

	// A/B Testing options for a campaign.
	ABSplitOpts interface{} `                 json:"ab_split_opts,omitempty"`

	// The preview for the campaign, rendered by social networks
	// like Facebook and Twitter. Learn more.
	SocialCard interface{} `                  json:"social_card,omitempty"`

	// For sent campaigns, a summary of opens, clicks, and unsubscribes.
	ReportSummary interface{} `               json:"report_summary,omitempty"`

	// Updates on campaigns in the process of sending.
	DeliveryStatus *CampaignDeliveryStatus `  json:"delivery_status,omitempty"`
}

// UpdateCampaign and CreateCampaign are the same but with
// different requiered fields (checked in function)
type UpdateCampaign CreateCampaign

// CampaignRecipient defines default fields for a reciptient
type CampaignRecipient struct {
	// The unique list id.
	ListID string `json:"list_id"`

	// A string marked-up with HTML explaining the
	// segment used for the campaign in plain English.
	SegmentText string `json:"segment_text"`

	// An object representing all segmentation options.
	SegmentOpts interface{} `json:"segment_opts"`
}

// CampaignCreateSettings Required fields for campaing creation settings
type CampaignCreateSettings struct {
	// The subject line for the campaign.
	SubjectLine string `json:"subject_line"`

	// The title of the campaign.
	Title string `json:"title"`

	// The ‘from’ name on the campaign (not an email address).
	FromName string `json:"from_name"`

	// The reply-to email address for the campaign.
	ReplyTo string `json:"reply_to"`

	// Use MailChimp Conversation feature to manage out-of-office replies.
	UseConversation bool `json:"use_conversation"`

	// The campaign’s custom ‘To’ name. Typically the first name merge field.
	ToName string `json:"to_name"`

	// If the campaign is listed in a folder, the id for that folder.
	FolderID string `json:"folder_id"`

	// Whether MailChimp authenticated the campaign. Defaults to true.
	Authenticate bool `json:"authenticate"`

	// Automatically append MailChimp’s default footer to the campaign.
	AutoFooter bool `json:"auto_footer"`

	// Automatically inline the CSS included with the campaign content.
	InlineCSS bool `json:"inline_css"`

	// Automatically tweet a link to the campaign archive page when the campaign is sent.
	AutoTweet bool `json:"auto_tweet"`

	// An array of Facebook page ids to auto-post to.
	AutoFbPost []string `json:"auto_fb_post"`

	// Allows Facebook comments on the campaign (also force-enables the Campaign Archive toolbar). Defaults to true.
	FbComments bool `json:"fb_comments"`
}

// NewCampaign creates a new campaign via mailchimp api v3
func (c *Client) NewCampaign(data *CreateCampaign) (*Campaign, error) {
	response, err := c.post(CAMPAIGNS_URL, nil, data)

	var campaign *Campaign
	err = json.Unmarshal(response, &campaign)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("response error", caller())
		return nil, err
	}

	campaign.client = c

	return campaign, nil
}

// GetCampaigns retrives all campaigns from mailchimp
func (c *Client) GetCampaigns() ([]*Campaign, error) {
	response, err := c.get(CAMPAIGNS_URL, nil)

	v, err := jason.NewObjectFromBytes(response)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("response error", caller())
		return nil, err
	}

	_lists, err := v.GetValue("campaigns")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("response error", caller())
		return nil, err
	}

	b, err := _lists.Marshal()
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("response error", caller())
		return nil, err
	}

	var campaigns []*Campaign
	err = json.Unmarshal(b, &campaigns)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("response error", caller())
		return nil, err
	}

	for _, l := range campaigns {
		l.client = c
	}

	return campaigns, nil
}

// GetCampaign retrives a single campaign by id
func (c *Client) GetCampaign(id string) (*Campaign, error) {
	response, err := c.get(slashJoin(CAMPAIGNS_URL, id), nil)

	var campaign *Campaign
	err = json.Unmarshal(response, &campaign)
	if err != nil {
		log.WithFields(log.Fields{
			"ID":    id,
			"error": err.Error(),
		}).Warn("response error", caller())
		return nil, err
	}

	campaign.client = c

	return campaign, nil
}

// Update sets new values on a campaign via mailchimp api
func (c *Campaign) Update(data *UpdateCampaign) (*Campaign, error) {
	response, err := c.client.patch(slashJoin(CAMPAIGNS_URL, c.ID), nil, data)

	var campaign *Campaign
	err = json.Unmarshal(response, &campaign)
	if err != nil {
		log.WithFields(log.Fields{
			"ID":    c.ID,
			"error": err.Error(),
		}).Warn("response error", caller())
		return nil, err
	}

	campaign.client = c.client

	return campaign, nil
}

// Delete removes a campaign
func (c *Campaign) Delete() error {
	err := c.client.delete(slashJoin(CAMPAIGNS_URL, c.ID))
	if err != nil {
		log.WithFields(log.Fields{
			"ID":    c.ID,
			"error": err.Error(),
		}).Warn("response error", caller())
		return err
	}
	return nil
}
