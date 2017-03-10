// Copyright (C) 2017 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by Sonny Vidfamn <sonny.vidfamn@greatbeyond.se>, 2017

package mailchimp

import (
	"context"
	"encoding/json"

	"github.com/Sirupsen/logrus"
)

const ReportURL = "/reports"
const SentToURL = "/sent-to"

// SentTo contains sent information about a single member for a campaign.
type SentTo struct {
	EmailID      string          `json:"email_id,omitempty"`
	EmailAddress string          `json:"email_address,omitempty"`
	MergeFields  interface{}     `json:"merge_fields,omitempty"`
	VIP          bool            `json:"vip,omitempty"`
	Status       string          `json:"status,omitempty"`
	OpenCount    int             `json:"open_count,omitempty"`
	LastOpen     string          `json:"last_open,omitempty"`
	ABSplitGroup json.RawMessage `json:"absplit_group,omitempty"`
	GMTOffset    int             `json:"gmt_offset,omitempty"`
	CampaignID   string          `json:"campaign_id,omitempty"`
	ListID       string          `json:"list_id,omitempty"`
	Links        json.RawMessage `json:"_links,omitempty"`
}

// GetSentTo contains sent information about members for a campaign.
// https://developer.mailchimp.com/documentation/mailchimp/reference/reports/sent-to/
type GetSentTo struct {
	SentTo     []SentTo        `json:"sent_to,omitempty"`
	CampaignID int             `json:"campaign_id,omitempty"`
	TotalItems int             `json:"total_items,omitempty"`
	Links      json.RawMessage `json:"_links,omitempty"`
}

// GetSentTo returns sent status for each member in a sent campaign.
// Optional params: fields, exclude_fields, count, offset.
// See: https://developer.mailchimp.com/documentation/mailchimp/reference/reports/sent-to/
func (c *Client) GetSentTo(ctx context.Context, campaignID string, params ...Parameters) (*GetSentTo, error) {
	p := requestParameters(params)
	response, err := c.Get(ctx, slashJoin(ReportURL, campaignID, SentToURL), p)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"campaign_id": campaignID,
			"error":       err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	var sentToResponse *GetSentTo
	err = json.Unmarshal(response, &sentToResponse)
	if err != nil {
		return nil, err
	}

	return sentToResponse, nil
}
