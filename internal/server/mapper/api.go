package mapper

import (
	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/db"
)

func AccountToAPI(account *db.Account) *api.Account {
	return &api.Account{
		Acct:        "", // TODO
		DisplayName: account.DisplayName.String,
		Id:          int(account.ID),
		Url:         account.Url,
		Username:    account.Username,
	}
}
func PostToAPI(post *db.Status) *api.Post {
	return &api.Post{
		CommentCount: -1,
		Content: post.Content,
		CreatedAt: post.CreatedAt,
		Id: int(post.ID),
		LikeCount: -1,
		ShareCount: -1,
		UpdatedAt: post.UpdatedAt,
		UserID: int(post.AccountID),
		Username: nil,
	}
}
func PostToAPIWithMetadata(post *db.GetStatusByIdWithMetadataRow) *api.Post {
	return &api.Post{
		CommentCount: int(post.CommentCount),
		Content: post.Content,
		CreatedAt: post.CreatedAt,
		Id: int(post.ID),
		LikeCount: int(post.LikeCount),
		ShareCount: int(post.ShareCount),
		UpdatedAt: post.UpdatedAt,
		UserID: int(post.AccountID),
		Username: &post.OwnerUsername,
	}
}
func LikeToAPI(like *db.Favourite) *api.Like {
	return &api.Like{
		CreatedAt: like.CreatedAt,
		Id: int(like.ID),
		PostID: int(like.StatusID),
		UserID: int(like.AccountID),
	}
}
