package ap

import "github.com/kibirisu/borg/internal/domain"

type Objecter interface {
	ObjectOrLink[Object]
}

type Object struct{}

type object struct {
	part *domain.ObjectOrLink
}
