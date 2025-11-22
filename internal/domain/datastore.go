package domain

import (
	"context"

	"borg/internal/db"
)

type DataStore interface {
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
	users    UserRepository
	posts    PostRepository
	comments CommentRepository
	likes    LikeRepository
	shares   ShareRepository
}

var _ DataStore = (*dataStore)(nil)

func NewDataStore(ctx context.Context, url string) (DataStore, error) {
	q, err := db.GetDB(ctx, url)
	if err != nil {
		return nil, err
	}
	ds := &dataStore{}
	ds.users = newUserRepository(q)
	ds.posts = newPostRepository(q)
	ds.comments = newCommentRepository(q)
	ds.likes = newLikeRepository(q)
	ds.shares = newShareRepository(q)
	return ds, nil
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
