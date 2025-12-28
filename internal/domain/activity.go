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
	Context     any
	ID          string
	Type        string
	Actor       ObjectOrLink[Actor]
	Object      ObjectOrLink[Object]
	Publication Maybe[Publication]
	Extra       map[string]any
}

type Object struct {
	ID          string
	Type        string
	Actor       ObjectOrLink[Actor]
	Object      ObjectOrLink[Object]
	Publication Maybe[Publication]
	Note        Maybe[Note]
	Extra       map[string]any
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
	Published    time.Time
	AttributedTo ObjectOrLink[Actor]
	To           []string
	CC           []string
}

type Note struct {
	Content   string
	InReplyTo ObjectOrLink[Object]
	Replies   ObjectOrLink[Object]
}

type Maybe[T any] struct {
	Maybe *T
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
				a.Publication.Maybe = &Publication{}
			}
		case "actor":
			if err = json.UnmarshalDecode(dec, &a.Actor); err != nil {
				return err
			}
		case "published":
			if err = json.UnmarshalDecode(dec, &a.Publication.Maybe.Published); err != nil {
				return err
			}
		case "to":
			if err = json.UnmarshalDecode(dec, &a.Publication.Maybe.To); err != nil {
				return err
			}
		case "cc":
			if err = json.UnmarshalDecode(dec, &a.Publication.Maybe.CC); err != nil {
				return err
			}
		case "attributedTo":
			if err = json.UnmarshalDecode(dec, &a.Publication.Maybe.AttributedTo); err != nil {
				return err
			}
		case "object":
			if err = json.UnmarshalDecode(dec, &a.Object); err != nil {
				return err
			}
		default:
			var extra any
			if err = json.UnmarshalDecode(dec, &extra); err != nil {
				return err
			}
			if a.Extra == nil {
				a.Extra = make(map[string]any)
			}
			a.Extra[property] = &extra
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
			case string(ActivityTypeCreate):
				o.Note.Maybe = &Note{}
				fallthrough
			case string(ActivityTypeAnnounce):
				o.Publication.Maybe = &Publication{}
			}
			if o.Type == string(ActivityTypeCreate) {
				o.Note.Maybe = &Note{}
			}
		case "published":
			if err = json.UnmarshalDecode(dec, &o.Publication.Maybe.Published); err != nil {
				return err
			}
		case "to":
			if err = json.UnmarshalDecode(dec, &o.Publication.Maybe.To); err != nil {
				return err
			}
		case "cc":
			if err = json.UnmarshalDecode(dec, &o.Publication.Maybe.CC); err != nil {
				return err
			}
		case "attributedTo":
			if err = json.UnmarshalDecode(dec, &o.Publication.Maybe.AttributedTo); err != nil {
				return err
			}
		case "content":
			if err = json.UnmarshalDecode(dec, &o.Note.Maybe.Content); err != nil {
				return err
			}
		case "inReplyTo":
			if err = json.UnmarshalDecode(dec, &o.Note.Maybe.InReplyTo); err != nil {
				return err
			}
		case "replies":
			if err = json.UnmarshalDecode(dec, &o.Note.Maybe.Replies); err != nil {
				return err
			}
		case "object":
			if err = json.UnmarshalDecode(dec, &o.Object); err != nil {
				return err
			}
		default:
			var extra any
			if err = json.UnmarshalDecode(dec, &extra); err != nil {
				return err
			}
			if o.Extra == nil {
				o.Extra = make(map[string]any)
			}
			o.Extra[property] = &extra
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
	default:
		return errors.New("expected JSON object or string")
	}
	return nil
}
