// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package mailchimp

import (
	"encoding/json"

	"github.com/Sirupsen/logrus"
	"github.com/antonholmquist/jason"
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

type Member struct {
	ID              string       `json:"id"`
	EmailAddress    string       `json:"email_address"`
	UniqueEmailID   string       `json:"unique_email_id"`
	EmailType       string       `json:"email_type"`
	Status          string       `json:"status"`
	MergeFields     interface{}  `json:"merge_fields"`
	Stats           *MemberStats `json:"stats"`
	IPSignup        string       `json:"ip_signup"`
	TimestampSignup string       `json:"timestamp_signup"`
	IPOpt           string       `json:"ip_opt"`
	TimestampOpt    string       `json:"timestamp_opt"`
	MemberRating    int          `json:"member_rating"`
	LastChanged     string       `json:"last_changed"`
	Language        string       `json:"language"`
	Vip             bool         `json:"vip"`
	EmailClient     string       `json:"email_client"`
	Location        *Location    `json:"location"`
	ListID          string       `json:"list_id"`

	// Internal
	client MailchimpClient
}

type CreateMember struct {
	// Email address for a subscriber.
	EmailAddress string `           json:"email_address,omitempty"`

	// Type of email this member asked to get (‘html’ or ‘text’).
	EmailType MailType `            json:"email_type,omitempty"`

	// Subscriber’s current status. Possible Values:
	// subscribed, unsubscribed, cleaned, pending
	Status MemberStatus `           json:"status,omitempty"`

	// An individual merge var and value for a member.
	MergeFields interface{} `       json:"merge_fields,omitempty"`

	// The key of this object’s properties is the ID of the interest in question.
	Interests interface{} `         json:"interests,omitempty"`

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
	AvgOpenRate  float64 `json:"avg_open_rate"`
	AvgClickRate float64 `json:"avg_click_rate"`
}

type Location struct {
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Gmtoff      float64 `json:"gmtoff"`
	Dstoff      float64 `json:"dstoff"`
	CountryCode string  `json:"country_code"`
	Timezone    string  `json:"timezone"`
}

// NewMember returns a empty member object
func (c *Client) NewMember() *Member {
	return &Member{
		client: c,
	}
}

// CreateMember Creates a member object and inserts it
func (c *Client) CreateMember(data *CreateMember, listID string) (*Member, error) {

	if err := missingField(listID, "listID"); err != nil {
		c.Log().Debug(err.Error, caller())
		return nil, err
	}

	if err := missingField(data.Status, "status"); err != nil {
		c.Log().Debug(err.Error, caller())
		return nil, err
	}

	response, err := c.Post(slashJoin(ListsURL, listID, MembersURL), nil, data)
	if err != nil {
		c.Log().WithFields(logrus.Fields{
			"listID": listID,
			"error":  err.Error(),
		}).Debug("response error", caller())
		return nil, err
	}

	var member *Member
	err = json.Unmarshal(response, &member)
	if err != nil {
		c.Log().WithFields(logrus.Fields{
			"listID": listID,
			"error":  err.Error(),
		}).Debug("response error", caller())
		return nil, err
	}

	member.client = c

	return member, nil
}

func (c *Client) GetMembers(listID string, params ...Parameters) ([]*Member, error) {

	p := requestParameters(params)
	response, err := c.Get(slashJoin(ListsURL, listID, MembersURL), p)
	if err != nil {
		c.Log().WithFields(logrus.Fields{
			"listID": listID,
			"error":  err.Error(),
		}).Debug("response error", caller())
		return nil, err
	}

	v, err := jason.NewObjectFromBytes(response)
	if err != nil {
		c.Log().WithFields(logrus.Fields{
			"listID": listID,
			"error":  err.Error(),
		}).Debug("response error", caller())
		return nil, err
	}

	_members, err := v.GetValue("members")
	if err != nil {
		c.Log().WithFields(logrus.Fields{
			"listID": listID,
			"error":  err.Error(),
		}).Debug("response error", caller())
		return nil, err
	}

	b, err := _members.Marshal()
	if err != nil {
		c.Log().WithFields(logrus.Fields{
			"listID": listID,
			"error":  err.Error(),
		}).Debug("response error", caller())
		return nil, err
	}

	var members []*Member
	err = json.Unmarshal(b, &members)
	if err != nil {
		c.Log().WithFields(logrus.Fields{
			"listID": listID,
			"error":  err.Error(),
		}).Debug("response error", caller())
		return nil, err
	}

	for _, l := range members {
		l.client = c
	}

	return members, nil

}

func (c *Client) GetMember(id string, listID string) (*Member, error) {
	response, err := c.Get(slashJoin(ListsURL, listID, MembersURL, id), nil)
	if err != nil {
		c.Log().WithFields(logrus.Fields{
			"listID":   listID,
			"memberID": id,
			"error":    err.Error(),
		}).Debug("response error", caller())
		return nil, err
	}

	var member *Member
	err = json.Unmarshal(response, &member)
	if err != nil {
		c.Log().WithFields(logrus.Fields{
			"listID":   listID,
			"memberID": id,
			"error":    err.Error(),
		}).Debug("response error", caller())
		return nil, err
	}

	member.client = c

	return member, nil
}

func (m *Member) Delete() error {

	if m.client == nil {
		return ErrorNoClient
	}
	err := m.client.Delete(slashJoin(ListsURL, m.ListID, MembersURL, m.ID))
	if err != nil {
		m.client.Log().WithFields(logrus.Fields{
			"listID":   m.ListID,
			"memberID": m.ID,
			"error":    err.Error(),
		}).Debug("response error", caller())
		return err
	}

	return nil
}

// Update Returns a new Member object with the updated values
func (m *Member) Update(data *UpdateMember) (*Member, error) {

	if m.client == nil {
		return nil, ErrorNoClient
	}

	// If the member was previously deleted we need to use a PUT request,
	// otherwhise the API will tell us it's gone.
	response, err := m.client.Put(slashJoin(ListsURL, m.ListID, MembersURL, m.ID), nil, data)
	if err != nil {
		m.client.Log().WithFields(logrus.Fields{
			"listID":   m.ListID,
			"memberID": m.ID,
			"error":    err.Error(),
		}).Debug("response error", caller())
		return nil, err
	}

	var member *Member
	err = json.Unmarshal(response, &member)
	if err != nil {
		m.client.Log().WithFields(logrus.Fields{
			"listID":   m.ListID,
			"memberID": m.ID,
			"error":    err.Error(),
		}).Debug("response error", caller())
		return nil, err
	}

	member.client = m.client

	return member, nil
}
