package domain

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
)

type URIer interface {
	URI() string
}

type ObjectOrLink struct {
	Object *Object
	Link   *string
}

func (o ObjectOrLink) GetURI() string {
	if o.Object != nil {
		return o.Object.ID
	}
	if o.Link != nil {
		return *o.Link
	}
	return ""
}

func (o *ObjectOrLink) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
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

func (o *ObjectOrLink) MarshalJSONTo(enc *jsontext.Encoder) error {
	if o.Object != nil {
		return json.MarshalEncode(enc, o.Object)
	}
	if o.Link != nil {
		return json.MarshalEncode(enc, o.Link)
	}
	return enc.WriteToken(jsontext.Null)
}
