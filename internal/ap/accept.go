package ap

import "github.com/kibirisu/borg/internal/domain"

type AcceptActivitier interface {
	Activiter[Activity[Actor]]
}

type acceptActivity struct {
	activity
}

var _ AcceptActivitier = (*acceptActivity)(nil)

func NewAcceptActivity(from *domain.ObjectOrLink) AcceptActivitier {
	return &acceptActivity{activity{object{from}}}
}

// GetObject implements AcceptActivitier.
// Subtle: this method shadows the method (activity).GetObject of acceptActivity.activity.
func (a *acceptActivity) GetObject() Activity[Activity[Actor]] {
	return Activity[Activity[Actor]]{
		ID:     a.raw.Object.ID,
		Type:   a.raw.Object.Type,
		Actor:  &actor{object{a.raw.Object.ActivityActor}},
		Object: &followActivity{activity{object{a.raw.Object.ActivityObject}}},
	}
}

// SetObject implements AcceptActivitier.
// Subtle: this method shadows the method (activity).SetObject of acceptActivity.activity.
func (a *acceptActivity) SetObject(activity Activity[Activity[Actor]]) {
	a.raw = &domain.ObjectOrLink{
		Object: &domain.Object{
			ID:             "",
			Type:           "",
			ActivityActor:  activity.Actor.GetRaw(),
			ActivityObject: activity.Object.GetRaw(),
		},
	}
}
