package util

import (
	"errors"
	"strings"
)

// FIXME: buggy

func ParseAccount(acct string) (string, string, error) {
	resource := strings.TrimPrefix(acct, "acct:")
	if resource == acct {
		return "", "", errors.New("not a valid resource")
	}
	arr := strings.Split(resource, "@")
	switch len(arr) {
	case 1:
		return "", "", errors.New("not a valid resource")
	case 2:
		return arr[0], arr[1], nil
	default:
		return "", "", errors.New("not a valid resource")
	}
}
