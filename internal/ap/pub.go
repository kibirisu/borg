package ap

import "github.com/kibirisu/borg/internal/domain"

type ValueType string

const (
	ObjectType  ValueType = "object"
	LinkType    ValueType = "link"
	NullType    ValueType = "null"
	InvalidType ValueType = "invalid"
)

type ObjectOrLink[T any] interface {
	GetObject() T
	GetLink() string
	SetObject(T)
	SetLink(string)
	SetNull()
	GetValueType() ValueType
	GetURI() string
	GetRaw() *domain.ObjectOrLink
}
