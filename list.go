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
	"time"
)

// ListVisibility denotes list visibility when creating and updating lists
type ListVisibility string

const (
	ListVisibilityPublic  ListVisibility = "pub"
	ListVisibilityPrivate ListVisibility = "prv"
)

const ListsURL = "/lists"

// List defines a Mailchimp list as received from server
type List struct {

	// A string that uniquely identifies this list.
	ID string `json:"id,omitempty"`

	// The name of the list.
	Name string `json:"name,omitempty"`

	// Contact information displayed in campaign footers to comply with international spam laws.
	Contact Contact `json:"contact,omitempty"`

	// The permission reminder for the list.
	PermissionReminder string `json:"permission_reminder,omitempty"`

	// Whether campaigns for this list use the Archive Bar in archives by default.
	UseArchiveBar bool `json:"use_archive_bar,omitempty"`

	// Default values for campaigns created for this list.
	CampaignDefaults CampaignDefaults `json:"campaign_defaults,omitempty"`

	// The email address to send subscribe notifications to.
	NotifyOnSubscribe string `json:"notify_on_subscribe,omitempty"`

	// The email address to send unsubscribe notifications to.
	NotifyOnUnsubscribe string `json:"notify_on_unsubscribe,omitempty"`

	// The date and time that this list was created.
	DateCreated string `json:"date_created,omitempty"`

	// An auto-generated activity score for the list (0-5).
	ListRating int `json:"list_rating,omitempty"`

	// Whether the list supports multiple formats for emails. When set to true,
	// subscribers can choose whether they want to receive HTML or plain-text
	// emails. When set to false, subscribers will receive HTML emails, with a plain-text alternative backup.
	EmailTypeOption bool `json:"email_type_option,omitempty"`

	// Our EepURL shortened version of this list’s subscribe form.
	SubscribeURLShort string `json:"subscribe_url_short,omitempty"`

	// The full version of this list’s subscribe form (host will vary).
	SubscribeURLLong string `json:"subscribe_url_long,omitempty"`

	// The list’s Email Beamer address.
	BeamerAddress string `json:"beamer_address,omitempty"`

	// Whether this list is public or private.
	Visibility ListVisibility `json:"visibility,omitempty"`

	// Any list-specific modules installed for this list.
	Modules []interface{} `json:"modules,omitempty"`

	// Stats for the list. Many of these are cached for at least five minutes.
	Stats ListStats `json:"stats,omitempty"`

	// Internal
	Client MailchimpClient `json:"-"`
}

// SetClient fulfills ClientType
func (l *List) SetClient(c MailchimpClient) { l.Client = c }

// Contact for list information
type Contact struct {
	Company  string `json:"company"    bson:"company"`
	Address1 string `json:"address1"   bson:"address1"`
	Address2 string `json:"address2"   bson:"address2"`
	City     string `json:"city"       bson:"city"`
	State    string `json:"state"      bson:"state"`
	Zip      string `json:"zip"        bson:"zip"`
	Country  string `json:"country"    bson:"country"`
	Phone    string `json:"phone"      bson:"phone"`
}

// CampaignDefaults is added to lists
type CampaignDefaults struct {
	FromName  string `json:"from_name"  bson:"from_name"`
	FromEmail string `json:"from_email" bson:"from_email"`
	Subject   string `json:"subject"    bson:"subject"`
	Language  string `json:"language"   bson:"language"`
}

// ListStats are stats recevied from server
type ListStats struct {
	MemberCount               int     `json:"member_count"`
	UnsubscribeCount          int     `json:"unsubscribe_count"`
	CleanedCount              int     `json:"cleaned_count"`
	MemberCountSinceSend      int     `json:"member_count_since_send"`
	UnsubscribeCountSinceSend int     `json:"unsubscribe_count_since_send"`
	CleanedCountSinceSend     int     `json:"cleaned_count_since_send"`
	CampaignCount             int     `json:"campaign_count"`
	CampaignLastSent          string  `json:"campaign_last_sent"`
	MergeFieldCount           int     `json:"merge_field_count"`
	AvgSubRate                float64 `json:"avg_sub_rate"`
	AvgUnsubRate              float64 `json:"avg_unsub_rate"`
	TargetSubRate             float64 `json:"target_sub_rate"`
	OpenRate                  float64 `json:"open_rate"`
	ClickRate                 float64 `json:"click_rate"`
	LastSubDate               string  `json:"last_sub_date"`
	LastUnsubDate             string  `json:"last_unsub_date"`
}

// CreateList defines fields neccessary to create a new list.
// Some fields are optional:
// 		UseArchiveBar, NotifyOnSubscribe, NotifyOnUnsubscribe, Visibility
//
type CreateList struct {
	// The name of the list.
	Name string `                        json:"name,omitempty"`

	// Contact information displayed in campaign footers to comply with international spam laws.
	Contact *Contact `                   json:"contact,omitempty"`

	// The permission reminder for the list.
	PermissionReminder string `          json:"permission_reminder,omitempty"`

	// Whether campaigns for this list use the Archive Bar in archives by default.
	UseArchiveBar bool `                 json:"use_archive_bar,omitempty"`

	// Default values for campaigns created for this list.
	CampaignDefaults *CampaignDefaults ` json:"campaign_defaults,omitempty"`

	// The email address to send subscribe notifications to.
	NotifyOnSubscribe string `           json:"notify_on_subscribe,omitempty"`

	// The email address to send unsubscribe notifications to.
	NotifyOnUnsubscribe string `         json:"notify_on_unsubscribe,omitempty"`

	// Whether the list supports multiple formats for emails. When set
	// to true, subscribers can choose whether they want to receive HTML
	// or plain-text emails. When set to false, subscribers will receive
	// HTML emails, with a plain-text alternative backup.
	EmailTypeOption bool `              json:"email_type_option,omitempty"`

	// Whether this list is public or private. Possible Values:
	// pub, prv
	Visibility ListVisibility `         json:"visibility,omitempty"`
}

// UpdateList and CreateList are the same but with slighlty
// different requiered fields (checked in function)
type UpdateList CreateList

// NewList returns a empty member object
// id is optional, with it you can do a bit of rudimentary chaining.
// Example:
//	c.NewList(23).Update(params)
func (c *Client) NewList(id ...string) *List {
	s := &List{
		Client: c,
	}
	if len(id) > 0 {
		s.ID = id[0]
	}
	return s
}

// CreateList Creates a member object and inserts it
func (c *Client) CreateList(ctx context.Context, data *CreateList) (*List, error) {
	fields := []string{"Name", "Contact", "PermissionReminder", "CampaignDefaults"}
	if err := hasFields(*data, fields...); err != nil {
		return nil, err
	}

	response, err := c.Post(ctx, ListsURL, nil, data)
	if err != nil {
		Log.Error(err.Error(), caller())
		return nil, err
	}

	var list *List
	err = json.Unmarshal(response, &list)
	if err != nil {
		Log.Error(err.Error(), caller())
		return nil, err
	}

	list.Client = c

	return list, nil
}

type getListsResponse struct {
	Lists      []*List `json:"lists"`
	TotalItems int     `json:"total_items"`
}

func (c *Client) GetLists(ctx context.Context) ([]*List, error) {
	response, err := c.Get(ctx, ListsURL, nil)
	if err != nil {
		Log.Error(err.Error(), caller())
		return nil, err
	}

	var listsResponse getListsResponse
	err = json.Unmarshal(response, &listsResponse)
	if err != nil {
		Log.Error(err.Error(), caller())
		return nil, err
	}

	// Add internal client
	lists := []*List{}
	for _, list := range listsResponse.Lists {
		list.Client = c
		lists = append(lists, list)
	}

	return lists, nil
}

// GetList returns a single list by id
func (c *Client) GetList(ctx context.Context, id string) (*List, error) {
	response, err := c.Get(ctx, slashJoin(ListsURL, id), nil)
	if err != nil {
		Log.Error(err.Error(), caller())
		return nil, err
	}

	var list *List
	err = json.Unmarshal(response, &list)
	if err != nil {
		return nil, err
	}

	if list == nil {
		return nil, fmt.Errorf("unable to unmarshal response")
	}

	list.Client = c

	return list, nil
}

//Update returns a new List object with the updated values
func (l *List) Update(ctx context.Context, data *UpdateList) (*List, error) {
	if l.Client == nil {
		return nil, ErrorNoClient
	}

	response, err := l.Client.Patch(ctx, slashJoin(ListsURL, l.ID), nil, data)
	if err != nil {
		Log.Error(err.Error(), caller())
		return nil, err
	}

	var list *List
	err = json.Unmarshal(response, &list)
	if err != nil {
		Log.Error(err.Error(), caller())
		return nil, err
	}

	list.Client = l.Client

	return list, nil
}

// Delete removes the list
func (l *List) Delete(ctx context.Context) error {
	if l.Client == nil {
		return ErrorNoClient
	}
	return l.Client.Delete(ctx, slashJoin(ListsURL, l.ID))
}

// TimeCreated converts DateCreated to a time.Time object
func (l *List) TimeCreated() time.Time {
	d, _ := StringToTime(l.DateCreated)
	return d
}
