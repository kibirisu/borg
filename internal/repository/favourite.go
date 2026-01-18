package repository

import (
	"context"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/db"
)

type FavouriteRepository interface {
	Create(context.Context, db.CreateFavouriteParams) (db.Favourite, error)
	CreateNew(context.Context, db.CreateFavouriteNewParams) (db.Favourite, error)
	GetByURI(context.Context, string) (db.Favourite, error)
	DeleteByID(context.Context, xid.ID) error
	DeleteByIDNew(context.Context, xid.ID) (db.Favourite, error)
	GetLikedPostsByAccountID(context.Context, xid.ID) ([]db.GetLikedPostsByAccountIdRow, error)
}

type favouriteRepository struct {
	q *db.Queries
}

var _ FavouriteRepository = (*favouriteRepository)(nil)

// Create implements FavouriteRepository.
func (r *favouriteRepository) Create(
	ctx context.Context,
	params db.CreateFavouriteParams,
) (db.Favourite, error) {
	return r.q.CreateFavourite(ctx, params)
}

// CreateNew implements FavouriteRepository.
func (r *favouriteRepository) CreateNew(
	ctx context.Context,
	favourite db.CreateFavouriteNewParams,
) (db.Favourite, error) {
	return r.q.CreateFavouriteNew(ctx, favourite)
}

// GetByURI implements FavouriteRepository.
func (r *favouriteRepository) GetByURI(ctx context.Context, uri string) (db.Favourite, error) {
	return r.q.GetFavouriteByURI(ctx, uri)
}

// DeleteByID implements FavouriteRepository.
func (r *favouriteRepository) DeleteByID(ctx context.Context, id xid.ID) error {
	return r.q.DeleteFavouriteByID(ctx, id)
}

// DeleteByIDNew implements FavouriteRepository.
func (r *favouriteRepository) DeleteByIDNew(ctx context.Context, id xid.ID) (db.Favourite, error) {
	return r.q.DeleteFavouriteByIDNew(ctx, id)
}

// GetLikedPostsByAccountId implements FavouriteRepository.
func (r *favouriteRepository) GetLikedPostsByAccountID(
	ctx context.Context,
	accountID xid.ID,
) ([]db.GetLikedPostsByAccountIdRow, error) {
	return r.q.GetLikedPostsByAccountId(ctx, accountID)
}
