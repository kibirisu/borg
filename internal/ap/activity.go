package ap

import "github.com/kibirisu/borg/internal/domain"

type Activiter interface {
	Objecter[Activity[any]]
}

type Activity[T any] struct {
	ID     string
	Type   string
	Actor  Actorer
	Object Objecter[T]
}

type activity struct {
	object
}

var _ Activiter = (*activity)(nil)

// GetObject implements Activiter.
// Subtle: this method shadows the method (object).GetObject of activity.object.
func (a *activity) GetObject() Activity[any] {
	return Activity[any]{
		ID:     a.raw.Object.ID,
		Type:   a.raw.Object.Type,
		Actor:  &actor{object{a.raw.Object.ActivityActor}},
		Object: &object{a.raw.Object.ActivityObject},
	}
}

// SetObject implements Activiter.
// Subtle: this method shadows the method (object).SetObject of activity.object.
func (a *activity) SetObject(activity Activity[any]) {
	a.raw = &domain.ObjectOrLink{
		Object: &domain.Object{
			ID:             activity.ID,
			Type:           activity.Type,
			ActivityActor:  activity.Actor.GetRaw(),
			ActivityObject: activity.Object.GetRaw(),
		},
	}
}
