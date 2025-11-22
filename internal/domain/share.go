package domain

import (
	"context"

	"borg/internal/db"
)

type ShareRepository interface {
	Repository[db.Share, db.AddShareParams, any]
	HasUserScope[db.Share]
	HasPostScope[db.Share]
}

type shareRepository struct {
	*db.Queries
}

func newShareRepository(q *db.Queries) ShareRepository {
	return &shareRepository{q}
}

func (r *shareRepository) Create(ctx context.Context, share db.AddShareParams) error {
	return r.AddShare(ctx, share)
}

func (r *shareRepository) Delete(ctx context.Context, id int32) error {
	return r.DeleteShare(ctx, id)
}

func (r *shareRepository) GetByID(ctx context.Context, id int32) (db.Share, error) {
	return r.GetShareByID(ctx, id)
}

func (r *shareRepository) Update(context.Context, any) error {
	panic("unimplemented")
}

func (r *shareRepository) GetByUserID(ctx context.Context, id int32) ([]db.Share, error) {
	return r.GetShareByUserID(ctx, id)
}

func (r *shareRepository) GetByPostID(ctx context.Context, id int32) ([]db.Share, error) {
	return r.GetSharesByPostID(ctx, id)
}
