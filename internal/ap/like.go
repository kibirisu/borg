package ap

import "github.com/kibirisu/borg/internal/domain"

type LikeActivitier interface {
	Activiter[Note]
}

type likeActivity struct {
	activity
}

var _ LikeActivitier = (*likeActivity)(nil)

func NewLikeActivity(from *domain.ObjectOrLink) LikeActivitier {
	return &likeActivity{activity{object{from}}}
}

func NewEmptyLikeActivity() LikeActivitier {
	return &likeActivity{activity{object{}}}
}

// GetObject implements likeActivitier.
// Subtle: this method shadows the method (activity).GetObject of likeActivity.activity.
func (f *likeActivity) GetObject() Activity[Note] {
	return Activity[Note]{
		ID:     f.raw.Object.ID,
		Type:   f.raw.Object.Type,
		Actor:  &actor{object{f.raw.Object.ActivityActor}},
		Object: &note{object{f.raw.Object.ActivityObject}},
	}
}

// SetObject implements LikeActivitier.
// Subtle: this method shadows the method (activity).SetObject of likeActivity.activity.
func (f *likeActivity) SetObject(activity Activity[Note]) {
	f.raw = &domain.ObjectOrLink{
		Object: &domain.Object{
			ID:             activity.ID,
			Type:           activity.Type,
			ActivityActor:  activity.Actor.GetRaw(),
			ActivityObject: activity.Object.GetRaw(),
		},
	}
}

// WithLink implements LikeActivitier.
// Subtle: this method shadows the method (activity).WithLink of likeActivity.activity.
func (f *likeActivity) WithLink(link string) Objecter[Activity[Note]] {
	f.SetLink(link)
	return f
}

// WithObject implements LikeActivitier.
// Subtle: this method shadows the method (activity).WithObject of likeActivity.activity.
func (f *likeActivity) WithObject(activity Activity[Note]) Objecter[Activity[Note]] {
	f.SetObject(activity)
	return f
}
