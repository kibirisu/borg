package ap

import "github.com/kibirisu/borg/internal/domain"

type Actorer interface {
	Objecter[Actor]
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
	object
}

var _ Actorer = (*actor)(nil)

func NewActor(from *domain.ObjectOrLink) Actorer {
	return &actor{object{from}}
}

func NewEmptyActor() Actorer {
	return &actor{object{}}
}

// GetObject implements Actorer.
// Subtle: this method shadows the method (object).GetObject of actor.object.
func (a *actor) GetObject() Actor {
	obj := a.raw.Object
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

// SetObject implements Actorer.
// Subtle: this method shadows the method (object).SetObject of actor.object.
func (a *actor) SetObject(actor Actor) {
	a.raw = &domain.ObjectOrLink{
		Object: &domain.Object{
			ID:   actor.ID,
			Type: actor.Type,
			Actor: &domain.Actor{
				PreferredUsername: actor.PreferredUsername,
				Inbox:             actor.Inbox,
				Outbox:            actor.Outbox,
				Following:         actor.Following,
				Followers:         actor.Followers,
			},
		},
	}
}

// WithLink implements Actorer.
// Subtle: this method shadows the method (object).WithLink of actor.object.
func (a *actor) WithLink(link string) Objecter[Actor] {
	a.SetLink(link)
	return a
}

// WithObject implements Actorer.
// Subtle: this method shadows the method (object).WithObject of actor.object.
func (a *actor) WithObject(actor Actor) Objecter[Actor] {
	a.SetObject(actor)
	return a
}
