package ap

import "github.com/kibirisu/borg/internal/domain"

type AnnounceActivitier interface {
	Activiter[Note]
}

type announceActivity struct {
	activity
}

var _ AnnounceActivitier = (*announceActivity)(nil)

func NewAnnounceActivity(from *domain.ObjectOrLink) AnnounceActivitier {
	return &announceActivity{activity{object{from}}}
}

// GetObject implements AnnounceActivitier.
// Subtle: this method shadows the method (activity).GetObject of announceActivity.activity.
func (a *announceActivity) GetObject() Activity[Note] {
	return Activity[Note]{
		ID:     a.raw.Object.ID,
		Type:   a.raw.Object.Type,
		Actor:  &actor{object{a.raw.Object.ActivityActor}},
		Object: &note{object{a.raw.Object.ActivityObject}},
	}
}

// SetObject implements AnnounceActivitier.
// Subtle: this method shadows the method (activity).SetObject of announceActivity.activity.
func (a *announceActivity) SetObject(activity Activity[Note]) {
	a.raw = &domain.ObjectOrLink{
		Object: &domain.Object{
			ID:             activity.ID,
			Type:           activity.Type,
			ActivityActor:  activity.Actor.GetRaw(),
			ActivityObject: activity.Object.GetRaw(),
		},
	}
}
