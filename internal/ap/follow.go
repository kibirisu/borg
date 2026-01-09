package ap

import "github.com/kibirisu/borg/internal/domain"

type FollowActivitier interface {
	Activiter[Actor]
}

type followActivity struct {
	activity
}

var _ FollowActivitier = (*followActivity)(nil)

func NewFollowActivity(from *domain.ObjectOrLink) FollowActivitier {
	return &followActivity{activity{object{from}}}
}

// GetObject implements FollowActivitier.
// Subtle: this method shadows the method (activity).GetObject of followActivity.activity.
func (f *followActivity) GetObject() Activity[Actor] {
	return Activity[Actor]{
		ID:     f.raw.Object.ID,
		Type:   f.raw.Object.Type,
		Actor:  &actor{object{f.raw.Object.ActivityActor}},
		Object: &actor{object{f.raw.Object.ActivityObject}},
	}
}

// SetObject implements FollowActivitier.
// Subtle: this method shadows the method (activity).SetObject of followActivity.activity.
func (f *followActivity) SetObject(activity Activity[Actor]) {
	f.raw = &domain.ObjectOrLink{
		Object: &domain.Object{
			ID:             activity.ID,
			Type:           activity.Type,
			ActivityActor:  activity.Actor.GetRaw(),
			ActivityObject: activity.Object.GetRaw(),
		},
	}
}
