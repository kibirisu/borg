package ap

import "github.com/kibirisu/borg/internal/domain"

type ActivityPubActor interface {
	APObjectOrLink[Actor]
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

type Activity struct {
	object domain.Object
}

var _ ActivityPubActor = (*Actor)(nil)

// GetObject implements ActivityPubActor.
func (a *Actor) GetObject() Actor {
	panic("unimplemented")
}

// GetURI implements ActivityPubActor.
func (a *Actor) GetURI() string {
	panic("unimplemented")
}

// IsValid implements ActivityPubActor.
func (a *Actor) IsValid() bool {
	panic("unimplemented")
}

// SetObject implements ActivityPubActor.
func (a *Actor) SetObject(Actor) {
	panic("unimplemented")
}

// SetURI implements ActivityPubActor.
func (a *Actor) SetURI(string) {
	panic("unimplemented")
}
