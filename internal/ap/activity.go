package ap

import "github.com/kibirisu/borg/internal/domain"

type Activity struct {
	ID     string
	Type   string
	Actor  Actorer
	Object Objecter
}

type activity struct {
	raw *domain.Object
}

func (a *activity) GetObject() Activity {
	return Activity{
		ID:    a.raw.ID,
		Type:  a.raw.Type,
		Actor: &actor{a.raw.ActivityActor},
	}
}
