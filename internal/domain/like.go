package domain

import (
	"context"

	"borg/internal/db"
)

type LikeRepository interface {
	Repository[db.Like, db.AddLikeParams, any]
	HasUserScope[db.Like]
	HasPostScope[db.Like]
}

type likeRepository struct {
	*db.Queries
}

func newLikeRepository(q *db.Queries) LikeRepository {
	return &likeRepository{q}
}

func (r *likeRepository) Create(ctx context.Context, like db.AddLikeParams) error {
	return r.AddLike(ctx, like)
}

func (r *likeRepository) Delete(ctx context.Context, id int32) error {
	return r.DeleteLike(ctx, id)
}

func (r *likeRepository) GetByID(ctx context.Context, id int32) (db.Like, error) {
	return r.GetLikeByID(ctx, id)
}

func (r *likeRepository) Update(context.Context, any) error {
	panic("unimplemented")
}

func (r *likeRepository) GetByUserID(ctx context.Context, id int32) ([]db.Like, error) {
	return r.GetLikesByUserID(ctx, id)
}

func (r *likeRepository) GetByPostID(ctx context.Context, id int32) ([]db.Like, error) {
	return r.GetLikesByPostID(ctx, id)
}
