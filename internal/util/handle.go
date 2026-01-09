package util

import (
	"errors"
	"fmt"
	"strings"
)

type HandleInfo struct {
	Username string // user part
	Domain   string // server part, e.g. borg.local
	Local    bool   // true if domain matches local host
}

func ParseHandle(raw string, localHost string) (*HandleInfo, error) {
	handle := strings.TrimSpace(raw)
	if handle == "" {
		return nil, errors.New("empty handle")
	}
	handle = strings.TrimPrefix(handle, "@")
	parts := strings.Split(handle, "@")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid handle format: %s", raw)
	}
	username := parts[0]
	domain := parts[1]
	if username == "" || domain == "" {
		return nil, fmt.Errorf("invalid handle format: %s", raw)
	}
	info := &HandleInfo{
		Username: username,
		Domain:   domain,
		Local:    strings.EqualFold(domain, localHost),
	}
	return info, nil
}
