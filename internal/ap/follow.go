package ap

import "github.com/kibirisu/borg/internal/domain"

type FollowActivitier interface {
	Activiter[Actor]
}

type UndoFollowActiviter interface {
	Activiter[Activity[Actor]]
}

type followActivity struct {
	activity
}

type undoFollowActivity struct {
	activity
}

var (
	_ FollowActivitier    = (*followActivity)(nil)
	_ UndoFollowActiviter = (*undoFollowActivity)(nil)
)

func NewFollowActivity(from *domain.ObjectOrLink) FollowActivitier {
	return &followActivity{activity{object{from}}}
}

func NewEmptyFollowActivity() FollowActivitier {
	return &followActivity{activity{object{}}}
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

// WithLink implements FollowActivitier.
// Subtle: this method shadows the method (activity).WithLink of followActivity.activity.
func (f *followActivity) WithLink(link string) Objecter[Activity[Actor]] {
	f.SetLink(link)
	return f
}

// WithObject implements FollowActivitier.
// Subtle: this method shadows the method (activity).WithObject of followActivity.activity.
func (f *followActivity) WithObject(activity Activity[Actor]) Objecter[Activity[Actor]] {
	f.SetObject(activity)
	return f
}

// GetObject implements UndoFollowActiviter.
// Subtle: this method shadows the method (activity).GetObject of undoFollowActivity.activity.
func (u *undoFollowActivity) GetObject() Activity[Activity[Actor]] {
	return Activity[Activity[Actor]]{
		ID:     u.raw.Object.ID,
		Type:   u.raw.Object.Type,
		Actor:  &actor{object{u.raw.Object.ActivityActor}},
		Object: &followActivity{activity{object{u.raw.Object.ActivityObject}}},
	}
}

// SetObject implements UndoFollowActiviter.
// Subtle: this method shadows the method (activity).SetObject of undoFollowActivity.activity.
func (u *undoFollowActivity) SetObject(activity Activity[Activity[Actor]]) {
	u.raw = &domain.ObjectOrLink{
		Object: &domain.Object{
			ID:             activity.ID,
			Type:           activity.Type,
			ActivityActor:  activity.Actor.GetRaw(),
			ActivityObject: activity.Object.GetRaw(),
		},
	}
}

// WithLink implements UndoFollowActiviter.
// Subtle: this method shadows the method (activity).WithLink of undoFollowActivity.activity.
func (u *undoFollowActivity) WithLink(link string) Objecter[Activity[Activity[Actor]]] {
	u.SetLink(link)
	return u
}

// WithObject implements UndoFollowActiviter.
// Subtle: this method shadows the method (activity).WithObject of undoFollowActivity.activity.
func (u *undoFollowActivity) WithObject(
	activity Activity[Activity[Actor]],
) Objecter[Activity[Activity[Actor]]] {
	u.SetObject(activity)
	return u
}
