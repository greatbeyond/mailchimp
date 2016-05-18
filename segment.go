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
	"strconv"

	"github.com/sirupsen/logrus"
)

const SegmentsURL = "/segments"

// Segment manages segments for a specific MailChimp list. A segment is
// a section of your list that includes only those subscribers
// who share specific common field information.
// http://developer.mailchimp.com/documentation/mailchimp/reference/lists/segments/#
type Segment struct {
	// The unique id for the segment.
	ID int `json:"id,omitempty"`

	// The name of the segment.
	Name string `json:"name,omitempty"`

	// The number of active subscribers currently included in the segment.
	MemberCount int `json:"member_count,omitempty"`

	// The type of segment.
	Type string `json:"type,omitempty"`

	// The date and time the segment was created.
	CreatedAt string `json:"created_at,omitempty"`

	// The date and time the segment was last updated.
	UpdatedAt string `json:"updated_at,omitempty"`

	// The conditions of the segment. Static and fuzzy segments don’t have conditions.
	Options map[string]interface{} `json:"options,omitempty"`

	// The list id.
	ListID string `json:"list_id,omitempty"`

	// Internal
	Client MailchimpClient `json:"-"`
}

// SetClient fulfills ClientType
func (s *Segment) SetClient(c MailchimpClient) { s.Client = c }

// CreateSegment sends a request to create a segment
type CreateSegment struct {
	// The name of the segment.
	Name string `json:"name,omitempty"`

	// An array of emails to be used for a static segment.
	// Any emails provided that are not present on the list will be ignored.
	// Passing an empty array will create a static segment without any subscribers.
	// This field cannot be provided with the options field.
	//
	// Due to how json-omitempty works, this needs to be a pointer in order
	// to make it work with the exclusive-or nature of StaticSegment and Options.
	StaticSegment *[]string `json:"static_segment,omitempty"`

	// The conditions of the segment. Static and fuzzy segments don’t have conditions.
	// See API reference for list of possible matching options.
	Options map[string]interface{} `json:"options,omitempty"`
}

// UpdateSegment is an alias since Create and Update share the same keys
type UpdateSegment CreateSegment

// NewSegment returns a empty segment object
// id is optional, with it you can do a bit of rudimentary chaining.
// Example:
//	c.NewSegment("abc23d", 23).Update(params)
func (c *Client) NewSegment(listID string, id ...int) *Segment {
	s := &Segment{
		Client: c,
	}
	if len(id) > 0 {
		s.ID = id[0]
	}

	s.ListID = listID

	return s
}

// CreateSegment Creates a segment object and inserts it
func (c *Client) CreateSegment(ctx context.Context, data *CreateSegment, listID string) (*Segment, error) {

	if listID == "" {
		return nil, fmt.Errorf("missing argument: listID")
	}

	if err := hasFields(*data, "Name"); err != nil {
		Log.Info(err.Error(), caller())
		return nil, err
	}

	response, err := c.Post(ctx, slashJoin(ListsURL, listID, SegmentsURL), nil, data)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"list_id": listID,
			"error":   err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	var segment *Segment
	err = json.Unmarshal(response, &segment)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"list_id": listID,
			"error":   err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	segment.Client = c

	return segment, nil
}

type getSegments struct {
	Segments   []*Segment `json:"segments"`
	ListID     string     `json:"list_id"`
	TotalItems int        `json:"total_items"`
}

func (c *Client) GetSegments(ctx context.Context, listID string, params ...Parameters) ([]*Segment, error) {

	p := requestParameters(params)
	response, err := c.Get(ctx, slashJoin(ListsURL, listID, SegmentsURL), p)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"list_id": listID,
			"error":   err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	var segmentsResponse *getSegments
	err = json.Unmarshal(response, &segmentsResponse)
	if err != nil {
		return nil, err
	}

	// Add internal client
	segments := []*Segment{}
	for _, segment := range segmentsResponse.Segments {
		segment.Client = c
		segments = append(segments, segment)
	}

	return segments, nil

}

func (c *Client) GetSegment(ctx context.Context, id string, listID string) (*Segment, error) {
	response, err := c.Get(ctx, slashJoin(ListsURL, listID, SegmentsURL, id), nil)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"list_id":    listID,
			"segment_id": id,
			"error":      err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	var segment *Segment
	err = json.Unmarshal(response, &segment)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"list_id":    listID,
			"segment_id": id,
			"error":      err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	segment.Client = c

	return segment, nil
}

// GetMembers returns all members in this Segment.
func (s *Segment) GetMembers(ctx context.Context, params ...Parameters) ([]*Member, error) {
	p := requestParameters(params)
	response, err := s.Client.Get(ctx, slashJoin(ListsURL, s.ListID, SegmentsURL, strconv.Itoa(s.ID), MembersURL), p)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"list_id":    s.ListID,
			"segment_id": s.ID,
			"error":      err.Error(),
		})
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
		member.Client = s.Client
		members = append(members, member)
	}

	return members, nil
}

func (s *Segment) Delete(ctx context.Context) error {

	if s.Client == nil {
		return ErrorNoClient
	}

	if err := hasFields(*s, "ID", "ListID"); err != nil {
		Log.Info(err.Error(), caller())
		return err
	}

	err := s.Client.Delete(ctx, slashJoin(ListsURL, s.ListID, SegmentsURL, strconv.Itoa(s.ID)))
	if err != nil {
		Log.WithFields(logrus.Fields{
			"list_id":    s.ListID,
			"segment_id": s.ID,
			"error":      err.Error(),
		}).Error("response error", caller())
		return err
	}

	return nil
}

// Update returns a new Segment object with the updated values
func (s *Segment) Update(ctx context.Context, data *UpdateSegment) (*Segment, error) {

	if s.Client == nil {
		return nil, ErrorNoClient
	}

	if err := hasFields(*s, "ID", "ListID"); err != nil {
		Log.Info(err.Error(), caller())
		return nil, err
	}

	// If the segment was previously deleted we need to use a PATCH request,
	// otherwhise the API will tell us it's gone.
	response, err := s.Client.Patch(ctx, slashJoin(ListsURL, s.ListID, SegmentsURL, strconv.Itoa(s.ID)), nil, data)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"list_id":    s.ListID,
			"segment_id": s.ID,
			"error":      err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	var segment *Segment
	err = json.Unmarshal(response, &segment)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"list_id":    s.ListID,
			"segment_id": s.ID,
			"error":      err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	segment.Client = s.Client

	return segment, nil
}
