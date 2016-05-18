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

package testing

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type MockClient struct{}

func (m *MockClient) Debug(set ...bool) bool {
	return false
}

func (m *MockClient) Log(level ...logrus.Level) *logrus.Logger {
	return logrus.New()
}

func (m *MockClient) Get(resource string, parameters map[string]interface{}) ([]byte, error) {
	return []byte("{}"), nil
}

func (m *MockClient) Post(resource string, parameters map[string]interface{}, data interface{}) ([]byte, error) {
	return []byte("{}"), nil
}

func (m *MockClient) Patch(resource string, parameters map[string]interface{}, data interface{}) ([]byte, error) {
	return []byte("{}"), nil
}

func (m *MockClient) Put(resource string, parameters map[string]interface{}, data interface{}) ([]byte, error) {
	return []byte("{}"), nil
}

func (m *MockClient) Delete(resource string) error {
	return nil
}

func (m *MockClient) Do(request *http.Request) ([]byte, error) {
	return []byte("{}"), nil
}
