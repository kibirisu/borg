package domain

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"time"
)

type ActivityType string

const (
	ActivityTypeAccept ActivityType = "Accept"
	ActivityTypeCreate ActivityType = "Create"
	ActivityTypeFollow ActivityType = "Follow"
	ActivityTypeUnimpl ActivityType = "Unimpl"
)

type Object struct {
	Context any    `json:"@context"`
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
	Base   Object         `json:",inline"`
	Actor  Actor          `json:"actor"`
	Object ActivityObject `json:"object"`
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
	for {
		tok, err := dec.ReadToken()
		if err != nil {
			return err
		}
		if tok.Kind() == '"' && tok.String() == "type" {
			t, err := dec.ReadToken()
			if err != nil {
				return err
			}
			// WARNING: works when assuming "string" type
			o.Type = ActivityType(t.String())
			break
		}
	}
	// TODO: decode object
	switch o.Type {
	case ActivityTypeAccept:
		var follow Activity
		o.Object = &follow
	case ActivityTypeCreate:
		var note Note
		o.Object = &note
		_ = note
	case ActivityTypeFollow:
		var actor Actor
		o.Object = &actor
		_ = actor
	default:
		o.Type = ActivityTypeUnimpl
	}
	return nil
}
