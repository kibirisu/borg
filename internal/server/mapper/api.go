package mapper

import (
	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/db"
)

func AccountToAPI(account *db.Account) *api.Account {
	return &api.Account{
		Acct:        "", // TODO
		DisplayName: account.DisplayName.String,
		Id:          account.ID.String(),
		Url:         account.Url,
		Username:    account.Username,
	}
}

func PostToAPI(post *db.Status) *api.Post {
	return &api.Post{
		CommentCount: -1,
		Content:      post.Content,
		CreatedAt:    post.CreatedAt,
		Id:           post.ID.String(),
		LikeCount:    -1,
		ShareCount:   -1,
		UpdatedAt:    post.UpdatedAt,
		UserID:       post.AccountID.String(),
		Username:     nil,
	}
}

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

func LikeToAPI(like *db.Favourite) *api.Like {
	return &api.Like{
		CreatedAt: like.CreatedAt,
		Id:        like.ID.String(),
		PostID:    like.StatusID.String(),
		UserID:    like.AccountID.String(),
	}
}

func AccountToUserAPI(account *db.Account, followersCount int, followingCount int) *api.User {
	origin := "local"
	if account.Domain.Valid && account.Domain.String != "" {
		origin = account.Domain.String
	}
	return &api.User{
		Id:             account.ID.String(),
		Username:       account.Username,
		Bio:            account.DisplayName.String,
		Origin:         origin,
		IsAdmin:        false, // TODO: add admin flag to accounts table if needed
		FollowersCount: followersCount,
		FollowingCount: followingCount,
		CreatedAt:      account.CreatedAt,
		UpdatedAt:      account.UpdatedAt,
	}
}
