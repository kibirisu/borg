package domain

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
)

type Type string

const (
	ObjectType Type = "type"
	LinkType   Type = "link"
	NullType   Type = "null"
)

type ObjectOrLink struct {
	Object *Object
	Link   *string
}

func (o ObjectOrLink) GetType() Type {
	if o.Object != nil {
		return ObjectType
	}
	if o.Link != nil {
		return LinkType
	}
	return NullType
}

func (o *ObjectOrLink) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	switch dec.PeekKind() {
	case '{':
		return json.UnmarshalDecode(dec, &o.Object)
	case '"':
		return json.UnmarshalDecode(dec, &o.Link)
	case 'n':
		return dec.SkipValue()
	default:
		return errors.New("expected JSON object, string or null")
	}
}

func (o *ObjectOrLink) MarshalJSONTo(enc *jsontext.Encoder) error {
	if o.Object != nil {
		return json.MarshalEncode(enc, o.Object)
	}
	if o.Link != nil {
		return json.MarshalEncode(enc, o.Link)
	}
	return enc.WriteToken(jsontext.Null)
}
