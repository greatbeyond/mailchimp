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

var _ = check.Suite(&ListSuite{})

type ListSuite struct {
	client *Client
	server *t.MockServer
}

func (s *ListSuite) SetUpSuite(c *check.C) {

}

func (s *ListSuite) SetUpTest(c *check.C) {
	s.server = t.NewMockServer()
	s.server.SetChecker(c)

	s.client = NewClient("b12824bd84759ef84abc67fd789e7570-us13")
	s.client.HTTPClient = s.server.HTTPClient
	s.client.APIURL = strings.Replace(s.client.APIURL, "https://", "http://", 1)
}

func (s *ListSuite) TearDownTest(c *check.C) {}

func (s *ListSuite) Test_NewList(c *check.C) {
	mem := s.client.NewList()
	c.Assert(mem.client, check.Not(check.IsNil))
}

// --------------------------------------------------------------
// Create

func (s *ListSuite) Test_CreateList_Normal(c *check.C) {

	create := &CreateList{
		Name: "Freddies Favorite Hats",
		Contact: &Contact{
			Company:  "MailChimp",
			Address1: "675 Ponce De Leon Ave NE",
			Address2: "Suite 5000",
			City:     "Atlanta",
			State:    "GA",
			Zip:      "30308",
			Country:  "US",
			Phone:    "",
		},
		PermissionReminder: "you signed up",
		CampaignDefaults: &CampaignDefaults{
			FromName:  "Freddie",
			FromEmail: "freddie@freddiehats.com",
			Subject:   "",
			Language:  "en",
		},
		EmailTypeOption: true,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "POST",
		Code:   200,
		Body:   `{"id":"1510500e0b","name":"Freddie's Favorite Hats","contact":{"company":"MailChimp","address1":"675 Ponce De Leon Ave NE","address2":"Suite 5000","city":"Atlanta","state":"GA","zip":"30308","country":"US","phone":""},"permission_reminder":"","use_archive_bar":true,"campaign_defaults":{"from_name":"Freddie","from_email":"freddie@freddiehats.com","subject":"","language":"en"},"notify_on_subscribe":"","notify_on_unsubscribe":"","date_created":"2015-09-16T14:55:51+00:00","list_rating":0,"email_type_option":true,"subscribe_url_short":"http://eepurl.com/xxxx","subscribe_url_long":"http://freddieshats.usX.list-manage.com/subscribe?u=8d3a3db4d97663a9074efcc16&id=1510500e0b","beamer_address":"usX-xxxx-xxxx@inbound.mailchimp.com","visibility":"pub","modules":[],"stats":{"member_count":0,"unsubscribe_count":0,"cleaned_count":0,"member_count_since_send":0,"unsubscribe_count_since_send":0,"cleaned_count_since_send":0,"campaign_count":0,"campaign_last_sent":"","merge_field_count":2,"avg_sub_rate":0,"avg_unsub_rate":0,"target_sub_rate":0,"open_rate":0,"click_rate":0,"last_sub_date":"","last_unsub_date":""}}`,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(body, check.Equals, `{"name":"Freddies Favorite Hats","contact":{"company":"MailChimp","address1":"675 Ponce De Leon Ave NE","address2":"Suite 5000","city":"Atlanta","state":"GA","zip":"30308","country":"US","phone":""},"permission_reminder":"you signed up","campaign_defaults":{"from_name":"Freddie","from_email":"freddie@freddiehats.com","subject":"","language":"en"},"email_type_option":true}`)
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists")
		},
	})

	list, err := s.client.CreateList(create)
	c.Assert(err, check.IsNil)
	c.Assert(list, check.DeepEquals, &List{
		ID:   "1510500e0b",
		Name: "Freddie's Favorite Hats",
		Contact: Contact{
			Company:  "MailChimp",
			Address1: "675 Ponce De Leon Ave NE",
			Address2: "Suite 5000",
			City:     "Atlanta",
			State:    "GA",
			Zip:      "30308",
			Country:  "US",
			Phone:    "",
		},
		PermissionReminder: "",
		UseArchiveBar:      true,
		CampaignDefaults: CampaignDefaults{
			FromName:  "Freddie",
			FromEmail: "freddie@freddiehats.com",
			Subject:   "",
			Language:  "en",
		},
		NotifyOnSubscribe:   "",
		NotifyOnUnsubscribe: "",
		DateCreated:         "2015-09-16T14:55:51+00:00",
		ListRating:          0,
		EmailTypeOption:     true,
		SubscribeURLShort:   "http://eepurl.com/xxxx",
		SubscribeURLLong:    "http://freddieshats.usX.list-manage.com/subscribe?u=8d3a3db4d97663a9074efcc16&id=1510500e0b",
		BeamerAddress:       "usX-xxxx-xxxx@inbound.mailchimp.com",
		Visibility:          "pub",
		Modules:             []interface{}{},
		Stats: ListStats{
			MemberCount:               0,
			UnsubscribeCount:          0,
			CleanedCount:              0,
			MemberCountSinceSend:      0,
			UnsubscribeCountSinceSend: 0,
			CleanedCountSinceSend:     0,
			CampaignCount:             0,
			CampaignLastSent:          "",
			MergeFieldCount:           2,
			AvgSubRate:                0,
			AvgUnsubRate:              0,
			TargetSubRate:             0,
			OpenRate:                  0,
			ClickRate:                 0,
			LastSubDate:               "",
			LastUnsubDate:             "",
		},

		client: s.client,
	})
}

func (s *ListSuite) Test_CreateList_MissingName(c *check.C) {
	create := &CreateList{
		// Name: "Freddies Favorite Hats",
		Contact: &Contact{
			Company:  "MailChimp",
			Address1: "675 Ponce De Leon Ave NE",
			Address2: "Suite 5000",
			City:     "Atlanta",
			State:    "GA",
			Zip:      "30308",
			Country:  "US",
			Phone:    "",
		},
		PermissionReminder: "you signed up",
		CampaignDefaults: &CampaignDefaults{
			FromName:  "Freddie",
			FromEmail: "freddie@freddiehats.com",
			Subject:   "",
			Language:  "en",
		},
		EmailTypeOption: true,
	}
	_, err := s.client.CreateList(create)
	c.Assert(err, check.ErrorMatches, "missing field: Name")
}

func (s *ListSuite) Test_CreateList_MissingContact(c *check.C) {
	create := &CreateList{
		Name: "Freddies Favorite Hats",
		// Contact: &Contact{
		// 	Company:  "MailChimp",
		// 	Address1: "675 Ponce De Leon Ave NE",
		// 	Address2: "Suite 5000",
		// 	City:     "Atlanta",
		// 	State:    "GA",
		// 	Zip:      "30308",
		// 	Country:  "US",
		// 	Phone:    "",
		// },
		PermissionReminder: "you signed up",
		CampaignDefaults: &CampaignDefaults{
			FromName:  "Freddie",
			FromEmail: "freddie@freddiehats.com",
			Subject:   "",
			Language:  "en",
		},
		EmailTypeOption: true,
	}
	_, err := s.client.CreateList(create)
	c.Assert(err, check.ErrorMatches, "missing field: Contact")
}

func (s *ListSuite) Test_CreateList_MissingPermissionReminder(c *check.C) {
	create := &CreateList{
		Name: "Freddies Favorite Hats",
		Contact: &Contact{
			Company:  "MailChimp",
			Address1: "675 Ponce De Leon Ave NE",
			Address2: "Suite 5000",
			City:     "Atlanta",
			State:    "GA",
			Zip:      "30308",
			Country:  "US",
			Phone:    "",
		},
		// PermissionReminder: "you signed up",
		CampaignDefaults: &CampaignDefaults{
			FromName:  "Freddie",
			FromEmail: "freddie@freddiehats.com",
			Subject:   "",
			Language:  "en",
		},
		EmailTypeOption: true,
	}
	_, err := s.client.CreateList(create)
	c.Assert(err, check.ErrorMatches, "missing field: PermissionReminder")
}

func (s *ListSuite) Test_CreateList_MissingCampaignDefaults(c *check.C) {
	create := &CreateList{
		Name: "Freddies Favorite Hats",
		Contact: &Contact{
			Company:  "MailChimp",
			Address1: "675 Ponce De Leon Ave NE",
			Address2: "Suite 5000",
			City:     "Atlanta",
			State:    "GA",
			Zip:      "30308",
			Country:  "US",
			Phone:    "",
		},
		PermissionReminder: "you signed up",
		// CampaignDefaults: &CampaignDefaults{
		// 	FromName:  "Freddie",
		// 	FromEmail: "freddie@freddiehats.com",
		// 	Subject:   "",
		// 	Language:  "en",
		// },
		EmailTypeOption: true,
	}
	_, err := s.client.CreateList(create)
	c.Assert(err, check.ErrorMatches, "missing field: CampaignDefaults")
}

func (s *ListSuite) Test_CreateList_BadResponse(c *check.C) {
	create := &CreateList{
		Name: "Freddies Favorite Hats",
		Contact: &Contact{
			Company:  "MailChimp",
			Address1: "675 Ponce De Leon Ave NE",
			Address2: "Suite 5000",
			City:     "Atlanta",
			State:    "GA",
			Zip:      "30308",
			Country:  "US",
			Phone:    "",
		},
		PermissionReminder: "you signed up",
		CampaignDefaults: &CampaignDefaults{
			FromName:  "Freddie",
			FromEmail: "freddie@freddiehats.com",
			Subject:   "",
			Language:  "en",
		},
		EmailTypeOption: true,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "POST",
		Code:   200,
		Body:   `{ bad json response`,
	})

	list, err := s.client.CreateList(create)
	c.Assert(err, check.ErrorMatches, "invalid character.*")
	c.Assert(list, check.IsNil)
}

func (s *ListSuite) Test_CreateList_UnknownResponse(c *check.C) {
	create := &CreateList{
		Name: "Freddies Favorite Hats",
		Contact: &Contact{
			Company:  "MailChimp",
			Address1: "675 Ponce De Leon Ave NE",
			Address2: "Suite 5000",
			City:     "Atlanta",
			State:    "GA",
			Zip:      "30308",
			Country:  "US",
			Phone:    "",
		},
		PermissionReminder: "you signed up",
		CampaignDefaults: &CampaignDefaults{
			FromName:  "Freddie",
			FromEmail: "freddie@freddiehats.com",
			Subject:   "",
			Language:  "en",
		},
		EmailTypeOption: true,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "POST",
		Code:   111,
		Body:   `{}`,
	})

	list, err := s.client.CreateList(create)
	c.Assert(err, check.ErrorMatches, "Response error.*")
	c.Assert(list, check.IsNil)
}

// --------------------------------------------------------------
// GetList

func (s *ListSuite) Test_GetLists_Normal(c *check.C) {

	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   `{"lists":[{"id":"1510500e0b","name":"Freddie's Favorite Hats","contact":{"company":"MailChimp","address1":"675 Ponce De Leon Ave NE","address2":"Suite 5000","city":"Atlanta","state":"GA","zip":"30308","country":"US","phone":""},"permission_reminder":"","use_archive_bar":true,"campaign_defaults":{"from_name":"Freddie","from_email":"freddie@freddiehats.com","subject":"","language":"en"},"notify_on_subscribe":"","notify_on_unsubscribe":"","date_created":"2015-09-16T14:55:51+00:00","list_rating":0,"email_type_option":true,"subscribe_url_short":"http://eepurl.com/xxxx","subscribe_url_long":"http://freddieshats.usX.list-manage.com/subscribe?u=8d3a3db4d97663a9074efcc16&id=1510500e0b","beamer_address":"usX-xxxx-xxxx@inbound.mailchimp.com","visibility":"pub","modules":[],"stats":{"member_count":0,"unsubscribe_count":0,"cleaned_count":0,"member_count_since_send":0,"unsubscribe_count_since_send":0,"cleaned_count_since_send":0,"campaign_count":0,"campaign_last_sent":"","merge_field_count":2,"avg_sub_rate":0,"avg_unsub_rate":0,"target_sub_rate":0,"open_rate":0,"click_rate":0,"last_sub_date":"","last_unsub_date":""}}],"total_items":1}`,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists")
		},
	})

	lists, err := s.client.GetLists()
	c.Assert(err, check.IsNil)
	c.Assert(len(lists), check.Equals, 1)
	c.Assert(lists[0].client, check.Not(check.IsNil))
	c.Assert(lists[0], check.DeepEquals, &List{
		ID:   "1510500e0b",
		Name: "Freddie's Favorite Hats",
		Contact: Contact{
			Company:  "MailChimp",
			Address1: "675 Ponce De Leon Ave NE",
			Address2: "Suite 5000",
			City:     "Atlanta",
			State:    "GA",
			Zip:      "30308",
			Country:  "US",
			Phone:    "",
		},
		PermissionReminder: "",
		UseArchiveBar:      true,
		CampaignDefaults: CampaignDefaults{
			FromName:  "Freddie",
			FromEmail: "freddie@freddiehats.com",
			Subject:   "",
			Language:  "en",
		},
		NotifyOnSubscribe:   "",
		NotifyOnUnsubscribe: "",
		DateCreated:         "2015-09-16T14:55:51+00:00",
		ListRating:          0,
		EmailTypeOption:     true,
		SubscribeURLShort:   "http://eepurl.com/xxxx",
		SubscribeURLLong:    "http://freddieshats.usX.list-manage.com/subscribe?u=8d3a3db4d97663a9074efcc16&id=1510500e0b",
		BeamerAddress:       "usX-xxxx-xxxx@inbound.mailchimp.com",
		Visibility:          "pub",
		Modules:             []interface{}{},
		Stats: ListStats{
			MemberCount:               0,
			UnsubscribeCount:          0,
			CleanedCount:              0,
			MemberCountSinceSend:      0,
			UnsubscribeCountSinceSend: 0,
			CleanedCountSinceSend:     0,
			CampaignCount:             0,
			CampaignLastSent:          "",
			MergeFieldCount:           2,
			AvgSubRate:                0,
			AvgUnsubRate:              0,
			TargetSubRate:             0,
			OpenRate:                  0,
			ClickRate:                 0,
			LastSubDate:               "",
			LastUnsubDate:             "",
		},
		client: s.client,
	})

}

func (s *ListSuite) Test_GetLists_BadResponse(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   `{ bad json response`,
	})

	List, err := s.client.GetLists()
	c.Assert(err, check.ErrorMatches, "invalid character.*")
	c.Assert(List, check.IsNil)
}

func (s *ListSuite) Test_GetLists_UnknownResponse(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   111,
		Body:   `{}`,
	})

	List, err := s.client.GetLists()
	c.Assert(err, check.ErrorMatches, "Response error.*")
	c.Assert(List, check.IsNil)
}

// --------------------------------------------------------------
// GetList

func (s *ListSuite) Test_GetList_Normal(c *check.C) {

	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   `{"id":"1510500e0b","name":"Freddie's Favorite Hats","contact":{"company":"MailChimp","address1":"675 Ponce De Leon Ave NE","address2":"Suite 5000","city":"Atlanta","state":"GA","zip":"30308","country":"US","phone":""},"permission_reminder":"","use_archive_bar":true,"campaign_defaults":{"from_name":"Freddie","from_email":"freddie@freddiehats.com","subject":"","language":"en"},"notify_on_subscribe":"","notify_on_unsubscribe":"","date_created":"2015-09-16T14:55:51+00:00","list_rating":0,"email_type_option":true,"subscribe_url_short":"http://eepurl.com/xxxx","subscribe_url_long":"http://freddieshats.usX.list-manage.com/subscribe?u=8d3a3db4d97663a9074efcc16&id=1510500e0b","beamer_address":"usX-xxxx-xxxx@inbound.mailchimp.com","visibility":"pub","modules":[],"stats":{"member_count":0,"unsubscribe_count":0,"cleaned_count":0,"member_count_since_send":0,"unsubscribe_count_since_send":0,"cleaned_count_since_send":0,"campaign_count":0,"campaign_last_sent":"","merge_field_count":2,"avg_sub_rate":0,"avg_unsub_rate":0,"target_sub_rate":0,"open_rate":0,"click_rate":0,"last_sub_date":"","last_unsub_date":""}}`,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/1510500e0b")
		},
	})

	list, err := s.client.GetList("1510500e0b")
	c.Assert(err, check.IsNil)
	c.Assert(list.client, check.Not(check.IsNil))
	c.Assert(list, check.DeepEquals, &List{
		ID:   "1510500e0b",
		Name: "Freddie's Favorite Hats",
		Contact: Contact{
			Company:  "MailChimp",
			Address1: "675 Ponce De Leon Ave NE",
			Address2: "Suite 5000",
			City:     "Atlanta",
			State:    "GA",
			Zip:      "30308",
			Country:  "US",
			Phone:    "",
		},
		PermissionReminder: "",
		UseArchiveBar:      true,
		CampaignDefaults: CampaignDefaults{
			FromName:  "Freddie",
			FromEmail: "freddie@freddiehats.com",
			Subject:   "",
			Language:  "en",
		},
		NotifyOnSubscribe:   "",
		NotifyOnUnsubscribe: "",
		DateCreated:         "2015-09-16T14:55:51+00:00",
		ListRating:          0,
		EmailTypeOption:     true,
		SubscribeURLShort:   "http://eepurl.com/xxxx",
		SubscribeURLLong:    "http://freddieshats.usX.list-manage.com/subscribe?u=8d3a3db4d97663a9074efcc16&id=1510500e0b",
		BeamerAddress:       "usX-xxxx-xxxx@inbound.mailchimp.com",
		Visibility:          "pub",
		Modules:             []interface{}{},
		Stats: ListStats{
			MemberCount:               0,
			UnsubscribeCount:          0,
			CleanedCount:              0,
			MemberCountSinceSend:      0,
			UnsubscribeCountSinceSend: 0,
			CleanedCountSinceSend:     0,
			CampaignCount:             0,
			CampaignLastSent:          "",
			MergeFieldCount:           2,
			AvgSubRate:                0,
			AvgUnsubRate:              0,
			TargetSubRate:             0,
			OpenRate:                  0,
			ClickRate:                 0,
			LastSubDate:               "",
			LastUnsubDate:             "",
		},

		client: s.client,
	})

}

func (s *ListSuite) Test_GetList_BadResponse(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   200,
		Body:   `{ bad json response`,
	})

	list, err := s.client.GetList("1510500e0b")
	c.Assert(err, check.ErrorMatches, "invalid character.*")
	c.Assert(list, check.IsNil)
}

func (s *ListSuite) Test_GetList_UnknownResponse(c *check.C) {
	s.server.AddResponse(&t.MockResponse{
		Method: "GET",
		Code:   111,
		Body:   `{}`,
	})

	list, err := s.client.GetList("1510500e0b")
	c.Assert(err, check.ErrorMatches, "Response error.*")
	c.Assert(list, check.IsNil)
}

// --------------------------------------------------------------
// Delete

func (s *ListSuite) Test_Delete_Normal(c *check.C) {

	list := &List{
		ID:     "1510500e0b",
		client: s.client,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "DELETE",
		Code:   http.StatusNoContent,
		Body:   ``,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/1510500e0b")
		},
	})

	err := list.Delete()
	c.Assert(err, check.IsNil)

}

func (s *ListSuite) Test_Delete_NoClient(c *check.C) {
	list := &List{
		ID: "1510500e0b",
	}
	err := list.Delete()
	c.Assert(err, check.ErrorMatches, "no client assigned by parent")
}

func (s *ListSuite) Test_Delete_UnknownResponse(c *check.C) {
	list := &List{
		ID:     "1510500e0b",
		client: s.client,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "DELETE",
		Code:   111,
		Body:   `{}`,
	})

	err := list.Delete()
	c.Assert(err, check.ErrorMatches, "Response error.*")

}

// --------------------------------------------------------------
// Update

func (s *ListSuite) Test_Update_Normal(c *check.C) {

	list := &List{
		ID: "1510500e0b",

		client: s.client,
	}

	update := &UpdateList{
		Name: "Updated",
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "PATCH",
		Code:   200,
		Body:   `{"id":"1510500e0b","name":"Updated","contact":{"company":"MailChimp","address1":"675 Ponce De Leon Ave NE","address2":"Suite 5000","city":"Atlanta","state":"GA","zip":"30308","country":"US","phone":""},"permission_reminder":"","use_archive_bar":true,"campaign_defaults":{"from_name":"Freddie","from_email":"freddie@freddiehats.com","subject":"","language":"en"},"notify_on_subscribe":"","notify_on_unsubscribe":"","date_created":"2015-09-16T14:55:51+00:00","list_rating":0,"email_type_option":true,"subscribe_url_short":"http://eepurl.com/xxxx","subscribe_url_long":"http://freddieshats.usX.list-manage.com/subscribe?u=8d3a3db4d97663a9074efcc16&id=1510500e0b","beamer_address":"usX-xxxx-xxxx@inbound.mailchimp.com","visibility":"pub","modules":[],"stats":{"member_count":0,"unsubscribe_count":0,"cleaned_count":0,"member_count_since_send":0,"unsubscribe_count_since_send":0,"cleaned_count_since_send":0,"campaign_count":0,"campaign_last_sent":"","merge_field_count":2,"avg_sub_rate":0,"avg_unsub_rate":0,"target_sub_rate":0,"open_rate":0,"click_rate":0,"last_sub_date":"","last_unsub_date":""}}`,
		CheckFn: func(r *http.Request, body string) {
			c.Assert(body, check.Equals, `{"name":"Updated"}`)
			c.Assert(r.RequestURI, check.Equals, "http://us13.api.mailchimp.com/3.0/lists/1510500e0b")
		},
	})

	upd, err := list.Update(update)
	c.Assert(err, check.IsNil)
	c.Assert(upd, check.DeepEquals, &List{
		ID:   "1510500e0b",
		Name: "Updated",
		Contact: Contact{
			Company:  "MailChimp",
			Address1: "675 Ponce De Leon Ave NE",
			Address2: "Suite 5000",
			City:     "Atlanta",
			State:    "GA",
			Zip:      "30308",
			Country:  "US",
			Phone:    "",
		},
		PermissionReminder: "",
		UseArchiveBar:      true,
		CampaignDefaults: CampaignDefaults{
			FromName:  "Freddie",
			FromEmail: "freddie@freddiehats.com",
			Subject:   "",
			Language:  "en",
		},
		NotifyOnSubscribe:   "",
		NotifyOnUnsubscribe: "",
		DateCreated:         "2015-09-16T14:55:51+00:00",
		ListRating:          0,
		EmailTypeOption:     true,
		SubscribeURLShort:   "http://eepurl.com/xxxx",
		SubscribeURLLong:    "http://freddieshats.usX.list-manage.com/subscribe?u=8d3a3db4d97663a9074efcc16&id=1510500e0b",
		BeamerAddress:       "usX-xxxx-xxxx@inbound.mailchimp.com",
		Visibility:          "pub",
		Modules:             []interface{}{},
		Stats: ListStats{
			MemberCount:               0,
			UnsubscribeCount:          0,
			CleanedCount:              0,
			MemberCountSinceSend:      0,
			UnsubscribeCountSinceSend: 0,
			CleanedCountSinceSend:     0,
			CampaignCount:             0,
			CampaignLastSent:          "",
			MergeFieldCount:           2,
			AvgSubRate:                0,
			AvgUnsubRate:              0,
			TargetSubRate:             0,
			OpenRate:                  0,
			ClickRate:                 0,
			LastSubDate:               "",
			LastUnsubDate:             "",
		},

		client: s.client,
	})
}

func (s *ListSuite) Test_Update_Missing_Client(c *check.C) {
	list := &List{
		ID: "1510500e0b",
	}
	update := &UpdateList{}
	_, err := list.Update(update)
	c.Assert(err, check.ErrorMatches, "no client assigned by parent")
}

func (s *ListSuite) Test_Update_BadResponse(c *check.C) {
	updSegm := &UpdateList{
		Name: "bad",
	}
	list := &List{
		ID:     "1510500e0b",
		client: s.client,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "PATCH",
		Code:   200,
		Body:   `{ bad json response`,
	})

	upd, err := list.Update(updSegm)
	c.Assert(err, check.ErrorMatches, "invalid character.*")
	c.Assert(upd, check.IsNil)
}

func (s *ListSuite) Test_Update_UnknownResponse(c *check.C) {
	updSegm := &UpdateList{
		Name: "bad",
	}
	list := &List{
		ID:     "1510500e0b",
		client: s.client,
	}

	s.server.AddResponse(&t.MockResponse{
		Method: "PATCH",
		Code:   111,
		Body:   `{}`,
	})

	upd, err := list.Update(updSegm)
	c.Assert(err, check.ErrorMatches, "Response error.*")
	c.Assert(upd, check.IsNil)
}
