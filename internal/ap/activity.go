package ap

import "github.com/kibirisu/borg/internal/domain"

type Activiter[T any] interface {
	Objecter[Activity[T]]
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

var _ Activiter[any] = (*activity)(nil)

func NewActivity(from *domain.ObjectOrLink) Activiter[any] {
	return &activity{object{from}}
}

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

// WithLink implements Activiter.
// Subtle: this method shadows the method (activity).WithLink of activity.object.
func (a *activity) WithLink(link string) Objecter[Activity[any]] {
	a.SetLink(link)
	return a
}

// WithObject implements Activiter.
// Subtle: this method shadows the method (activity).WithObject of activity.object.
func (a *activity) WithObject(
	activity Activity[any],
) Objecter[Activity[any]] {
	a.SetObject(activity)
	return a
}
