package repository

import (
	"context"

	"github.com/kibirisu/borg/internal/db"
)

type FavouriteRepository interface {
	Create(context.Context, db.CreateFavouriteParams) (db.Favourite, error)
	GetByURI(context.Context, string) (db.Favourite, error)
	GetByPost(context.Context, int) ([]db.Favourite, error)
	DeleteByID(context.Context, int32) error
	GetByPost(ctx context.Context, id int) ([]db.Favourite, error)
	GetLikedPostsByUser(ctx context.Context, accountID int) ([]db.GetLikedPostsByAccountIdRow, error)
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
	id int,
) ([]db.Favourite, error) {
	return r.q.GetStatusFavourites(ctx, int32(id))
}

// GetByURI implements FavouriteRepository.
func (r *favouriteRepository) GetByURI(ctx context.Context, uri string) (db.Favourite, error) {
	return r.q.GetFavouriteByURI(ctx, uri)
}

// DeleteByID implements FavouriteRepository.
func (r *favouriteRepository) DeleteByID(ctx context.Context, id int32) error {
	return r.q.DeleteFavouriteByID(ctx, id)
}

// GetByURI implements FavouriteRepository.
func (r *favouriteRepository) GetByURI(ctx context.Context, uri string) (db.Favourite, error) {
	return r.q.GetFavouriteByURI(ctx, uri)
}

// DeleteByID implements FavouriteRepository.
func (r *favouriteRepository) DeleteByID(ctx context.Context, id int32) error {
	return r.q.DeleteFavouriteByID(ctx, id)
}

// GetLikedPostsByUser implements FavouriteRepository.
func (r *favouriteRepository) GetLikedPostsByUser(
	ctx context.Context,
	accountID int,
) ([]db.GetLikedPostsByAccountIdRow, error) {
	return r.q.GetLikedPostsByAccountId(ctx, int32(accountID))
}