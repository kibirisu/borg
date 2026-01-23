package ap

import "github.com/kibirisu/borg/internal/domain"

type UndoActiviter interface {
	Activiter[Activity[Note]]
}

type undoActivity struct {
	activity
}

var _ UndoActiviter = (*undoActivity)(nil)

func NewUndoActivity(from *domain.ObjectOrLink) UndoActiviter {
	return &undoActivity{activity{object{from}}}
}

func NewEmptyUndoActivity() UndoActiviter {
	return &undoActivity{activity{object{}}}
}

// GetObject implements UndoActiviter.
// Subtle: this method shadows the method (activity).GetObject of undoActivity.activity.
func (u *undoActivity) GetObject() Activity[Activity[Note]] {
	return Activity[Activity[Note]]{
		ID:     u.raw.Object.ID,
		Type:   u.raw.Object.Type,
		Actor:  &actor{object{u.raw.Object.ActivityActor}},
		Object: &likeActivity{activity{object{u.raw.Object.ActivityObject}}},
	}
}

// SetObject implements UndoActiviter.
// Subtle: this method shadows the method (activity).SetObject of undoActivity.activity.
func (u *undoActivity) SetObject(activity Activity[Activity[Note]]) {
	u.raw = &domain.ObjectOrLink{
		Object: &domain.Object{
			ID:             activity.ID,
			Type:           activity.Type,
			ActivityActor:  activity.Actor.GetRaw(),
			ActivityObject: activity.Object.GetRaw(),
		},
	}
}

// WithLink implements UndoActiviter.
// Subtle: this method shadows the method (activity).WithLink of undoActivity.activity.
func (u *undoActivity) WithLink(link string) Objecter[Activity[Activity[Note]]] {
	u.SetLink(link)
	return u
}

// WithObject implements UndoActiviter.
// Subtle: this method shadows the method (activity).WithObject of undoActivity.activity.
func (u *undoActivity) WithObject(
	activity Activity[Activity[Note]],
) Objecter[Activity[Activity[Note]]] {
	u.SetObject(activity)
	return u
}
