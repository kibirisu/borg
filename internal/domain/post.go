package domain

import (
	"context"

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

func newPostRepository(q *db.Queries) PostRepository {
	return &postRepository{q}
}

func (r *postRepository) Create(ctx context.Context, post *api.NewPost) error {
	panic("unimplemented")
}

func (r *postRepository) Delete(ctx context.Context, id int32) error {
	panic("unimplemented")
}

func (r *postRepository) GetByID(ctx context.Context, id int32) (*api.Post, error) {
	panic("unimplemented")
}

func (r *postRepository) Update(ctx context.Context, post *api.UpdatePost) error {
	panic("unimplemented")
}

func (r *postRepository) GetByUserID(ctx context.Context, id int32) ([]*api.Post, error) {
	panic("unimplemented")
}

// func postToAPI(p *db.Post) *api.Post {
// 	return &api.Post{
// 		CommentCount: int(p.CommentCount.Int32),
// 		Content:      p.Content,
// 		CreatedAt:    p.CreatedAt.Time,
// 		Id:           int(p.ID),
// 		LikeCount:    int(p.LikeCount.Int32),
// 		ShareCount:   int(p.ShareCount.Int32),
// 		UpdatedAt:    p.UpdatedAt.Time,
// 		UserID:       int(p.UserID),
// 	}
// }
//
// func addPostToDB(p *api.NewPost) *db.AddPostParams {
// 	return &db.AddPostParams{
// 		UserID:  int32(p.UserID),
// 		Content: p.Content,
// 	}
// }
//
// func updatePostToDB(p *api.UpdatePost) *db.UpdatePostParams {
// 	return &db.UpdatePostParams{
// 		ID:           0,
// 		Content:      *p.Content,
// 		LikeCount:    sql.NullInt32{},
// 		ShareCount:   sql.NullInt32{},
// 		CommentCount: sql.NullInt32{},
// 	}
// }
//
// func (r *postRepository) GetAll(ctx context.Context) ([]*api.Post, error) {
// 	posts, err := r.GetAllPosts(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var res []*api.Post
// 	for _, p := range posts {
// 		res = append(res, postRowToAPI(&p))
// 	}
// 	return res, nil
// }
//
// func postRowToAPI(p *db.GetAllPostsRow) *api.Post {
// 	likeCount := 0
// 	if p.LikeCount.Valid {
// 		likeCount = int(p.LikeCount.Int32)
// 	}
// 	shareCount := 0
// 	if p.ShareCount.Valid {
// 		shareCount = int(p.ShareCount.Int32)
// 	}
// 	commentCount := 0
// 	if p.CommentCount.Valid {
// 		commentCount = int(p.CommentCount.Int32)
// 	}
// 	username := p.Username
// 	return &api.Post{
// 		CommentCount: commentCount,
// 		Content:      p.Content,
// 		CreatedAt:    p.CreatedAt.Time,
// 		Id:           int(p.ID),
// 		LikeCount:    likeCount,
// 		ShareCount:   shareCount,
// 		UpdatedAt:    p.UpdatedAt.Time,
// 		UserID:       int(p.UserID),
// 		Username:     &username,
// 	}
// }
