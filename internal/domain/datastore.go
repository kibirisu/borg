package domain

import (
	"context"

	"github.com/kibirisu/borg/internal/db"
)

type DataStore interface {
	Raw() *db.Queries
	UserRepository() UserRepository
	PostRepository() PostRepository
	CommentRepository() CommentRepository
	LikeRepository() LikeRepository
	ShareRepository() ShareRepository
}

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

type dataStore struct {
	raw      *db.Queries
	users    UserRepository
	posts    PostRepository
	comments CommentRepository
	likes    LikeRepository
	shares   ShareRepository
}

var _ DataStore = (*dataStore)(nil)

func NewDataStore(ctx context.Context, url string) DataStore {
	q := db.GetDB(ctx, url)
	ds := &dataStore{}
	ds.users = newUserRepository(q)
	ds.posts = newPostRepository(q)
	ds.comments = newCommentRepository(q)
	ds.likes = newLikeRepository(q)
	ds.shares = newShareRepository(q)
	ds.raw = q
	return ds
}

func (ds *dataStore) Raw() *db.Queries {
	return ds.raw
}

func (ds *dataStore) UserRepository() UserRepository {
	return ds.users
}

func (ds *dataStore) PostRepository() PostRepository {
	return ds.posts
}

func (ds *dataStore) CommentRepository() CommentRepository {
	return ds.comments
}

func (ds *dataStore) LikeRepository() LikeRepository {
	return ds.likes
}

func (ds *dataStore) ShareRepository() ShareRepository {
	return ds.shares
}
