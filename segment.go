// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package mailchimp

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Sirupsen/logrus"
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
	client MailchimpClient
}

// CreateSegment sends a request to create a segment
type CreateSegment struct {
	// The name of the segment.
	Name string `json:"name,omitempty"`

	// An array of emails to be used for a static segment.
	// Any emails provided that are not present on the list will be ignored.
	// Passing an empty array will create a static segment without any subscribers.
	// This field cannot be provided with the options field.
	StaticSegment []string `json:"static_segment,omitempty"`

	// The conditions of the segment. Static and fuzzy segments don’t have conditions.
	// See API reference for list of possible matching options.
	Options map[string]interface{} `json:"options,omitempty"`
}

// UpdateSegment is an alias since Create and Update share the same keys
type UpdateSegment CreateSegment

// NewSegment returns a empty segment object
func (c *Client) NewSegment() *Segment {
	return &Segment{
		client: c,
	}
}

// CreateSegment Creates a segment object and inserts it
func (c *Client) CreateSegment(data *CreateSegment, listID string) (*Segment, error) {

	if listID == "" {
		return nil, fmt.Errorf("missing argument: listID")
	}

	if err := missingField(*data, "Name"); err != nil {
		Log.Info(err.Error(), caller())
		return nil, err
	}

	response, err := c.Post(slashJoin(ListsURL, listID, SegmentsURL), nil, data)
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

	segment.client = c

	return segment, nil
}

type getSegments struct {
	Segments   []*Segment `json:"segments"`
	ListID     string     `json:"list_id"`
	TotalItems int        `json:"total_items"`
}

func (c *Client) GetSegments(listID string, params ...Parameters) ([]*Segment, error) {

	p := requestParameters(params)
	response, err := c.Get(slashJoin(ListsURL, listID, SegmentsURL), p)
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
		segment.client = c
		segments = append(segments, segment)
	}

	return segments, nil

}

func (c *Client) GetSegment(id string, listID string) (*Segment, error) {
	response, err := c.Get(slashJoin(ListsURL, listID, SegmentsURL, id), nil)
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

	segment.client = c

	return segment, nil
}

func (m *Segment) Delete() error {

	if m.client == nil {
		return ErrorNoClient
	}
	err := m.client.Delete(slashJoin(ListsURL, m.ListID, SegmentsURL, strconv.Itoa(m.ID)))
	if err != nil {
		Log.WithFields(logrus.Fields{
			"list_id":    m.ListID,
			"segment_id": m.ID,
			"error":      err.Error(),
		}).Error("response error", caller())
		return err
	}

	return nil
}

// Updatereturns a new Segment object with the updated values
func (m *Segment) Update(data *UpdateSegment) (*Segment, error) {

	if m.client == nil {
		return nil, ErrorNoClient
	}

	// If the segment was previously deleted we need to use a PUT request,
	// otherwhise the API will tell us it's gone.
	response, err := m.client.Put(slashJoin(ListsURL, m.ListID, SegmentsURL, strconv.Itoa(m.ID)), nil, data)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"list_id":    m.ListID,
			"segment_id": m.ID,
			"error":      err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	var segment *Segment
	err = json.Unmarshal(response, &segment)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"list_id":    m.ListID,
			"segment_id": m.ID,
			"error":      err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	segment.client = m.client

	return segment, nil
}
