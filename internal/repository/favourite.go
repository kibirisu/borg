package repository

import (
	"context"

	"github.com/kibirisu/borg/internal/db"
)

type FavouriteRepository interface {
	Create(context.Context, db.CreateFavouriteParams) (db.Favourite, error)
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

// GetLikedPostsByUser implements FavouriteRepository.
func (r *favouriteRepository) GetLikedPostsByUser(
	ctx context.Context,
	accountID int,
) ([]db.GetLikedPostsByAccountIdRow, error) {
	return r.q.GetLikedPostsByAccountId(ctx, int32(accountID))
}