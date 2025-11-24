package domain

import (
	"context"
	"database/sql"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/db"
)

type PostRepository interface {
	Repository[*api.Post, *api.NewPost, *api.UpdatePost]
	HasUserScope[*api.Post]
}

type postRepository struct {
	*db.Queries
}

var _ PostRepository = (*postRepository)(nil)

func newPostRepository(q *db.Queries) PostRepository {
	return &postRepository{q}
}

func (r *postRepository) Create(ctx context.Context, post *api.NewPost) error {
	return r.AddPost(ctx, *addPostToDB(post))
}

func (r *postRepository) Delete(ctx context.Context, id int32) error {
	return r.DeletePost(ctx, id)
}

func (r *postRepository) GetByID(ctx context.Context, id int32) (*api.Post, error) {
	p, err := r.GetPost(ctx, id)
	if err != nil {
		return nil, err
	}
	return postToAPI(&p), nil
}

func (r *postRepository) Update(ctx context.Context, post *api.UpdatePost) error {
	return r.UpdatePost(ctx, *updatePostToDB(post))
}

func (r *postRepository) GetByUserID(ctx context.Context, id int32) ([]*api.Post, error) {
	posts, err := r.GetPostsByUserID(ctx, id)
	if err != nil {
		return nil, err
	}
	var res []*api.Post
	for _, i := range posts {
		res = append(res, postToAPI(&i))
	}
	return res, nil
}

func postToAPI(p *db.Post) *api.Post {
	return &api.Post{
		CommentCount: 0,
		Content:      p.Content,
		CreatedAt:    p.CreatedAt.Time,
		Id:           int(p.ID),
		LikeCount:    int(p.LikeCount.Int32),
		ShareCount:   int(p.ShareCount.Int32),
		UpdatedAt:    p.UpdatedAt.Time,
		UserID:       int(p.UserID),
	}
}

func addPostToDB(p *api.NewPost) *db.AddPostParams {
	return &db.AddPostParams{
		UserID:  int32(p.UserID),
		Content: p.Content,
	}
}

func updatePostToDB(p *api.UpdatePost) *db.UpdatePostParams {
	return &db.UpdatePostParams{
		ID:           0,
		Content:      *p.Content,
		LikeCount:    sql.NullInt32{},
		ShareCount:   sql.NullInt32{},
		CommentCount: sql.NullInt32{},
	}
}
