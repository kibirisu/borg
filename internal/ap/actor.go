package ap

import "github.com/kibirisu/borg/internal/domain"

type Actorer interface {
	ObjectOrLink[Actor]
}

type Actor struct {
	ID                string
	Type              string
	PreferredUsername string
	Inbox             string
	Outbox            string
	Following         string
	Followers         string
}

type actor struct {
	part *domain.ObjectOrLink
}

var _ Actorer = (*actor)(nil)

// GetObject implements Actorer.
func (a *actor) GetObject() Actor {
	obj := a.part.Object
	actor := obj.Actor
	return Actor{
		ID:                obj.ID,
		Type:              obj.Type,
		PreferredUsername: actor.PreferredUsername,
		Inbox:             actor.Inbox,
		Outbox:            actor.Outbox,
		Following:         actor.Following,
		Followers:         actor.Followers,
	}
}

// GetURI implements Actorer.
func (a *actor) GetURI() string {
	return *a.part.Link
}

// GetValueType implements Actorer.
func (a *actor) GetValueType() ValueType {
	switch a.part.GetType() {
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

// SetNull implements Actorer.
func (a *actor) SetNull() {
	a.part.Link = nil
	a.part.Object = nil
}

// SetObject implements Actorer.
func (a *actor) SetObject(actor Actor) {
	a.part.Object.ID = actor.ID
	a.part.Object.Type = actor.Type
	a.part.Object.Actor = &domain.Actor{
		PreferredUsername: actor.PreferredUsername,
		Inbox:             actor.Inbox,
		Outbox:            actor.Outbox,
		Following:         actor.Following,
		Followers:         actor.Followers,
	}
}

// SetURI implements Actorer.
func (a *actor) SetURI(uri string) {
	a.part.Link = &uri
}

// GetRaw implements Actorer.
func (a *actor) GetRaw() *domain.ObjectOrLink {
	return a.part
}
