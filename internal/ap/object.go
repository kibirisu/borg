package ap

import "github.com/kibirisu/borg/internal/domain"

type Objecter interface {
	ObjectOrLink[Object]
}

type Object struct {
	ID    string
	Type  string
	Actor Actorer
}

type object struct {
	part *domain.ObjectOrLink
}

var _ Objecter = (*object)(nil)

// GetObject implements Objecter.
func (o *object) GetObject() Object {
	obj := o.part.Object
	return Object{
		ID:   obj.ID,
		Type: obj.Type,
		Actor: &actor{
			part: obj.ActivityActor,
		},
	}
}

// GetURI implements Objecter.
func (o *object) GetURI() string {
	return *o.part.Link
}

// GetValueType implements Objecter.
func (o *object) GetValueType() ValueType {
	switch o.part.GetType() {
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
	o.part.Link = nil
	o.part.Object = nil
}

// SetObject implements Objecter.
func (o *object) SetObject(object Object) {
	o.part.Object.ID = object.ID
	o.part.Object.Type = object.Type
	o.part.Object.ActivityActor = object.Actor.GetRaw()
}

// SetURI implements Objecter.
func (o *object) SetURI(uri string) {
	o.part.Link = &uri
}

// GetRaw implements Objecter.
func (o *object) GetRaw() *domain.ObjectOrLink {
	return o.part
}
