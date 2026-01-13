package repository

import (
	"context"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/db"
)

type FavouriteRepository interface {
	Create(context.Context, db.CreateFavouriteParams) (db.Favourite, error)
	GetByURI(context.Context, string) (db.Favourite, error)
	GetByPost(context.Context, xid.ID) ([]db.Favourite, error)
	DeleteByID(context.Context, xid.ID) error
	GetLikedPostsByAccountId(context.Context, xid.ID) ([]db.GetLikedPostsByAccountIdRow, error)
}

type favouriteRepository struct {
	q *db.Queries
}

var _ FavouriteRepository = (*favouriteRepository)(nil)

func NewFavouriteRepository(q *db.Queries) FavouriteRepository {
	return &favouriteRepository{q: q}
}

// Create implements FavouriteRepository.
func (r *favouriteRepository) Create(
	ctx context.Context,
	params db.CreateFavouriteParams,
) (db.Favourite, error) {
	return r.q.CreateFavourite(ctx, params)
}

// GetByPost implements FavouriteRepository.
func (r *favouriteRepository) GetByPost(
	ctx context.Context,
	id xid.ID,
) ([]db.Favourite, error) {
	return r.q.GetStatusFavourites(ctx, id)
}

// GetByURI implements FavouriteRepository.
func (r *favouriteRepository) GetByURI(ctx context.Context, uri string) (db.Favourite, error) {
	return r.q.GetFavouriteByURI(ctx, uri)
}

// DeleteByID implements FavouriteRepository.
func (r *favouriteRepository) DeleteByID(ctx context.Context, id xid.ID) error {
	return r.q.DeleteFavouriteByID(ctx, id)
}

// GetLikedPostsByAccountId implements FavouriteRepository.
func (r *favouriteRepository) GetLikedPostsByAccountId(
	ctx context.Context,
	accountID xid.ID,
) ([]db.GetLikedPostsByAccountIdRow, error) {
	return r.q.GetLikedPostsByAccountId(ctx, accountID)
}

