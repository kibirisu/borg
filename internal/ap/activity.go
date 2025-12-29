package ap

import "github.com/kibirisu/borg/internal/domain"

type Activiter interface {
	ObjectOrLink[Activity]
}

type Activity struct {
	ID     string
	Type   string
	Actor  Actorer
	Object Objecter
}

type activity struct {
	raw *domain.ObjectOrLink
}

var _ Activiter = (*activity)(nil)

// GetObject implements Activiter.
func (a *activity) GetObject() Activity {
	return Activity{
		ID:     a.raw.Object.ID,
		Type:   a.raw.Object.Type,
		Actor:  &actor{a.raw.Object.ActivityActor},
		Object: &object{a.raw.Object.ActivityObject},
	}
}

// GetRaw implements Activiter.
func (a *activity) GetRaw() *domain.ObjectOrLink {
	return a.raw
}

// GetURI implements Activiter.
func (a *activity) GetURI() string {
	return *a.raw.Link
}

// GetValueType implements Activiter.
func (a *activity) GetValueType() ValueType {
	switch a.raw.GetType() {
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

// SetNull implements Activiter.
func (a *activity) SetNull() {
	a.raw.Link = nil
	a.raw.Object = nil
}

// SetObject implements Activiter.
func (a *activity) SetObject(activity Activity) {
	a.raw.Object.ID = activity.ID
	a.raw.Object.Type = activity.Type
	a.raw.Object.ActivityActor = activity.Actor.GetRaw()
	a.raw.Object.ActivityObject = activity.Object.GetRaw()
}

// SetURI implements Activiter.
func (a *activity) SetURI(uri string) {
	a.raw.Link = &uri
}
