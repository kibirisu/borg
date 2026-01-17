package mapper

import (
	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/db"
)

func PostToAPIWithMetadata(
	post *db.Status,
	acc *db.Account,
	likeCount int,
	shareCount int,
	commentCount int,
) *api.Post {
	return &api.Post{
		CommentCount: commentCount,
		Content:      post.Content,
		CreatedAt:    post.CreatedAt,
		Id:           post.ID.String(),
		LikeCount:    likeCount,
		ShareCount:   shareCount,
		UpdatedAt:    post.UpdatedAt,
		UserID:       post.AccountID.String(),
		Username:     &acc.Username,
	}
}
