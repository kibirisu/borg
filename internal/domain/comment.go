package domain

import (
	"context"

	"github.com/kibirisu/borg/internal/db"
)

type CommentRepository interface {
	Repository[db.Comment, db.AddCommentParams, any]
	HasUserScope[db.Comment]
	HasPostScope[db.Comment]
}

type commentRepository struct {
	*db.Queries
}

func newCommentRepository(q *db.Queries) CommentRepository {
	return &commentRepository{q}
}

func (r *commentRepository) Create(ctx context.Context, comm db.AddCommentParams) error {
	return r.AddComment(ctx, comm)
}

func (r *commentRepository) Delete(ctx context.Context, id int32) error {
	return r.DeleteComment(ctx, id)
}

func (r *commentRepository) GetByID(context.Context, int32) (db.Comment, error) {
	panic("unimplemented")
}

func (r *commentRepository) GetByPostID(ctx context.Context, id int32) ([]db.Comment, error) {
	return r.GetPostComments(ctx, id)
}

func (r *commentRepository) GetByUserID(ctx context.Context, id int32) ([]db.Comment, error) {
	return r.GetUserComments(ctx, id)
}

func (r *commentRepository) Update(context.Context, any) error {
	panic("unimplemented")
}
