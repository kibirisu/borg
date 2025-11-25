package domain

import (
	"context"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/db"
)

type CommentRepository interface {
	Repository[*api.Comment, *api.NewComment, any]
	HasUserScope[*api.Comment]
	HasPostScope[*api.Comment]
}

type commentRepository struct {
	*db.Queries
}

func newCommentRepository(q *db.Queries) CommentRepository {
	return &commentRepository{q}
}

func (r *commentRepository) Create(ctx context.Context, comm *api.NewComment) error {
	return r.AddComment(ctx, *newCommentToDB(comm))
}

func (r *commentRepository) Delete(ctx context.Context, id int32) error {
	return r.DeleteComment(ctx, id)
}

func (r *commentRepository) GetByID(context.Context, int32) (*api.Comment, error) {
	panic("unimplemented")
}

func (r *commentRepository) GetByPostID(ctx context.Context, id int32) ([]*api.Comment, error) {
	comments, err := r.GetPostComments(ctx, id)
	if err != nil {
		return nil, err
	}
	var res []*api.Comment
	for _, i := range comments {
		res = append(res, commentToAPI(&i))
	}
	return res, nil
}

func (r *commentRepository) GetByUserID(ctx context.Context, id int32) ([]*api.Comment, error) {
	comments, err := r.GetUserComments(ctx, id)
	if err != nil {
		return nil, err
	}
	var res []*api.Comment
	for _, i := range comments {
		res = append(res, commentToAPI(&i))
	}
	return res, nil
}
func commentToAPI(p *db.Comment) *api.Comment {
	return &api.Comment{
		UpdatedAt: p.UpdatedAt.Time,
		Content:   p.Content,
		CreatedAt: p.UpdatedAt.Time,
		Id:        int(p.ID),
		ParentID:  int(p.ParentID.Int32),
		PostID:    int(p.PostID),
		UserID:    int(p.UserID),
	}
}
func newCommentToDB(p *api.NewComment) *db.AddCommentParams {
	return &db.AddCommentParams{
		PostID:  int32(p.PostID),
		UserID:  int32(p.UserID),
		Content: p.Content,
	}
}

func (r *commentRepository) Update(context.Context, any) error {
	panic("unimplemented")
}
