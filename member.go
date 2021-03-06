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
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

const MembersURL = "/members"

type MailType string

const (
	HTML         MailType = "html"
	MailTypeText MailType = "text"
)

type MemberStatus string

const (
	Subscribed   MemberStatus = "subscribed"
	Unsubscribed MemberStatus = "unsubscribed"
	Cleaned      MemberStatus = "cleaned"
	Pending      MemberStatus = "pending"
)

// Member manages members of a specific MailChimp list, including currently subscribed, unsubscribed, and bounced members.
// http://developer.mailchimp.com/documentation/mailchimp/reference/lists/members/#
type Member struct {
	// The MD5 hash of the lowercase version of the list member’s email address.
	ID string `                         json:"id,omitempty"`
	// Email address for a subscriber.
	EmailAddress string `               json:"email_address,omitempty"`
	// An identifier for the address across all of MailChimp.
	UniqueEmailID string `              json:"unique_email_id,omitempty"`
	// Type of email this member asked to get (‘html’ or ‘text’).
	EmailType MailType `                json:"email_type,omitempty"`
	// Subscriber’s current status.
	Status MemberStatus `               json:"status,omitempty"`
	// An individual merge var and value for a member.
	MergeFields map[string]interface{} `json:"merge_fields,omitempty"`
	// The key of this object’s properties is the ID of the interest in question.
	Interests map[string]bool `         json:"interests,omitempty"`
	// Open and click rates for this subscriber.
	Stats MemberStats `                json:"stats,omitempty"`
	// IP address the subscriber signed up from.
	IPSignup string `                   json:"ip_signup,omitempty"`
	// The date and time the subscriber signed up for the list.
	TimestampSignup string `            json:"timestamp_signup,omitempty"`
	// The IP address the subscriber used to confirm their opt-in status.
	IPOpt string `                      json:"ip_opt,omitempty"`
	// The date and time the subscribe confirmed their opt-in status.
	TimestampOpt string `               json:"timestamp_opt,omitempty"`
	// Star rating for this member, between 1 and 5.
	MemberRating int `                  json:"member_rating,omitempty"`
	// The date and time the member’s info was last changed.
	LastChanged string `                json:"last_changed,omitempty"`
	// If set/detected, the subscriber’s language.
	Language string `                   json:"language,omitempty"`
	// VIP status for subscriber.
	Vip bool `                          json:"vip,omitempty"`
	// The list member’s email client.
	EmailClient string `                json:"email_client,omitempty"`
	// Subscriber location information.
	Location Location `                json:"location,omitempty"`
	// The most recent Note added about this member.
	LastNote map[string]interface{} `   json:"last_note,omitempty"`
	// The list id.
	ListID string `                     json:"list_id,omitempty"`

	// Internal
	Client MailchimpClient `json:"-"`
}

// SetClient fulfills ClientType
func (m *Member) SetClient(c MailchimpClient) { m.Client = c }

// CreateMember contains fields to create or update memebrs
type CreateMember struct {
	// Email address for a subscriber. (required)
	EmailAddress string `           json:"email_address,omitempty"`

	// Type of email this member asked to get (‘html’ or ‘text’).
	EmailType MailType `            json:"email_type,omitempty"`

	// Subscriber’s current status. (Required) Possible Values:
	// subscribed, unsubscribed, cleaned, pending
	Status MemberStatus `           json:"status,omitempty"`

	// An individual merge var and value for a member.
	MergeFields map[string]interface{} `json:"merge_fields,omitempty"`

	// The key of this object’s properties is the ID of the interest in question.
	Interests map[string]bool `     json:"interests,omitempty"`

	// If set/detected, the subscriber’s language.
	Language string `               json:"language,omitempty"`

	// VIP status for subscriber.
	// http://kb.mailchimp.com/lists/managing-subscribers/designate-and-send-to-vip-subscribers
	Vip bool `                      json:"vip,omitempty"`

	// Subscriber location information.
	Location *Location `            json:"location,omitempty"`
}

// UpdateMember and CreateMember are the same but with slighlty
// different requiered fields (checked in function)
type UpdateMember CreateMember

type MemberStats struct {
	// A subscriber’s average open rate.
	AvgOpenRate float64 `json:"avg_open_rate,omitempty"`
	// A subscriber’s average clickthrough rate.
	AvgClickRate float64 `json:"avg_click_rate,omitempty"`
}

// Location points to a geo location and time zone
type Location struct {
	// The location latitude.
	Latitude float64 `json:"latitude,omitempty"`
	// The location longitude.
	Longitude float64 `json:"longitude,omitempty"`
	// The time difference in hours from GMT.
	GmtOff int `json:"gmtoff,omitempty"`
	// The offset for timezones where daylight saving time is observed.
	DstOff int `json:"dstoff,omitempty"`
	// The unique code for the location country.
	CountryCode string `json:"country_code,omitempty"`
	// The timezone for the location.
	Timezone string `json:"timezone,omitempty"`
}

// NewMember returns a empty member object
// id is optional, with it you can do a bit of rudimentary chaining.
// Example:
//	c.NewMember(23).Update(params)
func (c *Client) NewMember(listID string, id ...string) *Member {
	s := &Member{
		Client: c,
	}
	if len(id) > 0 {
		s.ID = id[0]
	}

	s.ListID = listID

	return s
}

// CreateMember Creates a member object and inserts it
func (c *Client) CreateMember(ctx context.Context, data *CreateMember, listID string) (*Member, error) {
	if listID == "" {
		return nil, fmt.Errorf("missing argument: listID")
	}

	if err := hasFields(*data, "EmailAddress", "Status"); err != nil {
		Log.Info(err.Error(), caller())
		return nil, err
	}

	response, err := c.Post(ctx, slashJoin(ListsURL, listID, MembersURL), nil, data)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"listID": listID,
			"error":  err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	var member *Member
	err = json.Unmarshal(response, &member)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"listID": listID,
			"error":  err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	member.Client = c

	return member, nil
}

type getMembers struct {
	Members    []*Member `json:"members"`
	ListID     string    `json:"list_id"`
	TotalItems int       `json:"total_items"`
}

func (c *Client) GetMembers(ctx context.Context, listID string, params ...Parameters) ([]*Member, error) {
	p := requestParameters(params)
	response, err := c.Get(ctx, slashJoin(ListsURL, listID, MembersURL), p)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"listID": listID,
			"error":  err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	var membersResponse *getMembers
	err = json.Unmarshal(response, &membersResponse)
	if err != nil {
		Log.Error(err.Error(), caller())
		return nil, err
	}

	// Add internal client
	members := []*Member{}
	for _, member := range membersResponse.Members {
		member.Client = c
		members = append(members, member)
	}

	return members, nil

}

func (c *Client) GetMember(ctx context.Context, id string, listID string) (*Member, error) {
	response, err := c.Get(ctx, slashJoin(ListsURL, listID, MembersURL, id), nil)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"listID":   listID,
			"memberID": id,
			"error":    err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	var member *Member
	err = json.Unmarshal(response, &member)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"listID":   listID,
			"memberID": id,
			"error":    err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	member.Client = c

	return member, nil
}

// Update Returns a new Member object with the updated values
func (m *Member) Update(ctx context.Context, data *UpdateMember) (*Member, error) {
	if m.Client == nil {
		return nil, ErrorNoClient
	}

	// If the member was previously deleted we need to use a PUT request,
	// otherwhise the API will tell us it's gone.
	response, err := m.Client.Put(ctx, slashJoin(ListsURL, m.ListID, MembersURL, m.ID), nil, data)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"listID":   m.ListID,
			"memberID": m.ID,
			"error":    err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	var member *Member
	err = json.Unmarshal(response, &member)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"listID":   m.ListID,
			"memberID": m.ID,
			"error":    err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	member.Client = m.Client

	return member, nil
}

func (m *Member) Delete(ctx context.Context) error {
	if m.Client == nil {
		return ErrorNoClient
	}
	err := m.Client.Delete(ctx, slashJoin(ListsURL, m.ListID, MembersURL, m.ID))
	if err != nil {
		Log.WithFields(logrus.Fields{
			"listID":   m.ListID,
			"memberID": m.ID,
			"error":    err.Error(),
		}).Error("response error", caller())
		return err
	}

	return nil
}
