// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package mailchimp

import (
	"encoding/json"
	"fmt"
	"strconv"
	"unicode/utf8"

	"github.com/Sirupsen/logrus"
)

const MergeFieldsURL = "/merge-fields"

// MergeField are custom var fields on lists
type MergeField struct {
	// An unchanging id for the merge field.
	MergeID int `json:"merge_id,omitempty"`

	// The tag used in MailChimp campaigns and for the /fields endpoint.
	Tag string `json:"tag,omitempty"`

	// The name of the merge field. Max 10 chars.
	Name string `json:"name"`

	// The type for the merge field.
	Type MergeFieldType `json:"type"`

	// The boolean value if the merge field is required.
	Required bool `json:"required,omitempty"`

	// The default value for the merge field if null.
	DefaultValue string `json:"default_value,omitempty"`

	// Whether the merge field is displayed on the signup form.
	Public bool `json:"public,omitempty"`

	// The order that the merge field displays on the list signup form.
	DisplayOrder int `json:"display_order,omitempty"`

	// Extra options for some merge field types.
	// In an address field, the default country code if none supplied.
	//   default_country   int
	// In a phone field, the phone number type: US or International.
	//   phone_format      string
	// In a date or birthday field, the format of the date.
	//   date_format       string
	// In a radio or dropdown non-group field, the available options for fields to pick from.
	//   choices           []string
	// In a text field, the default length of the text field.
	//   size              int
	Options map[string]interface{} `json:"options,omitempty"`

	// Extra text to help the subscriber fill out the form.
	HelpText string `json:"help_text,omitempty"`

	// A string that identifies this merge field collections’ list.
	ListID string `json:"list_id,omitempty"`

	// Internal
	client MailchimpClient
}

// SetClient fulfills ClientType
func (m *MergeField) SetClient(c MailchimpClient) { m.client = c }

type MergeFieldType string

const (
	MergeFieldTypeText       MergeFieldType = "text"
	MergeFieldTypeNumber     MergeFieldType = "number"
	MergeFieldTypeAddress    MergeFieldType = "address"
	MergeFieldTypePhone      MergeFieldType = "phone"
	MergeFieldTypeEmail      MergeFieldType = "email"
	MergeFieldTypeDate       MergeFieldType = "date"
	MergeFieldTypeURL        MergeFieldType = "url"
	MergeFieldTypeImageurl   MergeFieldType = "imageurl"
	MergeFieldTypeRadio      MergeFieldType = "radio"
	MergeFieldTypeDropdown   MergeFieldType = "dropdown"
	MergeFieldTypeCheckboxes MergeFieldType = "checkboxes"
	MergeFieldTypeBirthday   MergeFieldType = "birthday"
	MergeFieldTypeZip        MergeFieldType = "zip"
)

// CreateMergeField is a alias for MergeField, the keys are the same.
type CreateMergeField MergeField

// UpdateMergeField and CreateMergeField are the same but with slighlty
// different requiered fields (checked in function)
type UpdateMergeField CreateMergeField

// NewMergeField returns a empty field object
// id is optional, with it you can do a bit of rudimentary chaining.
// Example:
//	c.NewMergeField(23).Update(params)
func (c *Client) NewMergeField(id ...int) *MergeField {
	s := &MergeField{
		client: c,
	}
	if len(id) > 0 {
		s.MergeID = id[0]
	}
	return s
}

// CreateMergeField Creates a field object and inserts it
func (c *Client) CreateMergeField(data *CreateMergeField, listID string) (*MergeField, error) {

	if listID == "" {
		return nil, fmt.Errorf("missing argument: listID")
	}

	if err := hasFields(*data, "Name", "Type"); err != nil {
		Log.Info(err.Error(), caller())
		return nil, err
	}

	if utf8.RuneCountInString(data.Tag) > 10 {
		return nil, fmt.Errorf("tag length over limit (10)")
	}

	response, err := c.Post(slashJoin(ListsURL, listID, MergeFieldsURL), nil, data)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"listID": listID,
			"error":  err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	var field *MergeField
	err = json.Unmarshal(response, &field)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"listID": listID,
			"error":  err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	field.client = c

	return field, nil
}

type getMergeField struct {
	MergeField []*MergeField `json:"merge_fields"`
	ListID     string        `json:"list_id"`
	TotalItems int           `json:"total_items"`
}

// GetMergeFields fetches all merge fields
func (c *Client) GetMergeFields(listID string, params ...Parameters) ([]*MergeField, error) {

	p := requestParameters(params)
	response, err := c.Get(slashJoin(ListsURL, listID, MergeFieldsURL), p)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"listID": listID,
			"error":  err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	var mergefieldsResponse *getMergeField
	err = json.Unmarshal(response, &mergefieldsResponse)
	if err != nil {
		return nil, err
	}

	// Add internal client
	mergefields := []*MergeField{}
	for _, mergefield := range mergefieldsResponse.MergeField {
		mergefield.client = c
		mergefields = append(mergefields, mergefield)
	}

	return mergefields, nil

}

// GetMergeField retrives a single merge field
func (c *Client) GetMergeField(id int, listID string) (*MergeField, error) {
	response, err := c.Get(slashJoin(ListsURL, listID, MergeFieldsURL, strconv.Itoa(id)), nil)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"listID":   listID,
			"merge_id": id,
			"error":    err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	var field *MergeField
	err = json.Unmarshal(response, &field)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"listID":   listID,
			"merge_id": id,
			"error":    err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	field.client = c

	return field, nil
}

//Delete remvoes the merge field
func (m *MergeField) Delete() error {

	if m.client == nil {
		return ErrorNoClient
	}
	err := m.client.Delete(slashJoin(ListsURL, m.ListID, MergeFieldsURL, strconv.Itoa(m.MergeID)))
	if err != nil {
		Log.WithFields(logrus.Fields{
			"listID":   m.ListID,
			"merge_id": m.MergeID,
			"error":    err.Error(),
		}).Error("response error", caller())
		return err
	}

	return nil
}

// Update returns a existing MergeField object with the updated values
func (m *MergeField) Update(data *UpdateMergeField) (*MergeField, error) {

	if m.client == nil {
		return nil, ErrorNoClient
	}

	// If the field was previously deleted we need to use a PUT request,
	// otherwhise the API will tell us it's gone.
	response, err := m.client.Put(slashJoin(ListsURL, m.ListID, MergeFieldsURL, strconv.Itoa(m.MergeID)), nil, data)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"listID":   m.ListID,
			"merge_id": m.MergeID,
			"error":    err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	var field *MergeField
	err = json.Unmarshal(response, &field)
	if err != nil {
		Log.WithFields(logrus.Fields{
			"listID":   m.ListID,
			"merge_id": m.MergeID,
			"error":    err.Error(),
		}).Error("response error", caller())
		return nil, err
	}

	field.client = m.client

	return field, nil
}
