package mapper

import (
	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/db"
)

func PostToAPIWithMetadata(
	post *db.Status,
	acc *db.Account,
	LikeCount int,
	ShareCount int,
	CommentCount int,
) *api.Post {
	return &api.Post{
		CommentCount: CommentCount,
		Content:      post.Content,
		CreatedAt:    post.CreatedAt,
		Id:           post.ID.String(),
		LikeCount:    LikeCount,
		ShareCount:   ShareCount,
		UpdatedAt:    post.UpdatedAt,
		UserID:       post.AccountID.String(),
		Username:     &acc.Username,
	}
}
