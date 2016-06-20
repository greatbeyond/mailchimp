// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package mailchimp

import (
	"bytes"
	"net/http"
	"testing"

	check "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { check.TestingT(t) }

var _ = check.Suite(&BatchTestSuite{})

type BatchTestSuite struct {
	Batch *BatchQueue
}

func (suite *BatchTestSuite) SetUpSuite(c *check.C) {}

func (suite *BatchTestSuite) SetUpTest(c *check.C) {
	suite.Batch = &BatchQueue{}
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
