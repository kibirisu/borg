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

	ActivityTypeNote ActivityType = "Note"
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

func (o *ObjectOrLink[T]) GetURI() string {
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

func (a *Activity) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if dec.PeekKind() != '{' {
		return errors.New("expected JSON object")
	}
	_, err := dec.ReadToken()
	if err != nil {
		return err
	}
	if err = a.unmarshalProperties(dec); err != nil {
		return err
	}

	return nil
}

func (a *Activity) unmarshalProperties(dec *jsontext.Decoder) error {
	for dec.PeekKind() != '}' {
		if dec.PeekKind() != '"' {
			return errors.New("expected JSON string")
		}
		tok, err := dec.ReadToken()
		if err != nil {
			return err
		}
		property := tok.String()
		switch property {
		case "@context":
			if err = json.UnmarshalDecode(dec, &a.Context); err != nil {
				return err
			}
		case "id":
			if err = json.UnmarshalDecode(dec, &a.ID); err != nil {
				return err
			}
		case "type":
			if err = json.UnmarshalDecode(dec, &a.Type); err != nil {
				return err
			}
			switch a.Type {
			case string(ActivityTypeAnnounce), string(ActivityTypeCreate):
				a.Publication = &Publication{}
			}
		case "actor":
			if err = json.UnmarshalDecode(dec, &a.Actor); err != nil {
				return err
			}
		case "published":
			if err = json.UnmarshalDecode(dec, &a.Publication.Published); err != nil {
				return err
			}
		case "to":
			if err = json.UnmarshalDecode(dec, &a.Publication.To); err != nil {
				return err
			}
		case "cc":
			if err = json.UnmarshalDecode(dec, &a.Publication.CC); err != nil {
				return err
			}
		case "attributedTo":
			if err = json.UnmarshalDecode(dec, &a.Publication.AttributedTo); err != nil {
				return err
			}
		case "object":
			if err = json.UnmarshalDecode(dec, &a.Object); err != nil {
				return err
			}
		default:
			if err = dec.SkipValue(); err != nil {
				return err
			}
		}
	}
	_, err := dec.ReadToken()
	return err
}

func (o *Object) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	if dec.PeekKind() != '{' {
		return errors.New("expected JSON object")
	}
	_, err := dec.ReadToken()
	if err != nil {
		return err
	}
	for dec.PeekKind() != '}' {
		if dec.PeekKind() != '"' {
			return errors.New("expected JSON string")
		}
		tok, err := dec.ReadToken()
		if err != nil {
			return err
		}
		property := tok.String()
		switch property {
		case "id":
			if err = json.UnmarshalDecode(dec, &o.ID); err != nil {
				return err
			}
		case "type":
			if err = json.UnmarshalDecode(dec, &o.Type); err != nil {
				return err
			}
			switch o.Type {
			case string(ActivityTypeNote):
				o.Note = &Note{}
				fallthrough
			case string(ActivityTypeAnnounce):
				o.Publication = &Publication{}
			}
		case "actor":
			if err = json.UnmarshalDecode(dec, &o.Actor); err != nil {
				return err
			}
		case "published":
			if err = json.UnmarshalDecode(dec, &o.Publication.Published); err != nil {
				return err
			}
		case "to":
			if err = json.UnmarshalDecode(dec, &o.Publication.To); err != nil {
				return err
			}
		case "cc":
			if err = json.UnmarshalDecode(dec, &o.Publication.CC); err != nil {
				return err
			}
		case "attributedTo":
			if err = json.UnmarshalDecode(dec, &o.Publication.AttributedTo); err != nil {
				return err
			}
		case "content":
			if err = json.UnmarshalDecode(dec, &o.Note.Content); err != nil {
				return err
			}
		case "inReplyTo":
			var obj ObjectOrLink[Actor]
			if err = json.UnmarshalDecode(dec, &obj); err != nil {
				return err
			}
			o.Publication.AttributedTo = &obj
		case "replies":
			if err = json.UnmarshalDecode(dec, &o.Note.Replies); err != nil {
				return err
			}
		case "object":
			if err = json.UnmarshalDecode(dec, &o.Object); err != nil {
				return err
			}
		default:
			if err = dec.SkipValue(); err != nil {
				return err
			}
		}
	}
	_, err = dec.ReadToken()
	return err
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
