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

func (s *CommonSuite) Test_HasField_FieldInStruct(c *check.C) {

	ok := TestTypeA{StringField: "test"}
	empty := TestTypeA{}

	c.Assert(hasField(ok, "StringField"), check.IsNil)
	c.Assert(hasField(empty, "StringField"), check.ErrorMatches, "missing field: StringField")

}

func (s *CommonSuite) Test_HasField_TimeField(c *check.C) {

	now := time.Now()
	ok := TestTypeA{TimeField: now, TimePtr: &now}

	empty := TestTypeA{}
	empty2 := TestTypeA{TimeField: time.Time{}, TimePtr: nil}

	c.Assert(hasField(ok, "TimeField"), check.IsNil)
	c.Assert(hasField(ok, "TimePtr"), check.IsNil)

	c.Assert(hasField(empty, "TimeField"), check.ErrorMatches, "missing field: TimeField")
	c.Assert(hasField(empty, "TimePtr"), check.ErrorMatches, "missing field: TimePtr")

	c.Assert(hasField(empty2, "TimeField"), check.ErrorMatches, "missing field: TimeField")
	c.Assert(hasField(empty2, "TimePtr"), check.ErrorMatches, "missing field: TimePtr")

}

func (s *CommonSuite) Test_HasField_ObjAndPtr(c *check.C) {

	ok := TestTypeA{
		PtrField:   &TestTypeB{Content: "Test"},
		ValueField: TestTypeB{Content: "Test"},
	}

	empty := TestTypeA{}

	c.Assert(hasField(ok, "PtrField"), check.IsNil)
	c.Assert(hasField(ok, "ValueField"), check.IsNil)

	c.Assert(hasField(empty, "PtrField"), check.ErrorMatches, "missing field: PtrField")
	c.Assert(hasField(empty, "ValueField"), check.ErrorMatches, "missing field: ValueField")

}

func (s *CommonSuite) Test_HasFields_FieldInStruct(c *check.C) {
	ok := TestTypeA{
		PtrField:   &TestTypeB{Content: "Test"},
		ValueField: TestTypeB{Content: "Test"},
	}

	empty := TestTypeA{}

	c.Assert(hasFields(ok, "PtrField", "ValueField"), check.IsNil)

	c.Assert(hasFields(empty, "PtrField", "ValueField"), check.ErrorMatches, "missing field: PtrField")

	empty.PtrField = ok.PtrField
	c.Assert(hasFields(empty, "PtrField", "ValueField"), check.ErrorMatches, "missing field: ValueField")
}
