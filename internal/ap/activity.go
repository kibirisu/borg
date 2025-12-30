package ap

type Activiter interface {
	Objecter[Activity]
}

type Activity struct {
	ID     string
	Type   string
	Actor  Actorer
	Object Objecter[any]
}

type activity struct {
	object
}

var _ Activiter = (*activity)(nil)

// GetObject implements Activiter.
// Subtle: this method shadows the method (object).GetObject of activity.object.
func (a *activity) GetObject() Activity {
	return Activity{
		ID:     a.raw.Object.ID,
		Type:   a.raw.Object.Type,
		Actor:  &actor{object{a.raw.Object.ActivityActor}},
		Object: &object{a.raw.Object.ActivityObject},
	}
}

// SetObject implements Activiter.
// Subtle: this method shadows the method (object).SetObject of activity.object.
func (a *activity) SetObject(activity Activity) {
	a.raw.Object.ActivityActor = activity.Actor.GetRaw()
	a.raw.Object.ActivityObject = activity.Object.GetRaw()
}
