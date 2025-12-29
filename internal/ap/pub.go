package ap

import "github.com/kibirisu/borg/internal/domain"

type Pub interface {
	GetActor()
	GetObject()
	SetActor()
	SetObject()
}

type APObjectOrLink[T any] interface {
	IsValid() bool
	GetObject() T
	GetURI() string
	SetObject(T)
	SetURI(string)
}

var _ Pub = (*activityPub)(nil)

type activityPub struct {
	raw domain.Object
}

// GetActor implements Pub.
func (a *activityPub) GetActor() {
	panic("unimplemented")
}

// GetObject implements Pub.
func (a *activityPub) GetObject() {
	panic("unimplemented")
}

// SetActor implements Pub.
func (a *activityPub) SetActor() {
	panic("unimplemented")
}

// SetObject implements Pub.
func (a *activityPub) SetObject() {
	panic("unimplemented")
}
