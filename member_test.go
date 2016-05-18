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
	"net/http"
	"os"

	check "gopkg.in/check.v1"

	t "github.com/greatbeyond/mailchimp/testing"
)

var _ = check.Suite(&MemberSuite{})

type MemberSuite struct {
	client *Client
	server *t.MockServer
	ctx    context.Context
}

func (s *MemberSuite) SetUpSuite(c *check.C) {

}

func (s *MemberSuite) SetUpTest(c *check.C) {
	s.server = t.NewMockServer()
	s.server.SetChecker(c)

	s.client = NewClient()
	s.client.HTTPClient = s.server.HTTPClient

	s.ctx = NewContextWithToken(context.Background(), os.Getenv("MAILCHIMP_TEST_TOKEN"))
	// We need http to use the mock server
	s.ctx = NewContextWithURL(s.ctx, "http://us13.api.mailchimp.com/3.0/")
}

func (s *MemberSuite) TearDownTest(c *check.C) {}

func (s *MemberSuite) Test_NewMember(c *check.C) {
	mem := s.client.NewMember("abc23d")
	c.Assert(mem.ListID, check.Equals, "abc23d")
	c.Assert(mem.Client, check.Not(check.IsNil))
}

// --------------------------------------------------------------
// Create

func (s *MemberSuite) Test_CreateMember_Normal(c *check.C) {

	create := &CreateMember{
		EmailAddress: "urist.mcvankab+3@freddiesjokes.com",
		EmailType:    HTML,
		Status:       Subscribed,
		Interests: map[string]bool{
			"9143cf3bd1": false,
		},
		Vip: false,
		Location: &Location{
			Latitude:    55.30192,
			Longitude:   13.2928438,
			GmtOff:      1,
			DstOff:      0,
			CountryCode: "se",
		},
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "POST",
		Code:   200,
		Body:   `{"id":"852aaa9532cb36adfb5e9fef7a4206a9","email_address":"urist.mcvankab+3@freddiesjokes.com","unique_email_id":"fab20fa03d","email_type":"html","status":"subscribed","status_if_new":"","merge_fields":{"FNAME":"","LNAME":""},"interests":{"9143cf3bd1":false},"stats":{"avg_open_rate":0,"avg_click_rate":0},"ip_signup":"","timestamp_signup":"","ip_opt":"198.2.191.34","timestamp_opt":"2015-09-16 19:24:29","member_rating":2,"last_changed":"2015-09-16 19:24:29","language":"","vip":false,"email_client":"","location":{"latitude":0,"longitude":0,"gmtoff":0,"dstoff":0,"country_code":"","timezone":""},"list_id":"57afe96172"}`,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(body, check.Equals, `{"email_address":"urist.mcvankab+3@freddiesjokes.com","email_type":"html","status":"subscribed","interests":{"9143cf3bd1":false},"location":{"latitude":55.30192,"longitude":13.2928438,"gmtoff":1,"country_code":"se"}}`)
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/57afe96172/members")
		},
	})

	member, err := s.client.CreateMember(s.ctx, create, "57afe96172")
	c.Assert(err, check.IsNil)
	c.Assert(member, check.DeepEquals, &Member{
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

		Client: s.client,
	})
}

func (s *MemberSuite) Test_CreateMember_MissingStatus(c *check.C) {
	create := &CreateMember{
		EmailAddress: "urist.mcvankab+3@freddiesjokes.com",
		EmailType:    HTML,
		// Status: Subscribed,
	}
	_, err := s.client.CreateMember(s.ctx, create, "57afe96172")
	c.Assert(err, check.ErrorMatches, "missing field: Status")
}

func (s *MemberSuite) Test_CreateMember_MissingEmailAddress(c *check.C) {
	create := &CreateMember{
		// EmailAddress: "urist.mcvankab+3@freddiesjokes.com",
		EmailType: HTML,
		Status:    Subscribed,
	}
	_, err := s.client.CreateMember(s.ctx, create, "57afe96172")
	c.Assert(err, check.ErrorMatches, "missing field: EmailAddress")
}

func (s *MemberSuite) Test_CreateMember_MissingListID(c *check.C) {
	create := &CreateMember{}
	_, err := s.client.CreateMember(s.ctx, create, "")
	c.Assert(err, check.ErrorMatches, "missing argument: listID")
}

func (s *MemberSuite) Test_CreateMember_BadResponse(c *check.C) {
	create := &CreateMember{
		EmailAddress: "urist.mcvankab+3@freddiesjokes.com",
		EmailType:    HTML,
		Status:       Subscribed,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "POST",
		Code:   200,
		Body:   `{ bad json response`,
	})

	member, err := s.client.CreateMember(s.ctx, create, "57afe96172")
	c.Assert(err, check.ErrorMatches, "invalid character.*")
	c.Assert(member, check.IsNil)
}

func (s *MemberSuite) Test_CreateMember_UnknownResponse(c *check.C) {
	create := &CreateMember{
		EmailAddress: "urist.mcvankab+3@freddiesjokes.com",
		EmailType:    HTML,
		Status:       Subscribed,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "POST",
		Code:   111,
		Body:   `{}`,
	})

	member, err := s.client.CreateMember(s.ctx, create, "57afe96172")
	c.Assert(err, check.ErrorMatches, "Response error.*")
	c.Assert(member, check.IsNil)
}

// --------------------------------------------------------------
// GetMember

func (s *MemberSuite) Test_GetMembers_Normal(c *check.C) {

	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   `{"members":[{"id":"852aaa9532cb36adfb5e9fef7a4206a9","email_address":"urist.mcvankab+3@freddiesjokes.com","unique_email_id":"fab20fa03d","email_type":"html","status":"subscribed","status_if_new":"","merge_fields":{"FNAME":"","LNAME":""},"interests":{"9143cf3bd1":false,"3a2a927344":false,"f9c8f5f0ff":false},"stats":{"avg_open_rate":0,"avg_click_rate":0},"ip_signup":"","timestamp_signup":"","ip_opt":"198.2.191.34","timestamp_opt":"2015-09-16 19:24:29","member_rating":2,"last_changed":"2015-09-16 19:24:29","language":"","vip":false,"email_client":"","location":{"latitude":0,"longitude":0,"gmtoff":0,"dstoff":0,"country_code":"","timezone":""},"list_id":"57afe96172"}]}`,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/57afe96172/members")
		},
	})

	members, err := s.client.GetMembers(s.ctx, "57afe96172")
	c.Assert(err, check.IsNil)
	c.Assert(len(members), check.Equals, 1)
	c.Assert(members[0].Client, check.Not(check.IsNil))
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

		Client: s.client,
	})

}

func (s *MemberSuite) Test_GetMembers_BadResponse(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   `{ bad json response`,
	})

	Member, err := s.client.GetMembers(s.ctx, "57afe96172")
	c.Assert(err, check.ErrorMatches, "invalid character.*")
	c.Assert(Member, check.IsNil)
}

func (s *MemberSuite) Test_GetMembers_UnknownResponse(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   111,
		Body:   `{}`,
	})

	Member, err := s.client.GetMembers(s.ctx, "57afe96172")
	c.Assert(err, check.ErrorMatches, "Response error.*")
	c.Assert(Member, check.IsNil)
}

// --------------------------------------------------------------
// GetMember

func (s *MemberSuite) Test_GetMember_Normal(c *check.C) {

	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   `{"id":"852aaa9532cb36adfb5e9fef7a4206a9","email_address":"urist.mcvankab+3@freddiesjokes.com","unique_email_id":"fab20fa03d","email_type":"html","status":"subscribed","status_if_new":"","merge_fields":{"FNAME":"","LNAME":""},"interests":{"9143cf3bd1":false,"3a2a927344":false,"f9c8f5f0ff":false},"stats":{"avg_open_rate":0,"avg_click_rate":0},"ip_signup":"","timestamp_signup":"","ip_opt":"198.2.191.34","timestamp_opt":"2015-09-16 19:24:29","member_rating":2,"last_changed":"2015-09-16 19:24:29","language":"","vip":false,"email_client":"","location":{"latitude":0,"longitude":0,"gmtoff":0,"dstoff":0,"country_code":"","timezone":""},"list_id":"57afe96172"}`,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/57afe96172/members/852aaa9532cb36adfb5e9fef7a4206a9")
		},
	})

	member, err := s.client.GetMember(s.ctx, "852aaa9532cb36adfb5e9fef7a4206a9", "57afe96172")
	c.Assert(err, check.IsNil)
	c.Assert(member.Client, check.Not(check.IsNil))
	c.Assert(member, check.DeepEquals, &Member{
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

		Client: s.client,
	})

}

func (s *MemberSuite) Test_GetMember_BadResponse(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   `{ bad json response`,
	})

	member, err := s.client.GetMember(s.ctx, "0", "57afe96172")
	c.Assert(err, check.ErrorMatches, "invalid character.*")
	c.Assert(member, check.IsNil)
}

func (s *MemberSuite) Test_GetMember_UnknownResponse(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   111,
		Body:   `{}`,
	})

	member, err := s.client.GetMember(s.ctx, "0", "57afe96172")
	c.Assert(err, check.ErrorMatches, "Response error.*")
	c.Assert(member, check.IsNil)
}

// --------------------------------------------------------------
// Delete

func (s *MemberSuite) Test_Delete_Normal(c *check.C) {

	member := &Member{
		ID:     "852aaa9532cb36adfb5e9fef7a4206a9",
		ListID: "57afe96172",
		Client: s.client,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "DELETE",
		Code:   http.StatusNoContent,
		Body:   ``,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/57afe96172/members/852aaa9532cb36adfb5e9fef7a4206a9")
		},
	})

	err := member.Delete(s.ctx)
	c.Assert(err, check.IsNil)

}

func (s *MemberSuite) Test_Delete_NoClient(c *check.C) {
	member := &Member{
		ID:     "852aaa9532cb36adfb5e9fef7a4206a9",
		ListID: "57afe96172",
	}
	err := member.Delete(s.ctx)
	c.Assert(err, check.ErrorMatches, "no client assigned by parent")
}

func (s *MemberSuite) Test_Delete_UnknownResponse(c *check.C) {
	member := &Member{
		ID:     "852aaa9532cb36adfb5e9fef7a4206a9",
		ListID: "57afe96172",
		Client: s.client,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "DELETE",
		Code:   111,
		Body:   `{}`,
	})

	err := member.Delete(s.ctx)
	c.Assert(err, check.ErrorMatches, "Response error.*")

}

// --------------------------------------------------------------
// Update

func (s *MemberSuite) Test_Update_Normal(c *check.C) {

	member := &Member{
		ID:     "852aaa9532cb36adfb5e9fef7a4206a9",
		ListID: "57afe96172",
		Client: s.client,
	}

	update := &UpdateMember{
		Status: Unsubscribed,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "PUT",
		Code:   200,
		Body:   `{"id":"852aaa9532cb36adfb5e9fef7a4206a9","email_address":"urist.mcvankab+3@freddiesjokes.com","unique_email_id":"fab20fa03d","email_type":"html","status":"unsubscribed","status_if_new":"","merge_fields":{"FNAME":"","LNAME":""},"interests":{"9143cf3bd1":false,"3a2a927344":false,"f9c8f5f0ff":false},"stats":{"avg_open_rate":0,"avg_click_rate":0},"ip_signup":"","timestamp_signup":"","ip_opt":"198.2.191.34","timestamp_opt":"2015-09-16 19:24:29","member_rating":2,"last_changed":"2015-09-16 19:24:29","language":"","vip":false,"email_client":"","location":{"latitude":0,"longitude":0,"gmtoff":0,"dstoff":0,"country_code":"","timezone":""},"list_id":"57afe96172"}`,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(body, check.Equals, `{"status":"unsubscribed"}`)
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/57afe96172/members/852aaa9532cb36adfb5e9fef7a4206a9")
		},
	})

	upd, err := member.Update(s.ctx, update)
	c.Assert(err, check.IsNil)
	c.Assert(upd, check.DeepEquals, &Member{
		ID:            "852aaa9532cb36adfb5e9fef7a4206a9",
		EmailAddress:  "urist.mcvankab+3@freddiesjokes.com",
		UniqueEmailID: "fab20fa03d",
		EmailType:     HTML,
		Status:        Unsubscribed,
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

		Client: s.client,
	})
}

func (s *MemberSuite) Test_Update_Missing_Client(c *check.C) {
	member := &Member{
		ID:     "852aaa9532cb36adfb5e9fef7a4206a9",
		ListID: "57afe96172",
	}
	update := &UpdateMember{}
	_, err := member.Update(s.ctx, update)
	c.Assert(err, check.ErrorMatches, "no client assigned by parent")
}

func (s *MemberSuite) Test_Update_BadResponse(c *check.C) {
	updSegm := &UpdateMember{
		Status: Subscribed,
	}
	member := &Member{
		ID:     "852aaa9532cb36adfb5e9fef7a4206a9",
		ListID: "57afe96172",
		Client: s.client,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "PUT",
		Code:   200,
		Body:   `{ bad json response`,
	})

	upd, err := member.Update(s.ctx, updSegm)
	c.Assert(err, check.ErrorMatches, "invalid character.*")
	c.Assert(upd, check.IsNil)
}

func (s *MemberSuite) Test_Update_UnknownResponse(c *check.C) {
	updSegm := &UpdateMember{
		Status: Subscribed,
	}
	member := &Member{
		ID:     "852aaa9532cb36adfb5e9fef7a4206a9",
		ListID: "57afe96172",
		Client: s.client,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "PUT",
		Code:   111,
		Body:   `{}`,
	})

	upd, err := member.Update(s.ctx, updSegm)
	c.Assert(err, check.ErrorMatches, "Response error.*")
	c.Assert(upd, check.IsNil)
}
