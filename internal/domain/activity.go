//go:build goexperiment.jsonv2

package domain

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"time"
)

type ActivityType string

const (
	ActivityTypeAccept        ActivityType = "Accept"
	ActivityTypeAnnounce      ActivityType = "Announce"
	ActivityTypeCreate        ActivityType = "Create"
	ActivityTypeFollow        ActivityType = "Follow"
	ActivityTypeLike          ActivityType = "Like"
	ActivityTypeUnimplemented ActivityType = "Unimplemented"
)

type URIer interface {
	URI() string
}

type Activity struct {
	Context     any                  `json:"@context"`
	ID          string               `json:"id"`
	Type        string               `json:"type"`
	Publication *Publication         `json:",inline"`
	Actor       ObjectOrLink[Actor]  `json:"actor"`
	Object      ObjectOrLink[Object] `json:"object"`
}

type Object struct {
	ID          string                `json:"id"`
	Type        string                `json:"type"`
	Publication *Publication          `json:",inline"`
	Note        *Note                 `json:",inline"`
	Actor       *ObjectOrLink[Actor]  `json:"actor,omitempty"`
	Object      *ObjectOrLink[Object] `json:"object,omitempty"`
}

type Actor struct {
	ID                string `json:"id"`
	Type              string `json:"type"`
	PreferredUsername string `json:"preferredUsername"`
	Inbox             string `json:"inbox"`
	Outbox            string `json:"outbox"`
	Following         string `json:"following"`
	Followers         string `json:"followers"`
}

type Publication struct {
	Published    time.Time            `json:"published"`
	AttributedTo *ObjectOrLink[Actor] `json:"attributedTo,omitempty"`
	To           []string             `json:"to"`
	CC           []string             `json:"cc"`
}

type Note struct {
	Content   string                `json:"content"`
	InReplyTo *ObjectOrLink[Object] `json:"inReplyTo"`
	Replies   ObjectOrLink[Object]  `json:"replies"`
}

type ObjectOrLink[T URIer] struct {
	Object *T
	Link   *string
}

func (o ObjectOrLink[T]) GetURI() string {
	if o.Object != nil {
		return (*o.Object).URI()
	}
	if o.Link != nil {
		return *o.Link
	}
	return ""
}

func (a Actor) URI() string {
	return a.ID
}

func (o Object) URI() string {
	return o.ID
}

func (o *ObjectOrLink[T]) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch dec.PeekKind() {
	case '{':
		if err := json.UnmarshalDecode(dec, &o.Object); err != nil {
			return err
		}
	case '"':
		if err := json.UnmarshalDecode(dec, &o.Link); err != nil {
			return err
		}
	case 'n':
		if err := dec.SkipValue(); err != nil {
			return err
		}
	default:
		return errors.New("expected JSON object or string")
	}
	return nil
}

func (o *ObjectOrLink[T]) MarshalJSONTo(enc *jsontext.Encoder) error {
	if o.Object != nil {
		return json.MarshalEncode(enc, o.Object)
	}
	if o.Link != nil {
		return json.MarshalEncode(enc, o.Link)
	}
	return enc.WriteToken(jsontext.Null)
}
