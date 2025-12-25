//go:build goexperiment.jsonv2

package domain

import (
	"bytes"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"io"
	"time"
)

type ActivityType string

const (
	ActivityTypePerson ActivityType = "Person"
	ActivityTypeNote   ActivityType = "Note"
	ActivityTypeFollow ActivityType = "Follow"
	ActivityTypeUnimpl ActivityType = "Unimpl"
)

type Object struct {
	Context any    `json:"@context,omitempty"`
	Type    string `json:"type"`
	ID      string `json:"id"`
}

type Actor struct {
	Base              Object `json:",inline"`
	PreferredUsername string `json:"preferredUsername"`
	Inbox             string `json:"inbox"`
	Outbox            string `json:"outbox"`
	Following         string `json:"following"`
	Followers         string `json:"followers"`
}

type Activity struct {
	Object `json:",inline"`
	Actor  Actor          `json:"actor"`
	Target ActivityObject `json:"object"`
}

type ActivityObject struct {
	Type   ActivityType
	Object any
}

type NoteProperties struct {
	Published    time.Time `json:"published"`
	AttributedTo string    `json:"attributedTo"`
	To           []string  `json:"to"`
	CC           []string  `json:"cc"`
	Content      string    `json:"content"`
}

type Note struct {
	Base       Object         `json:",inline"`
	Properties NoteProperties `json:",inline"`
	InReplyTo  struct {
		Base       Object         `json:",inline"`
		Properties NoteProperties `json:",inline"`
	} `json:"inReplyTo"`
}

var _ json.UnmarshalerFrom = (*ActivityObject)(nil)

// UnmarshalJSONFrom implements json.UnmarshalerFrom.
func (o *ActivityObject) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	val, err := dec.ReadValue()
	if err != nil {
		return err
	}
	if val.Kind() == '"' {
		// TODO: handle AP Link type
		return nil
	}
	objType, err := scanForType(val)
	if err != nil {
		return err
	}
	switch objType {
	// case Ac
	case ActivityTypePerson:
		var person Actor
		if err := json.Unmarshal(val, &person); err != nil {
			return err
		}
		o.Object = &person
	case ActivityTypeNote:
		var note Note
		if err := json.Unmarshal(val, &note); err != nil {
			return err
		}
		o.Object = &note
	case ActivityTypeFollow:
		var follow Activity
		if err := json.Unmarshal(val, &follow); err != nil {
			return err
		}
		o.Object = &follow
	default:
		o.Type = ActivityTypeUnimpl
	}
	o.Type = objType
	return nil
}

func scanForType(val jsontext.Value) (ActivityType, error) {
	dec := jsontext.NewDecoder(bytes.NewReader(val))
	if dec.PeekKind() != '{' {
		return "", errors.New("expected JSON object")
	}
	if _, err := dec.ReadToken(); err != nil {
		return "", err
	}
	for {
		tok, err := dec.ReadToken()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
		// TODO: check if object and discard it
		if tok.Kind() == '"' && tok.String() == "type" {
			typeToken, err := dec.ReadToken()
			if err != nil {
				return "", err
			}
			if typeToken.Kind() == '"' {
				return ActivityType(typeToken.String()), nil
			}
			return "", errors.New("expected JSON string")
		}
	}
	return ActivityTypeUnimpl, nil
}
