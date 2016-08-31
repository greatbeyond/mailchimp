// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package mailchimp

import (
	"net/http"
	"strings"

	t "github.com/greatbeyond/mailchimp/testing"

	check "gopkg.in/check.v1"
)

var _ = check.Suite(&SegmentSuite{})

type SegmentSuite struct {
	client *Client
	server *t.MockServer
}

func (s *SegmentSuite) SetUpSuite(c *check.C) {

}

func (s *SegmentSuite) SetUpTest(c *check.C) {
	s.server = t.NewMockServer()
	s.server.SetChecker(c)

	s.client = NewClient("b12824bd84759ef84abc67fd789e7570-us13")
	s.client.HTTPClient = s.server.HTTPClient
	s.client.APIURL = strings.Replace(s.client.APIURL, "https://", "http://", 1)
}

func (s *SegmentSuite) TearDownTest(c *check.C) {}

func (s *SegmentSuite) Test_NewSegment(c *check.C) {
	seg := s.client.NewSegment("abc23d")
	c.Assert(seg.ListID, check.Equals, "abc23d")
	c.Assert(seg.client, check.Not(check.IsNil))
}

// --------------------------------------------------------------
// Create

func (s *SegmentSuite) Test_CreateSegment_Normal(c *check.C) {

	create := &CreateSegment{
		Name: "Segment in list",
		StaticSegment: &[]string{
			"hsims0@ihg.com",
			"acox1@alibaba.com",
			"jlopez2@deliciousdays.com",
		},
		Options: map[string]interface{}{
			"test": "value",
		},
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "POST",
		Code:   200,
		Body:   `{"id":49377,"name":"Freddie'sMostPopularJokes","member_count":9,"type":"static","created_at":"2015-09-16 21:14:46","updated_at":"2015-09-16 21:14:47","options":{"match":"any","conditions":[]},"list_id":"57afe96172"}`,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(body, check.Equals, `{"name":"Segment in list","static_segment":["hsims0@ihg.com","acox1@alibaba.com","jlopez2@deliciousdays.com"],"options":{"test":"value"}}`)
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/57afe96172/segments")
		},
	})

	segment, err := s.client.CreateSegment(create, "57afe96172")
	c.Assert(err, check.IsNil)
	c.Assert(segment, check.DeepEquals, &Segment{
		ID:          49377,
		Name:        "Freddie'sMostPopularJokes",
		MemberCount: 9,
		Type:        "static",
		CreatedAt:   "2015-09-16 21:14:46",
		UpdatedAt:   "2015-09-16 21:14:47",
		Options: map[string]interface{}{
			"match":      "any",
			"conditions": []interface{}{},
		},
		ListID: "57afe96172",
		client: s.client,
	})
}

func (s *SegmentSuite) Test_CreateSegment_MissingName(c *check.C) {
	create := &CreateSegment{
		// Name: "Segment in list",
		StaticSegment: &[]string{
			"hsims0@ihg.com",
			"acox1@alibaba.com",
			"jlopez2@deliciousdays.com",
		},
		Options: map[string]interface{}{
			"test": "value",
		},
	}
	_, err := s.client.CreateSegment(create, "57afe96172")
	c.Assert(err, check.ErrorMatches, "missing field: Name")
}

func (s *SegmentSuite) Test_CreateSegment_MissingListID(c *check.C) {
	create := &CreateSegment{}
	_, err := s.client.CreateSegment(create, "")
	c.Assert(err, check.ErrorMatches, "missing argument: listID")
}

func (s *SegmentSuite) Test_CreateSegment_BadResponse(c *check.C) {
	create := &CreateSegment{
		Name: "Segment in list",
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "POST",
		Code:   200,
		Body:   `{ bad json response`,
	})

	segment, err := s.client.CreateSegment(create, "57afe96172")
	c.Assert(err, check.ErrorMatches, "invalid character.*")
	c.Assert(segment, check.IsNil)
}

func (s *SegmentSuite) Test_CreateSegment_UnknownResponse(c *check.C) {
	create := &CreateSegment{
		Name: "Segment in list",
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "POST",
		Code:   111,
		Body:   `{}`,
	})

	segment, err := s.client.CreateSegment(create, "57afe96172")
	c.Assert(err, check.ErrorMatches, "Response error.*")
	c.Assert(segment, check.IsNil)
}

// --------------------------------------------------------------
// GetSegments

func (s *SegmentSuite) Test_GetSegments_Normal(c *check.C) {

	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   `{"segments":[{"id":49377,"name":"Freddie'sMostPopularJokes","member_count":9,"type":"static","created_at":"2015-09-1621:14:46","updated_at":"2015-09-1621:14:47","options":{"match":"any","conditions":[]},"list_id":"57afe96172"}],"list_id":"57afe96172","total_items":1}`,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/57afe96172/segments")
		},
	})

	segments, err := s.client.GetSegments("57afe96172")
	c.Assert(err, check.IsNil)
	c.Assert(len(segments), check.Equals, 1)
	c.Assert(segments[0].client, check.Not(check.IsNil))
	c.Assert(segments[0], check.DeepEquals, &Segment{
		ID:          49377,
		Name:        "Freddie'sMostPopularJokes",
		MemberCount: 9,
		Type:        "static",
		CreatedAt:   "2015-09-1621:14:46",
		UpdatedAt:   "2015-09-1621:14:47",
		Options: map[string]interface{}{
			"match":      "any",
			"conditions": []interface{}{},
		},
		ListID: "57afe96172",
		client: s.client,
	})

}

func (s *SegmentSuite) Test_GetSegments_BadResponse(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   `{ bad json response`,
	})

	segments, err := s.client.GetSegments("57afe96172")
	c.Assert(err, check.ErrorMatches, "invalid character.*")
	c.Assert(segments, check.IsNil)
}

func (s *SegmentSuite) Test_GetSegments_UnknownResponse(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   111,
		Body:   `{}`,
	})

	segments, err := s.client.GetSegments("57afe96172")
	c.Assert(err, check.ErrorMatches, "Response error.*")
	c.Assert(segments, check.IsNil)
}

// --------------------------------------------------------------
// GetSegment

func (s *SegmentSuite) Test_GetSegment_Normal(c *check.C) {

	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   `{"id":49377,"name":"Freddie'sMostPopularJokes","member_count":9,"type":"static","created_at":"2015-09-16 21:14:46","updated_at":"2015-09-16 21:14:47","options":{"match":"any","conditions":[]},"list_id":"57afe96172"}`,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/57afe96172/segments/49377")
		},
	})

	segment, err := s.client.GetSegment("49377", "57afe96172")
	c.Assert(err, check.IsNil)
	c.Assert(segment.client, check.Not(check.IsNil))
	c.Assert(segment, check.DeepEquals, &Segment{
		ID:          49377,
		Name:        "Freddie'sMostPopularJokes",
		MemberCount: 9,
		Type:        "static",
		CreatedAt:   "2015-09-16 21:14:46",
		UpdatedAt:   "2015-09-16 21:14:47",
		Options: map[string]interface{}{
			"match":      "any",
			"conditions": []interface{}{},
		},
		ListID: "57afe96172",
		client: s.client,
	})

}

func (s *SegmentSuite) Test_GetSegment_BadResponse(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   `{ bad json response`,
	})

	segment, err := s.client.GetSegment("0", "57afe96172")
	c.Assert(err, check.ErrorMatches, "invalid character.*")
	c.Assert(segment, check.IsNil)
}

func (s *SegmentSuite) Test_GetSegment_UnknownResponse(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   111,
		Body:   `{}`,
	})

	segment, err := s.client.GetSegment("0", "57afe96172")
	c.Assert(err, check.ErrorMatches, "Response error.*")
	c.Assert(segment, check.IsNil)
}

// --------------------------------------------------------------
// GetMembers
func (s *SegmentSuite) Test_GetMembers_Normal(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   `{"members":[{"id":"852aaa9532cb36adfb5e9fef7a4206a9","email_address":"urist.mcvankab+3@freddiesjokes.com","unique_email_id":"fab20fa03d","email_type":"html","status":"subscribed","status_if_new":"","merge_fields":{"FNAME":"","LNAME":""},"interests":{"9143cf3bd1":false,"3a2a927344":false,"f9c8f5f0ff":false},"stats":{"avg_open_rate":0,"avg_click_rate":0},"ip_signup":"","timestamp_signup":"","ip_opt":"198.2.191.34","timestamp_opt":"2015-09-16 19:24:29","member_rating":2,"last_changed":"2015-09-16 19:24:29","language":"","vip":false,"email_client":"","location":{"latitude":0,"longitude":0,"gmtoff":0,"dstoff":0,"country_code":"","timezone":""},"list_id":"57afe96172"}]}`,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/57afe96172/segments/23/members")
		},
	})

	segment := s.client.NewSegment("57afe96172", 23)
	members, err := segment.GetMembers()

	c.Assert(err, check.IsNil)
	c.Assert(len(members), check.Equals, 1)
	c.Assert(members[0].client, check.Not(check.IsNil))
	c.Assert(members[0], check.DeepEquals, &Member{
		ID:            "852aaa9532cb36adfb5e9fef7a4206a9",
		EmailAddress:  "urist.mcvankab+3@freddiesjokes.com",
		UniqueEmailID: "fab20fa03d",
		EmailType:     HTML,
		Status:        Subscribed,
		MergeFields: map[string]interface{}{
			"FNAME": "",
			"LNAME": "",
		},
		Interests: map[string]bool{
			"9143cf3bd1": false,
			"3a2a927344": false,
			"f9c8f5f0ff": false,
		},
		Stats: MemberStats{
			AvgOpenRate:  0,
			AvgClickRate: 0,
		},
		IPSignup:        "",
		TimestampSignup: "",
		IPOpt:           "198.2.191.34",
		TimestampOpt:    "2015-09-16 19:24:29",
		MemberRating:    2,
		LastChanged:     "2015-09-16 19:24:29",
		Language:        "",
		Vip:             false,
		EmailClient:     "",
		Location: Location{
			Latitude:    0,
			Longitude:   0,
			GmtOff:      0,
			DstOff:      0,
			CountryCode: "",
			Timezone:    "",
		},
		ListID: "57afe96172",

		client: s.client,
	})
}

// --------------------------------------------------------------
// Delete

func (s *SegmentSuite) Test_Delete_Normal(c *check.C) {

	segment := &Segment{
		ID:     49377,
		ListID: "57afe96172",
		client: s.client,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "DELETE",
		Code:   http.StatusNoContent,
		Body:   ``,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/57afe96172/segments/49377")
		},
	})

	err := segment.Delete()
	c.Assert(err, check.IsNil)

}

func (s *SegmentSuite) Test_Delete_NoClient(c *check.C) {
	segment := &Segment{
		ID:     49377,
		ListID: "57afe96172",
	}
	err := segment.Delete()
	c.Assert(err, check.ErrorMatches, "no client assigned by parent")
}

func (s *SegmentSuite) Test_Delete_UnknownResponse(c *check.C) {
	segment := &Segment{
		ID:     49377,
		ListID: "57afe96172",
		client: s.client,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "DELETE",
		Code:   111,
		Body:   `{}`,
	})

	err := segment.Delete()
	c.Assert(err, check.ErrorMatches, "Response error.*")

}

// --------------------------------------------------------------
// Update

func (s *SegmentSuite) Test_Update_Normal(c *check.C) {

	segment := &Segment{
		ID:     49377,
		ListID: "57afe96172",
		client: s.client,
	}

	update := &UpdateSegment{
		Name: "Segment in list",
		StaticSegment: &[]string{
			"hsims0@ihg.com",
			"acox1@alibaba.com",
			"jlopez2@deliciousdays.com",
		},
		Options: map[string]interface{}{
			"test": "value",
		},
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "PATCH",
		Code:   200,
		Body:   `{"id":49377,"name":"Freddie'sMostPopularJokes","member_count":9,"type":"static","created_at":"2015-09-16 21:14:46","updated_at":"2015-09-16 21:14:47","options":{"match":"any","conditions":[]},"list_id":"57afe96172"}`,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(body, check.Equals, `{"name":"Segment in list","static_segment":["hsims0@ihg.com","acox1@alibaba.com","jlopez2@deliciousdays.com"],"options":{"test":"value"}}`)
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/57afe96172/segments/49377")
		},
	})

	upd, err := segment.Update(update)
	c.Assert(err, check.IsNil)
	c.Assert(upd, check.DeepEquals, &Segment{
		ID:          49377,
		Name:        "Freddie'sMostPopularJokes",
		MemberCount: 9,
		Type:        "static",
		CreatedAt:   "2015-09-16 21:14:46",
		UpdatedAt:   "2015-09-16 21:14:47",
		Options: map[string]interface{}{
			"match":      "any",
			"conditions": []interface{}{},
		},
		ListID: "57afe96172",
		client: s.client,
	})
}

func (s *SegmentSuite) Test_Update_Missing_Client(c *check.C) {
	segment := &Segment{
		ID:     49377,
		ListID: "57afe96172",
	}
	update := &UpdateSegment{}
	_, err := segment.Update(update)
	c.Assert(err, check.ErrorMatches, "no client assigned by parent")
}

func (s *SegmentSuite) Test_Update_BadResponse(c *check.C) {
	updSegm := &UpdateSegment{
		Name: "Segment in list",
	}
	segment := &Segment{
		ID:     49377,
		ListID: "57afe96172",
		client: s.client,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "PATCH",
		Code:   200,
		Body:   `{ bad json response`,
	})

	upd, err := segment.Update(updSegm)
	c.Assert(err, check.ErrorMatches, "invalid character.*")
	c.Assert(upd, check.IsNil)
}

func (s *SegmentSuite) Test_Update_UnknownResponse(c *check.C) {
	updSegm := &UpdateSegment{
		Name: "Segment in list",
	}
	segment := &Segment{
		ID:     49377,
		ListID: "57afe96172",
		client: s.client,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "PATCH",
		Code:   111,
		Body:   `{}`,
	})

	upd, err := segment.Update(updSegm)
	c.Assert(err, check.ErrorMatches, "Response error.*")
	c.Assert(upd, check.IsNil)
}
