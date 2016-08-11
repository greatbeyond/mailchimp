// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package mailchimp

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"
)

const (
	// CampaignsURL is the url endpoint for campaign on mailchimp v3
	CampaignsURL = "/campaigns"
	// CampaignContentURL is the url endpoint for campaign content on mailchimp v3
	CampaignContentURL = "/content"

	CampaignActionCancel     = "/actions/cancel-send" //	Cancel a campaign
	CampaignActionPause      = "/actions/pause"       //	Pause an RSS-Driven campaign
	CampaignActionResume     = "/actions/resume"      //	Resume an RSS-Driven campaign
	CampaignActionSchedule   = "/actions/schedule"    //	Schedule a campaign
	CampaignActionSend       = "/actions/send"        //	Send a campaign
	CampaignActionTest       = "/actions/test"        //	Send a test email
	CampaignActionUnschedule = "/actions/unschedule"  //	Unschedule a campaign
	// CampaignActionReplicate  = "/actions/replicate"   //	Replicate a campaign ( not implemented )
)

// Campaign defines a campaign on mailchimp
type Campaign struct {
	// A string that uniquely identifies this campaign.
	ID string `                              json:"id"`

	// There are four types of campaigns you can create in MailChimp. A/B Split campaigns have been deprecated and variate campaigns should be used instead.
	Type string `                            json:"type"`

	// The date and time the campaign was created.
	CreateTime string `                      json:"create_time"`

	// The link to the campaign’s archive version.
	ArchiveURL string `                      json:"archive_url"`

	// The original link to the campaign’s archive version.
	LongArchiveURL string `                  json:"long_archive_url"`

	// The current status of the campaign.
	Status string `                          json:"status"`

	// The total number of emails sent for this campaign.
	EmailsSent int `                         json:"emails_sent"`

	// The date and time a campaign was sent.
	SendTime string `                        json:"send_time"`

	// How the campaign’s content is put together (‘template’, ‘drag_and_drop’, ‘html’, ‘url’).
	ContentType string `                     json:"content_type"`

	// List settings for the campaign.
	Recipients CampaignRecipients `          json:"recipients"`

	// The settings for your campaign, including subject, from name, reply-to address, and more.
	Settings CampaignSettings `             json:"settings"`

	// The settings specific to A/B test campaigns.
	VariateSettings interface{} `            json:"variate_settings"`

	// The tracking options for a campaign.
	Tracking CampaignTracking `             json:"tracking"`

	// RSS options for a campaign.
	RSSOpts interface{} `                    json:"rss_opts"`

	// A/B Testing options for a campaign.
	ABSplitOpts interface{} `                json:"ab_split_opts"`

	// The preview for the campaign, rendered by social networks like Facebook and Twitter. Learn more.
	SocialCard interface{} `                 json:"social_card"`

	// For sent campaigns, a summary of opens, clicks, and unsubscribes.
	ReportSummary interface{} `              json:"report_summary"`

	// Updates on campaigns in the process of sending.
	DeliveryStatus CampaignDeliveryStatus ` json:"delivery_status"`

	// Internal
	client MailchimpClient
}

// SetClient fulfills ClientType
func (m *Campaign) SetClient(c MailchimpClient) { m.client = c }

// CampaignSettings defines settings for a campaign
type CampaignSettings struct {
	// The subject line for the campaign.
	SubjectLine string `         json:"subject_line,omitempty"`

	// The title of the campaign.
	Title string `               json:"title,omitempty"`

	// The ‘from’ name on the campaign (not an email address).
	FromName string `            json:"from_name,omitempty"`

	// The reply-to email address for the campaign.
	ReplyTo string `             json:"reply_to,omitempty"`

	// Use MailChimp Conversation feature to manage out-of-office replies.
	UseConversation bool `       json:"use_conversation,omitempty"`

	// The campaign’s custom ‘To’ name. Typically the first name merge field.
	ToName string `              json:"to_name,omitempty"`

	// If the campaign is listed in a folder, the id for that folder.
	FolderID string `            json:"folder_id,omitempty"`

	// Whether MailChimp authenticated the campaign. Defaults to true.
	Authenticate bool `          json:"authenticate,omitempty"`

	// Automatically append MailChimp’s default footer to the campaign.
	AutoFooter bool `            json:"auto_footer,omitempty"`

	// Automatically inline the CSS included with the campaign content.
	InlineCSS bool `             json:"inline_css,omitempty"`

	// Automatically tweet a link to the campaign archive page when the campaign is sent.
	AutoTweet bool `             json:"auto_tweet,omitempty"`

	// An array of Facebook page ids to auto-post to.
	AutoFbPost []string `        json:"auto_fb_post,omitempty"`

	// Allows Facebook comments on the campaign (also force-enables the Campaign Archive toolbar). Defaults to true.
	FbComments bool `            json:"fb_comments,omitempty"`

	// Send this campaign using Timewarp.
	Timewarp bool `              json:"timewarp,omitempty"`

	// The id for the template used in this campaign.
	TemplateID int `             json:"template_id,omitempty"`

	// Whether the campaign uses the drag-and-drop editor.
	DragAndDrop bool `           json:"drag_and_drop,omitempty"`
}

// CampaignTracking settings
type CampaignTracking struct {
	// Whether to track opens. Defaults to true. Cannot be set to false for variate campaigns.
	Opens bool `                 json:"opens,omitempty"`

	// Whether to track clicks in the HTML version of the campaign. Defaults to true. Cannot be set to false for variate campaigns.
	HTMLClicks bool `            json:"html_clicks,omitempty"`

	// Whether to track clicks in the plain-text version of the campaign. Defaults to true. Cannot be set to false for variate campaigns.
	TextClicks bool `            json:"text_clicks,omitempty"`

	// Whether to enable Goal tracking.
	GoalTracking bool `          json:"goal_tracking,omitempty"`

	// Whether to enable E-commerce tracking.
	Ecomm360 bool `              json:"ecomm360,omitempty"`

	// The custom slug for Google Analytics tracking (max of 50 bytes).
	GoogleAnalytics string `     json:"google_analytics,omitempty"`

	// The custom slug for ClickTale tracking (max of 50 bytes).
	Clicktale string `           json:"clicktale,omitempty"`

	// Salesforce tracking options for a campaign. Must be using MailChimp’s built-in Salesforce integration.
	Salesforce interface{} `     json:"salesforce,omitempty"`

	// Highrise tracking options for a campaign. Must be using MailChimp’s built-in Highrise integration.
	Highrise interface{} `       json:"highrise,omitempty"`

	// Capsule tracking options for a campaign. Must be using MailChimp’s built-in Capsule integration.
	Capsule interface{} `        json:"capsule,omitempty"`
}

// CampaignDeliveryStatus updates on campaigns in the process of sending.
type CampaignDeliveryStatus struct {
	Enabled bool `json:"enabled"`
}

// CreateCampaign reference:
// http://developer.mailchimp.com/documentation/mailchimp/reference/campaigns/#
type CreateCampaign struct {
	// There are four types of campaigns you can create
	// in MailChimp. A/B Split campaigns have been
	// deprecated and variate campaigns should be used instead.
	// Possible values: regular, plaintext, absplit, rss, variate
	Type string `                             json:"type"`

	// List settings for the campaign.
	Recipients *CampaignRecipients `         json:"recipients,omitempty"`

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

// CampaignRecipients defines default fields for a reciptient
type CampaignRecipients struct {
	// The unique list id.
	ListID string `              json:"list_id,omitempty"`

	// The name of the list.
	ListName string `            json:"list_name,omitempty"`

	// A string marked-up with HTML explaining the
	// segment used for the campaign in plain English.
	SegmentText string `         json:"segment_text,omitempty"`

	// Count of the recipients on the associated list. Formatted as an integer.
	RecipientCount int `         json:"recipient_count,omitempty"`

	// An object representing all segmentation options.
	SegmentOpts interface{} `    json:"segment_opts,omitempty"`
}

// CampaignCreateSettings Required fields for campaing creation settings
type CampaignCreateSettings struct {
	// The subject line for the campaign.
	SubjectLine string `         json:"subject_line,omitempty"`

	// The title of the campaign.
	Title string `               json:"title,omitempty"`

	// The ‘from’ name on the campaign (not an email address).
	FromName string `            json:"from_name,omitempty"`

	// The reply-to email address for the campaign.
	ReplyTo string `             json:"reply_to,omitempty"`

	// Use MailChimp Conversation feature to manage out-of-office replies.
	UseConversation bool `       json:"use_conversation,omitempty"`

	// The campaign’s custom ‘To’ name. Typically the first name merge field.
	ToName string `              json:"to_name,omitempty"`

	// If the campaign is listed in a folder, the id for that folder.
	FolderID string `            json:"folder_id,omitempty"`

	// Whether MailChimp authenticated the campaign. Defaults to true.
	Authenticate bool `          json:"authenticate,omitempty"`

	// Automatically append MailChimp’s default footer to the campaign.
	AutoFooter bool `            json:"auto_footer,omitempty"`

	// Automatically inline the CSS included with the campaign content.
	InlineCSS bool `             json:"inline_css,omitempty"`

	// Automatically tweet a link to the campaign archive page when the campaign is sent.
	AutoTweet bool `             json:"auto_tweet,omitempty"`

	// An array of Facebook page ids to auto-post to.
	AutoFbPost []string `        json:"auto_fb_post,omitempty"`

	// Allows Facebook comments on the campaign (also force-enables the Campaign Archive toolbar). Defaults to true.
	FbComments bool `            json:"fb_comments,omitempty"`
}

// NewCampaign creates a new campaign with the client addressed
// id is optional, with it you can do a bit of rudimentary chaining.
// Example:
//	c.NewCampaign(23).Update(params)
func (c *Client) NewCampaign(id ...string) *Campaign {
	s := &Campaign{
		client: c,
	}
	if len(id) > 0 {
		s.ID = id[0]
	}
	return s
}

// CreateCampaign creates a new campaign via mailchimp api v3
func (c *Client) CreateCampaign(data *CreateCampaign) (*Campaign, error) {
	response, err := c.Post(CampaignsURL, nil, data)

	var campaign *Campaign
	err = json.Unmarshal(response, &campaign)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	campaign.client = c

	return campaign, nil
}

// -----------------------------------------------------------------
// Retrive and update

type getCampaigns struct {
	Campaigns  []*Campaign `json:"campaigns"`
	ListID     string      `json:"list_id"`
	TotalItems int         `json:"total_items"`
}

// GetCampaigns retrives all campaigns from mailchimp
func (c *Client) GetCampaigns() ([]*Campaign, error) {
	response, err := c.Get(CampaignsURL, nil)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	var campaignsResponse *getCampaigns
	err = json.Unmarshal(response, &campaignsResponse)
	if err != nil {
		return nil, err
	}

	// Add internal client
	campaigns := []*Campaign{}
	for _, campaign := range campaignsResponse.Campaigns {
		campaign.client = c
		campaigns = append(campaigns, campaign)
	}

	return campaigns, nil
}

// GetCampaign retrives a single campaign by id
func (c *Client) GetCampaign(id string) (*Campaign, error) {
	response, err := c.Get(slashJoin(CampaignsURL, id), nil)

	var campaign *Campaign
	err = json.Unmarshal(response, &campaign)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"ID":    id,
			"error": err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	campaign.client = c

	return campaign, nil
}

// Update sets new values on a campaign via mailchimp api
func (c *Campaign) Update(data *UpdateCampaign) (*Campaign, error) {
	response, err := c.client.Patch(slashJoin(CampaignsURL, c.ID), nil, data)

	var campaign *Campaign
	err = json.Unmarshal(response, &campaign)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"ID":    c.ID,
			"error": err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	campaign.client = c.client

	return campaign, nil
}

// -----------------------------------------------------------------
// Actions on campaign

// Cancel a campaign
func (c *Campaign) Cancel() error {
	_, err := c.client.Post(slashJoin(CampaignsURL, c.ID, CampaignActionCancel), nil, nil)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"ID":    c.ID,
			"error": err.Error(),
		}).Error("response error", caller())
		return err
	}
	return nil
}

// Pause a campaign
func (c *Campaign) Pause() error {
	_, err := c.client.Post(slashJoin(CampaignsURL, c.ID, CampaignActionPause), nil, nil)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"ID":    c.ID,
			"error": err.Error(),
		}).Error("response error", caller())
		return err
	}
	return nil
}

// Resume a campaign
func (c *Campaign) Resume() error {
	_, err := c.client.Post(slashJoin(CampaignsURL, c.ID, CampaignActionResume), nil, nil)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"ID":    c.ID,
			"error": err.Error(),
		}).Error("response error", caller())
		return err
	}
	return nil
}

// Schedule a campaign
func (c *Campaign) Schedule() error {
	_, err := c.client.Post(slashJoin(CampaignsURL, c.ID, CampaignActionSchedule), nil, nil)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"ID":    c.ID,
			"error": err.Error(),
		}).Error("response error", caller())
		return err
	}
	return nil
}

// Send a campaign
func (c *Campaign) Send() error {
	_, err := c.client.Post(slashJoin(CampaignsURL, c.ID, CampaignActionSend), nil, nil)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"ID":    c.ID,
			"error": err.Error(),
		}).Error("response error", caller())
		return err
	}
	return nil
}

// Test a campaign
func (c *Campaign) Test() error {
	_, err := c.client.Post(slashJoin(CampaignsURL, c.ID, CampaignActionTest), nil, nil)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"ID":    c.ID,
			"error": err.Error(),
		}).Error("response error", caller())
		return err
	}
	return nil
}

// Unschedule a campaign
func (c *Campaign) Unschedule() error {
	_, err := c.client.Post(slashJoin(CampaignsURL, c.ID, CampaignActionUnschedule), nil, nil)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"ID":    c.ID,
			"error": err.Error(),
		}).Error("response error", caller())
		return err
	}
	return nil
}

// Delete removes a campaign
func (c *Campaign) Delete() error {
	err := c.client.Delete(slashJoin(CampaignsURL, c.ID))
	if err != nil {
		Log.WithFields(logrus.Fields{
			"ID":    c.ID,
			"error": err.Error(),
		}).Error("response error", caller())
		return err
	}
	return nil
}

// -----------------------------------------------------------------
// Content manipulation on the campaign

// CampaignContent can be retrived
type CampaignContent struct {
	// Content options for multivariate campaigns.
	VariateContents struct {
		// Label used to identify the content option.
		ContentLabel string `     json:"content_label,omitempty"`

		// The plain-text portion of the campaign. If left unspecified, we’ll generate this automatically.
		PlainText string `        json:"plain_text,omitempty"`

		// The raw HTML for the campaign.
		HTML string `             json:"html,omitempty"`
	} `json:"variate_contents,omitempty"`

	// The plain-text portion of the campaign. If left unspecified, we’ll generate this automatically.
	PlainText string `            json:"plain_text,omitempty"`

	// The raw HTML for the campaign.
	HTML string `                 json:"html,omitempty"`
}

// CampaignContentEdit is documented here:
// http://developer.mailchimp.com/documentation/mailchimp/reference/campaigns/content/
type CampaignContentEdit struct {
	// The plain-text portion of the campaign. If left unspecified, we’ll generate this automatically.
	PlainText string `          json:"plain_text,omitempty"`

	// The raw HTML for the campaign.
	HTML string `               json:"html,omitempty"`

	// When importing a campaign, the URL where the HTML lives.
	URL string `                json:"url,omitempty"`

	// Use this template to generate the HTML content of the campaign
	Template struct {
		// The id of the template to use.
		ID int `                json:"id,omitempty"`

		// Content for the sections of the template. Each key should be the unique mc:edit area name from the template.
		Sections interface{} `  json:"sections,omitempty"`
	} `json:"template,omitempty"`

	// Available when uploading an archive to create campaign content. The archive should include all campaign content and images. Learn more.
	Archive struct {
		// he base64-encoded representation of the archive file.
		ArchiveContent string `json:"archive_content,omitempty"`

		// The type of encoded file. Defaults to zip.
		// Possible Values:
		// 	zip tar.gz tar.bz2 tar tgz tbz
		ArchiveType string `   json:"archive_type,omitempty"`
	} `json:"archive,omitempty"`

	// Content options for Multivariate Campaigns. Each content option must provide HTML content and may optionally provide plain text. For campaigns not testing content, only one object should be provided.
	VariateContents []interface{} `json:"variate_contents,omitempty"`
}

// GetContent retrives the content for a campaign
func (c *Campaign) GetContent() (interface{}, error) {

	response, err := c.client.Get(slashJoin(CampaignsURL, c.ID, CampaignContentURL), nil)

	var content *CampaignContent
	err = json.Unmarshal(response, &content)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"ID":    c.ID,
			"error": err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	return content, nil
}

// SetContent updates the content for the campaign
func (c *Campaign) SetContent(content *CampaignContentEdit) (*CampaignContent, error) {
	response, err := c.client.Put(slashJoin(CampaignsURL, c.ID, CampaignContentURL), nil, content)

	var responseContent *CampaignContent
	err = json.Unmarshal(response, &responseContent)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"ID":    c.ID,
			"error": err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	return responseContent, nil
}
