// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package mailchimp

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

type MailchimpClient interface {
	Debug(set ...bool) bool
	Log() *logrus.Logger

	Get(resource string, parameters map[string]interface{}) ([]byte, error)
	Post(resource string, parameters map[string]interface{}, data interface{}) ([]byte, error)
	Patch(resource string, parameters map[string]interface{}, data interface{}) ([]byte, error)
	Put(resource string, parameters map[string]interface{}, data interface{}) ([]byte, error)
	Delete(resource string) error

	Do(request *http.Request) ([]byte, error)
}
