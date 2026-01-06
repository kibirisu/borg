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

// GetRaw implements Objecter.
func (o *object) GetRaw() *domain.ObjectOrLink {
	return o.raw
}

// GetURI implements Objecter.
func (o *object) GetURI() string {
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
		panic("unexpected domain.Type")
	}
}

// SetNull implements Objecter.
func (o *object) SetNull() {
	*o.raw = domain.ObjectOrLink{}
}

// SetObject implements Objecter.
func (o *object) SetObject(object any) {
	obj := object.(Object)
	*o.raw = domain.ObjectOrLink{
		Object: &domain.Object{
			ID:   obj.ID,
			Type: obj.Type,
		},
	}
}

// SetURI implements Objecter.
func (o *object) SetURI(uri string) {
	*o.raw = domain.ObjectOrLink{
		Link: &uri,
	}
}
