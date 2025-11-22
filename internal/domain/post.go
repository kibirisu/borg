package domain

import (
	"context"
	"database/sql"

	"borg/internal/api"
	"borg/internal/db"
)

type PostRepository interface {
	Repository[*api.Post, *api.NewPost, *api.UpdatePost]
	HasUserScope[*api.Post]
	GetAll(context.Context) ([]*api.Post, error)
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

func (r *postRepository) GetAll(ctx context.Context) ([]*api.Post, error) {
	rows, err := r.GetAllPosts(ctx)
	if err != nil {
		return nil, err
	}
	var res []*api.Post
	for _, row := range rows {
		res = append(res, postRowToAPI(&row))
	}
	return res, nil
}

func postToAPI(p *db.Post) *api.Post {
	likeCount := 0
	if p.LikeCount.Valid {
		likeCount = int(p.LikeCount.Int32)
	}
	shareCount := 0
	if p.ShareCount.Valid {
		shareCount = int(p.ShareCount.Int32)
	}
	commentCount := 0
	if p.CommentCount.Valid {
		commentCount = int(p.CommentCount.Int32)
	}
	return &api.Post{
		CommentCount: commentCount,
		Content:      p.Content,
		CreatedAt:    p.CreatedAt.Time,
		Id:           int(p.ID),
		LikeCount:    likeCount,
		ShareCount:   shareCount,
		UpdatedAt:    p.UpdatedAt.Time,
		UserID:       int(p.UserID),
	}
}

func postRowToAPI(row *db.GetAllPostsRow) *api.Post {
	likeCount := 0
	if row.LikeCount.Valid {
		likeCount = int(row.LikeCount.Int32)
	}
	shareCount := 0
	if row.ShareCount.Valid {
		shareCount = int(row.ShareCount.Int32)
	}
	commentCount := 0
	if row.CommentCount.Valid {
		commentCount = int(row.CommentCount.Int32)
	}
	username := row.Username
	return &api.Post{
		CommentCount: commentCount,
		Content:      row.Content,
		CreatedAt:    row.CreatedAt.Time,
		Id:           int(row.ID),
		LikeCount:    likeCount,
		ShareCount:   shareCount,
		UpdatedAt:    row.UpdatedAt.Time,
		UserID:       int(row.UserID),
		Username:     &username,
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
