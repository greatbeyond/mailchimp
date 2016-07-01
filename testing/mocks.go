// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package testing

import (
	"net/http"

	"github.com/Sirupsen/logrus"
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
