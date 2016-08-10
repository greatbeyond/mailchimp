// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package mailchimp

import (
	"time"

	check "gopkg.in/check.v1"
)

var _ = check.Suite(&CommonSuite{})

type CommonSuite struct {
}

func (s *CommonSuite) SetUpSuite(c *check.C) {

}

func (s *CommonSuite) SetUpTest(c *check.C) {

}

func (s *CommonSuite) TearDownTest(c *check.C) {}

func (s *CommonSuite) Test_SingleJoiningSlash(c *check.C) {

	c.Assert(singleJoiningSlash("no", "slash"), check.Equals, "no/slash")
	c.Assert(singleJoiningSlash("single/", "slash"), check.Equals, "single/slash")
	c.Assert(singleJoiningSlash("single", "/slash"), check.Equals, "single/slash")
	c.Assert(singleJoiningSlash("dual/", "/slash"), check.Equals, "dual/slash")

}

func (s *CommonSuite) Test_SlashJoin(c *check.C) {
	c.Assert(slashJoin("a", "b", "c", "d"), check.Equals, "a/b/c/d")
	c.Assert(slashJoin("a/", "/b", "c/", "d"), check.Equals, "a/b/c/d")
	c.Assert(slashJoin("a/", "b", "c", "d/"), check.Equals, "a/b/c/d")
	c.Assert(slashJoin("//a/", "b", "c", "d//"), check.Equals, "/a/b/c/d/")
}

type TestTypeA struct {
	StringField string
	TimeField   time.Time
	TimePtr     *time.Time
	PtrField    *TestTypeB
	ValueField  TestTypeB
}

type TestTypeB struct {
	Content string
}

func (s *CommonSuite) Test_MissingField_FieldInStruct(c *check.C) {

	ok := TestTypeA{StringField: "test"}
	empty := TestTypeA{}

	c.Assert(missingField(ok, "StringField"), check.IsNil)
	c.Assert(missingField(empty, "StringField"), check.ErrorMatches, "missing field: StringField")

}

func (s *CommonSuite) Test_MissingField_TimeField(c *check.C) {

	now := time.Now()
	ok := TestTypeA{TimeField: now, TimePtr: &now}

	empty := TestTypeA{}
	empty2 := TestTypeA{TimeField: time.Time{}, TimePtr: nil}

	c.Assert(missingField(ok, "TimeField"), check.IsNil)
	c.Assert(missingField(ok, "TimePtr"), check.IsNil)

	c.Assert(missingField(empty, "TimeField"), check.ErrorMatches, "missing field: TimeField")
	c.Assert(missingField(empty, "TimePtr"), check.ErrorMatches, "missing field: TimePtr")

	c.Assert(missingField(empty2, "TimeField"), check.ErrorMatches, "missing field: TimeField")
	c.Assert(missingField(empty2, "TimePtr"), check.ErrorMatches, "missing field: TimePtr")

}

func (s *CommonSuite) Test_MissingField_ObjAndPtr(c *check.C) {

	ok := TestTypeA{
		PtrField:   &TestTypeB{Content: "Test"},
		ValueField: TestTypeB{Content: "Test"},
	}

	empty := TestTypeA{}

	c.Assert(missingField(ok, "PtrField"), check.IsNil)
	c.Assert(missingField(ok, "ValueField"), check.IsNil)

	c.Assert(missingField(empty, "PtrField"), check.ErrorMatches, "missing field: PtrField")
	c.Assert(missingField(empty, "ValueField"), check.ErrorMatches, "missing field: ValueField")

}
