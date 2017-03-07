// Copyright (C) 2017 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Created by David Högborg <d@greatbeyond.se>, 2017

package mailchimp

import "context"

type contextKey int

const (
	tokenKey contextKey = iota
	urlKey   contextKey = iota
)

func TokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(tokenKey).(string)
	return token, ok
}

func NewContextWithToken(ctx context.Context, token string) context.Context {
	return context.WithValue(ctx, tokenKey, token)
}

func URLFromContext(ctx context.Context) (string, bool) {
	url, ok := ctx.Value(urlKey).(string)
	return url, ok
}

func NewContextWithURL(ctx context.Context, url string) context.Context {
	return context.WithValue(ctx, urlKey, url)
}
