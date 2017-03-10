// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package mailchimp

import (
	"context"
	"net/http"
)

type MailchimpClient interface {
	Get(ctx context.Context, resource string, parameters map[string]interface{}) ([]byte, error)
	Post(ctx context.Context, resource string, parameters map[string]interface{}, data interface{}) ([]byte, error)
	Patch(ctx context.Context, resource string, parameters map[string]interface{}, data interface{}) ([]byte, error)
	Put(ctx context.Context, resource string, parameters map[string]interface{}, data interface{}) ([]byte, error)
	Delete(ctx context.Context, resource string) error

	Do(request *http.Request) ([]byte, error)
}
