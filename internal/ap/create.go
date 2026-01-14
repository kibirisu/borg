package ap

import "github.com/kibirisu/borg/internal/domain"

type CreateActivitier interface {
	Activiter[Note]
}

type createActivity struct {
	activity
}

var _ CreateActivitier = (*createActivity)(nil)

func NewCreateActivity(from *domain.ObjectOrLink) CreateActivitier {
	return &createActivity{activity{object{from}}}
}

func NewEmptyCreateActivity() CreateActivitier {
	return &createActivity{activity{object{}}}
}

// GetObject implements CreateActivitier.
// Subtle: this method shadows the method (activity).GetObject of createActivity.activity.
func (c *createActivity) GetObject() Activity[Note] {
	return Activity[Note]{
		ID:     c.raw.Object.ID,
		Type:   c.raw.Object.Type,
		Actor:  &actor{object{c.raw.Object.ActivityActor}},
		Object: &note{object{c.raw.Object.ActivityObject}},
	}
}

// SetObject implements CreateActivitier.
// Subtle: this method shadows the method (activity).SetObject of createActivity.activity.
func (c *createActivity) SetObject(activity Activity[Note]) {
	c.raw = &domain.ObjectOrLink{
		Object: &domain.Object{
			ID:             activity.ID,
			Type:           activity.Type,
			ActivityActor:  activity.Actor.GetRaw(),
			ActivityObject: activity.Object.GetRaw(),
		},
	}
}

// WithLink implements CreateActivitier.
// Subtle: this method shadows the method (activity).WithLink of createActivity.activity.
func (c *createActivity) WithLink(link string) Objecter[Activity[Note]] {
	c.SetLink(link)
	return c
}

// WithObject implements CreateActivitier.
// Subtle: this method shadows the method (activity).WithObject of createActivity.activity.
func (c *createActivity) WithObject(activity Activity[Note]) Objecter[Activity[Note]] {
	c.SetObject(activity)
	return c
}
