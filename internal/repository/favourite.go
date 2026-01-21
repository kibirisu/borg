package repository

import (
	"context"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/db"
)

type FavouriteRepository interface {
	GetLocalLikeByID(context.Context, xid.ID) (db.GetLocalLikeByIDRow, error)
	Create(context.Context, db.CreateFavouriteParams) (db.CreateFavouriteRow, error)
	AddLike(context.Context, db.AddLikeParams) (db.Favourite, error)
	GetByURI(context.Context, string) (db.Favourite, error)
	DeleteByID(context.Context, xid.ID) error
	DeleteByStatusID(context.Context, xid.ID, xid.ID) (db.Favourite, error)
}

type favouriteRepository struct {
	q *db.Queries
}

var _ FavouriteRepository = (*favouriteRepository)(nil)

// GetLocalLikeByID implements FavouriteRepository.
func (r *favouriteRepository) GetLocalLikeByID(
	ctx context.Context,
	id xid.ID,
) (db.GetLocalLikeByIDRow, error) {
	return r.q.GetLocalLikeByID(ctx, id)
}

// Create implements FavouriteRepository.
func (r *favouriteRepository) Create(
	ctx context.Context,
	favourite db.CreateFavouriteParams,
) (db.CreateFavouriteRow, error) {
	return r.q.CreateFavourite(ctx, favourite)
}

// AddLike implements FavouriteRepository.
func (r *favouriteRepository) AddLike(
	ctx context.Context,
	like db.AddLikeParams,
) (db.Favourite, error) {
	return r.q.AddLike(ctx, like)
}

// GetByURI implements FavouriteRepository.
func (r *favouriteRepository) GetByURI(ctx context.Context, uri string) (db.Favourite, error) {
	return r.q.GetFavouriteByURI(ctx, uri)
}

// DeleteByID implements FavouriteRepository.
func (r *favouriteRepository) DeleteByID(ctx context.Context, id xid.ID) error {
	return r.q.DeleteFavouriteByID(ctx, id)
}

func (r *favouriteRepository) DeleteByStatusID(
	ctx context.Context,
	accountID xid.ID,
	statusID xid.ID,
) (db.Favourite, error) {
	params := db.DeleteFavouriteByStatusIDParams{
		AccountID: accountID,
		StatusID:  statusID,
	}
	return r.q.DeleteFavouriteByStatusID(ctx, params)
}
