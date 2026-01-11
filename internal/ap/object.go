package ap

import "github.com/kibirisu/borg/internal/domain"

type Objecter[T any] interface {
	ObjectOrLink[T]
}

type Object struct {
	ID   string
	Type string
}

type object struct {
	raw *domain.ObjectOrLink
}

var _ Objecter[any] = (*object)(nil)

// GetObject implements Objecter.
func (o *object) GetObject() any {
	obj := o.raw.Object
	return Object{
		ID:   obj.ID,
		Type: obj.Type,
	}
}

// GetURI implements Objecter.
func (o *object) GetURI() string {
	if o.raw == nil {
		return ""
	}
	switch o.raw.GetType() {
	case domain.LinkType:
		return *o.raw.Link
	case domain.NullType:
		return ""
	case domain.ObjectType:
		return o.raw.Object.ID
	default:
		return ""
	}
}

// GetRaw implements Objecter.
func (o *object) GetRaw() *domain.ObjectOrLink {
	return o.raw
}

// GetLink implements Objecter.
func (o *object) GetLink() string {
	return *o.raw.Link
}

// GetValueType implements Objecter.
func (o *object) GetValueType() ValueType {
	if o.raw == nil {
		return NullType
	}
	switch o.raw.GetType() {
	case domain.LinkType:
		return LinkType
	case domain.NullType:
		return NullType
	case domain.ObjectType:
		return ObjectType
	default:
		return InvalidType
	}
}

// SetNull implements Objecter.
func (o *object) SetNull() {
	*o.raw = domain.ObjectOrLink{}
}

// SetObject implements Objecter.
func (o *object) SetObject(object any) {
	obj := object.(Object)
	o.raw = &domain.ObjectOrLink{
		Object: &domain.Object{
			ID:   obj.ID,
			Type: obj.Type,
		},
	}
}

// SetLink implements Objecter.
func (o *object) SetLink(link string) {
	o.raw = &domain.ObjectOrLink{
		Link: &link,
	}
}
