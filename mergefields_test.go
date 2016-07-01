// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package mailchimp

import (
	"net/http"
	"strings"

	t "github.com/greatbeyond/mailchimp.v3/testing"

	check "gopkg.in/check.v1"
)

var _ = check.Suite(&MergeFieldSuite{})

type MergeFieldSuite struct {
	client *Client
	server *t.MockServer
}

func (s *MergeFieldSuite) SetUpSuite(c *check.C) {

}

func (s *MergeFieldSuite) SetUpTest(c *check.C) {
	s.server = t.NewMockServer()
	s.server.SetChecker(c)

	s.client = NewClient("b12824bd84759ef84abc67fd789e7570-us13")
	s.client.HTTPClient = s.server.HTTPClient
	s.client.APIURL = strings.Replace(s.client.APIURL, "https://", "http://", 1)
}

func (s *MergeFieldSuite) TearDownTest(c *check.C) {}

func (s *MergeFieldSuite) Test_NewMergeField(c *check.C) {
	seg := s.client.NewMergeField()
	c.Assert(seg.client, check.Not(check.IsNil))
}

// --------------------------------------------------------------
// Create

func (s *MergeFieldSuite) Test_CreateMergeField_Normal(c *check.C) {

	create := &CreateMergeField{
		MergeID:      3,
		Tag:          "MMERGE3",
		Name:         "FAVORITEJOKE",
		Type:         MergeFieldTypeText,
		Required:     false,
		Public:       false,
		DefaultValue: "",
		DisplayOrder: 6,
		Options: map[string]interface{}{
			"size": 25.0,
		},
		HelpText: "",
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "POST",
		Code:   200,
		Body:   `{"merge_id":3,"tag":"MMERGE3","name":"FAVORITEJOKE","type":"text","required":false,"default_value":"","public":false,"display_order":6,"options":{"size":25},"help_text":"","list_id":"57afe96172"}`,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(body, check.Equals, `{"merge_id":3,"tag":"MMERGE3","name":"FAVORITEJOKE","type":"text","display_order":6,"options":{"size":25}}`)
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/57afe96172/merge-fields")
		},
	})

	mergefield, err := s.client.CreateMergeField(create, "57afe96172")
	c.Assert(err, check.IsNil)
	c.Assert(mergefield, check.DeepEquals, &MergeField{
		MergeID:      3,
		Tag:          "MMERGE3",
		Name:         "FAVORITEJOKE",
		Type:         MergeFieldTypeText,
		Required:     false,
		Public:       false,
		DefaultValue: "",
		DisplayOrder: 6,
		Options: map[string]interface{}{
			"size": 25.0,
		},
		HelpText: "",
		ListID:   "57afe96172",
		client:   s.client,
	})
}

func (s *MergeFieldSuite) Test_CreateMergeField_MissingName(c *check.C) {
	create := &CreateMergeField{
		MergeID: 3,
		Tag:     "MMERGE3",
		// Name:         "FAVORITEJOKE",
		Type: MergeFieldTypeText,
	}
	_, err := s.client.CreateMergeField(create, "57afe96172")
	c.Assert(err, check.ErrorMatches, "missing field: Name")
}

func (s *MergeFieldSuite) Test_CreateMergeField_MissingType(c *check.C) {
	create := &CreateMergeField{
		MergeID: 3,
		Tag:     "MMERGE3",
		Name:    "FAVORITEJOKE",
		// Type: MergeFieldTypeText,
	}
	_, err := s.client.CreateMergeField(create, "57afe96172")
	c.Assert(err, check.ErrorMatches, "missing field: Type")
}

func (s *MergeFieldSuite) Test_CreateMergeField_LongTag(c *check.C) {
	create := &CreateMergeField{
		MergeID: 3,
		Tag:     "MMERGE3OVERSIZED",
		Name:    "FAVORITEJOKE",
		Type:    MergeFieldTypeText,
	}
	_, err := s.client.CreateMergeField(create, "57afe96172")
	c.Assert(err, check.ErrorMatches, "tag length over limit \\(10\\)")
}

func (s *MergeFieldSuite) Test_CreateMergeField_MissingListID(c *check.C) {
	create := &CreateMergeField{}
	_, err := s.client.CreateMergeField(create, "")
	c.Assert(err, check.ErrorMatches, "missing field: listID")
}

func (s *MergeFieldSuite) Test_CreateMergeField_BadResponse(c *check.C) {
	create := &CreateMergeField{
		Name: "MergeField in list",
		Type: MergeFieldTypeText,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "POST",
		Code:   200,
		Body:   `{ bad json response`,
	})

	mergefield, err := s.client.CreateMergeField(create, "57afe96172")
	c.Assert(err, check.ErrorMatches, "invalid character.*")
	c.Assert(mergefield, check.IsNil)
}

func (s *MergeFieldSuite) Test_CreateMergeField_UnknownResponse(c *check.C) {
	create := &CreateMergeField{
		Name: "MergeField in list",
		Type: MergeFieldTypeText,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "POST",
		Code:   111,
		Body:   `{}`,
	})

	mergefield, err := s.client.CreateMergeField(create, "57afe96172")
	c.Assert(err, check.ErrorMatches, "Response error.*")
	c.Assert(mergefield, check.IsNil)
}

// --------------------------------------------------------------
// GetMergeFields

func (s *MergeFieldSuite) Test_GetMergeFields_Normal(c *check.C) {

	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   `{"merge_fields":[{"merge_id":1,"tag":"FNAME","name":"FirstName","type":"text","required":false,"default_value":"","public":true,"display_order":2,"options":{"size":25},"help_text":"","list_id":"57afe96172"}],"list_id":"57afe96172","total_items":1}`,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/57afe96172/merge-fields")
		},
	})

	mergefield, err := s.client.GetMergeFields("57afe96172")
	c.Assert(err, check.IsNil)
	c.Assert(len(mergefield), check.Equals, 1)
	c.Assert(mergefield[0].client, check.Not(check.IsNil))
	c.Assert(mergefield[0], check.DeepEquals, &MergeField{
		MergeID:      1,
		Tag:          "FNAME",
		Name:         "FirstName",
		Type:         MergeFieldTypeText,
		Required:     false,
		Public:       true,
		DefaultValue: "",
		DisplayOrder: 2,
		Options: map[string]interface{}{
			"size": 25.0,
		},
		HelpText: "",
		ListID:   "57afe96172",
		client:   s.client,
	})

}

func (s *MergeFieldSuite) Test_GetMergeFields_BadResponse(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   `{ bad json response`,
	})

	MergeField, err := s.client.GetMergeFields("57afe96172")
	c.Assert(err, check.ErrorMatches, "invalid character.*")
	c.Assert(MergeField, check.IsNil)
}

func (s *MergeFieldSuite) Test_GetMergeFields_UnknownResponse(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   111,
		Body:   `{}`,
	})

	mergefield, err := s.client.GetMergeFields("57afe96172")
	c.Assert(err, check.ErrorMatches, "Response error.*")
	c.Assert(mergefield, check.IsNil)
}

// --------------------------------------------------------------
// GetMergeField

func (s *MergeFieldSuite) Test_GetMergeField_Normal(c *check.C) {

	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   `{"merge_id":3,"tag":"MMERGE3","name":"FAVORITEJOKE","type":"text","required":false,"default_value":"","public":false,"display_order":6,"options":{"size":25},"help_text":"","list_id":"57afe96172"}`,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/57afe96172/merge-fields/3")
		},
	})

	mergefield, err := s.client.GetMergeField(3, "57afe96172")
	c.Assert(err, check.IsNil)
	c.Assert(mergefield.client, check.Not(check.IsNil))
	c.Assert(mergefield, check.DeepEquals, &MergeField{
		MergeID:      3,
		Tag:          "MMERGE3",
		Name:         "FAVORITEJOKE",
		Type:         MergeFieldTypeText,
		Required:     false,
		Public:       false,
		DefaultValue: "",
		DisplayOrder: 6,
		Options: map[string]interface{}{
			"size": 25.0,
		},
		HelpText: "",
		ListID:   "57afe96172",
		client:   s.client,
	})

}

func (s *MergeFieldSuite) Test_GetMergeField_BadResponse(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   `{ bad json response`,
	})

	mergefield, err := s.client.GetMergeField(0, "57afe96172")
	c.Assert(err, check.ErrorMatches, "invalid character.*")
	c.Assert(mergefield, check.IsNil)
}

func (s *MergeFieldSuite) Test_GetMergeField_UnknownResponse(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   111,
		Body:   `{}`,
	})

	mergefield, err := s.client.GetMergeField(0, "57afe96172")
	c.Assert(err, check.ErrorMatches, "Response error.*")
	c.Assert(mergefield, check.IsNil)
}

// --------------------------------------------------------------
// Delete

func (s *MergeFieldSuite) Test_Delete_Normal(c *check.C) {

	mergefield := &MergeField{
		MergeID: 3,
		ListID:  "57afe96172",
		client:  s.client,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "DELETE",
		Code:   http.StatusNoContent,
		Body:   ``,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/57afe96172/merge-fields/3")
		},
	})

	err := mergefield.Delete()
	c.Assert(err, check.IsNil)

}

func (s *MergeFieldSuite) Test_Delete_NoClient(c *check.C) {
	mergefield := &MergeField{
		MergeID: 49377,
		ListID:  "57afe96172",
	}
	err := mergefield.Delete()
	c.Assert(err, check.ErrorMatches, "no client assigned by parent")
}

func (s *MergeFieldSuite) Test_Delete_UnknownResponse(c *check.C) {
	mergefield := &MergeField{
		MergeID: 49377,
		ListID:  "57afe96172",
		client:  s.client,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "DELETE",
		Code:   111,
		Body:   `{}`,
	})

	err := mergefield.Delete()
	c.Assert(err, check.ErrorMatches, "Response error.*")

}

// --------------------------------------------------------------
// Update

func (s *MergeFieldSuite) Test_Update_Normal(c *check.C) {

	mergefield := &MergeField{
		MergeID: 49377,
		ListID:  "57afe96172",
		client:  s.client,
	}

	update := &UpdateMergeField{
		Tag:  "MMERGE3",
		Name: "FAVORITEJOKE",
		Type: MergeFieldTypeText,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "PUT",
		Code:   200,
		Body:   `{"merge_id":3,"tag":"MMERGE3","name":"FAVORITEJOKE","type":"text","required":false,"default_value":"","public":false,"display_order":6,"options":{"size":25},"help_text":"","list_id":"57afe96172"}`,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(body, check.Equals, `{"tag":"MMERGE3","name":"FAVORITEJOKE","type":"text"}`)
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/57afe96172/merge-fields/49377")
		},
	})

	upd, err := mergefield.Update(update)
	c.Assert(err, check.IsNil)
	c.Assert(upd, check.DeepEquals, &MergeField{
		MergeID:      3,
		Tag:          "MMERGE3",
		Name:         "FAVORITEJOKE",
		Type:         MergeFieldTypeText,
		Required:     false,
		Public:       false,
		DefaultValue: "",
		DisplayOrder: 6,
		Options: map[string]interface{}{
			"size": 25.0,
		},
		HelpText: "",
		ListID:   "57afe96172",
		client:   s.client,
	})
}

func (s *MergeFieldSuite) Test_Update_Missing_Client(c *check.C) {
	mergefield := &MergeField{
		MergeID: 49377,
		ListID:  "57afe96172",
	}
	update := &UpdateMergeField{}
	_, err := mergefield.Update(update)
	c.Assert(err, check.ErrorMatches, "no client assigned by parent")
}

func (s *MergeFieldSuite) Test_Update_BadResponse(c *check.C) {
	updSegm := &UpdateMergeField{
		Name: "MergeField in list",
	}
	mergefield := &MergeField{
		MergeID: 49377,
		ListID:  "57afe96172",
		client:  s.client,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "PUT",
		Code:   200,
		Body:   `{ bad json response`,
	})

	upd, err := mergefield.Update(updSegm)
	c.Assert(err, check.ErrorMatches, "invalid character.*")
	c.Assert(upd, check.IsNil)
}

func (s *MergeFieldSuite) Test_Update_UnknownResponse(c *check.C) {
	updSegm := &UpdateMergeField{
		Name: "MergeField in list",
	}
	mergefield := &MergeField{
		MergeID: 49377,
		ListID:  "57afe96172",
		client:  s.client,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "PUT",
		Code:   111,
		Body:   `{}`,
	})

	upd, err := mergefield.Update(updSegm)
	c.Assert(err, check.ErrorMatches, "Response error.*")
	c.Assert(upd, check.IsNil)
}
