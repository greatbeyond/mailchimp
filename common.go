// Copyright (C) 2016 Great Beyond AB - All Rights Reserved
// Unauthorized copying of this file, via any medium is strictly prohibited
// Proprietary and confidential
// Written by David Högborg <d@greatbeyond.se>, 2016

package mailchimp

import (
	"crypto/md5"
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// MemberEmailToID converts member email address to mailchimp ID (md5 hashed lowercase version of email address)
func MemberEmailToID(email string) string {
	lcemail := strings.ToLower(email)
	hash := md5.Sum([]byte(lcemail))
	return fmt.Sprintf("%x", hash)
}

// from go httputils
func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

// slashJoin converts slashJoin("/hello", "world", "how/", "/you/", "doing/") to
// "hello/world/how/you/doing". A leading or trailing / is removed if present.
func slashJoin(components ...string) string {
	for i, c := range components {
		if strings.HasPrefix(c, "/") {
			c = c[1:]
		}
		if strings.HasSuffix(c, "/") {
			c = c[:len(c)-1]
		}
		components[i] = c
	}
	return strings.Join(components, "/")
}

const TimeFormat = "2006-01-02 15:04:05"

func TimeToString(t time.Time) string {
	return t.Format(TimeFormat)
}

func StringToTime(str string) (time.Time, error) {
	return time.Parse(TimeFormat, str)
}

func missingField(field interface{}, name string) error {
	err := fmt.Errorf("missing field: %s", name)
	switch f := field.(type) {
	case string:
		if f == "" {
			return err
		}
	default:
		if f == nil {
			return err
		}
	}

	return nil
}

func requestParameters(filters []Parameters) map[string]interface{} {
	params := map[string]interface{}{}
	for _, filter := range filters {
		for k, v := range filter {
			params[k] = v
		}
	}

	return params
}

func caller() string {
	_, file, line, _ := runtime.Caller(1)
	matcher := regexp.MustCompile("^(.*)/(.*?)\\.go$")
	matches := matcher.FindAllStringSubmatch(file, -1)
	msg := fmt.Sprintf(" [%s.go:%d]", matches[0][2], line)

	return msg
}
