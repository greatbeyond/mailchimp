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
