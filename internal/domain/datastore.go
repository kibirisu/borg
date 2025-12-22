package domain

import "context"

type Repository[T, C, U any] interface {
	Create(context.Context, C) error
	GetByID(context.Context, int32) (T, error)
	Update(context.Context, U) error
	Delete(context.Context, int32) error
}

type UserScopedRepository[T, C, U any] interface {
	Repository[T, C, U]
	HasUserScope[T]
}

type PostScopedRepository[T, C, U any] interface {
	Repository[T, C, U]
	HasPostScope[T]
}

type HasUserScope[T any] interface {
	GetByUserID(context.Context, int32) ([]T, error)
}

type HasPostScope[T any] interface {
	GetByPostID(context.Context, int32) ([]T, error)
}
