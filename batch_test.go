// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package mailchimp

import (
	"bytes"
	"net/http"

	"github.com/Sirupsen/logrus"
	t "github.com/greatbeyond/mailchimp/testing"

	check "gopkg.in/check.v1"
)

var _ = check.Suite(&BatchTestSuite{})

type BatchTestSuite struct {
	Batch *BatchQueue
}

func (suite *BatchTestSuite) SetUpSuite(c *check.C) {
	Log.Level = logrus.DebugLevel
}

func (suite *BatchTestSuite) SetUpTest(c *check.C) {
	suite.Batch = &BatchQueue{
		client: &t.MockClient{},
	}
}

func (suite *BatchTestSuite) TearDownTest(c *check.C) {}

func (suite *BatchTestSuite) Test_Batch_Do_GET(c *check.C) {

	req, _ := http.NewRequest("GET", "http://example.net/3.0/resoruce/id", nil)
	resp, err := suite.Batch.Do(req)
	c.Assert(err, check.IsNil)

	c.Assert(resp, check.DeepEquals, []byte("{}"))
	c.Assert(len(suite.Batch.Operations), check.Equals, 1)
	c.Assert(suite.Batch.Operations[0], check.DeepEquals, &BatchOperation{
		Method: "GET",
		Path:   "/resoruce/id",
		Params: nil,
		Body:   "",
	})

}

func (suite *BatchTestSuite) Test_Batch_Do_GET_Parms(c *check.C) {

	req, _ := http.NewRequest("GET", "http://example.net/3.0/resoruce/id?r=K", nil)
	resp, err := suite.Batch.Do(req)
	c.Assert(err, check.IsNil)

	c.Assert(resp, check.DeepEquals, []byte("{}"))
	c.Assert(len(suite.Batch.Operations), check.Equals, 1)
	c.Assert(suite.Batch.Operations[0], check.DeepEquals, &BatchOperation{
		Method: "GET",
		Path:   "/resoruce/id",
		Params: map[string]string{
			"r": "K",
		},
		Body: "",
	})

}

func (suite *BatchTestSuite) Test_Batch_Do_POST(c *check.C) {

	body := bytes.NewBuffer([]byte(`{"key":"value"}`))
	req, _ := http.NewRequest("POST", "http://example.net/3.0/resoruce/id", body)

	resp, err := suite.Batch.Do(req)
	c.Assert(err, check.IsNil)

	c.Assert(resp, check.DeepEquals, []byte("{}"))
	c.Assert(len(suite.Batch.Operations), check.Equals, 1)
	c.Assert(suite.Batch.Operations[0], check.DeepEquals, &BatchOperation{
		Method: "POST",
		Path:   "/resoruce/id",
		Params: nil,
		Body:   `{"key":"value"}`,
	})

}

func (suite *BatchTestSuite) Test_Batch_CreateMember(c *check.C) {

	client := NewClient("arandomtoken-us0")
	client.NewBatch()

	_, err := client.CreateMember(&CreateMember{
		EmailAddress: "test@example.net",
	}, "123456")
	c.Assert(err, check.IsNil)

	c.Assert(len(client.Batch.Operations), check.Equals, 1)
	c.Assert(client.Batch.Operations[0], check.DeepEquals, &BatchOperation{
		Method: "POST",
		Path:   "/lists/123456/members",
		Body:   `{"email_address":"test@example.net"}`,
	})

}

func (suite *BatchTestSuite) Test_Batch_CreateMergeField(c *check.C) {

	client := NewClient("arandomtoken-us0")
	client.NewBatch()

	_, err := client.CreateMergeField(&CreateMergeField{
		Name: "TAGNAME",
		Type: MergeFieldTypeText,
	}, "123456")
	c.Assert(err, check.IsNil)

	c.Assert(len(client.Batch.Operations), check.Equals, 1)
	c.Assert(client.Batch.Operations[0], check.DeepEquals, &BatchOperation{
		Method: "POST",
		Path:   "/lists/123456/merge-fields",
		Body:   `{"name":"TAGNAME","type":"text"}`,
	})

}
