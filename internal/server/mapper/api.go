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
		Content:      post.Content,
		CreatedAt:    post.CreatedAt,
		Id:           int(post.ID),
		LikeCount:    -1,
		ShareCount:   -1,
		UpdatedAt:    post.UpdatedAt,
		UserID:       int(post.AccountID),
		Username:     nil,
	}
}

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
		Id:           int(post.ID),
		LikeCount:    likeCount,
		ShareCount:   shareCount,
		UpdatedAt:    post.UpdatedAt,
		UserID:       int(post.AccountID),
		Username:     &acc.Username,
	}
}

func LikeToAPI(like *db.Favourite) *api.Like {
	return &api.Like{
		CreatedAt: like.CreatedAt,
		Id:        int(like.ID),
		PostID:    int(like.StatusID),
		UserID:    int(like.AccountID),
	}
}

func AccountToUserAPI(account *db.Account, followersCount int, followingCount int) *api.User {
	origin := "local"
	if account.Domain.Valid && account.Domain.String != "" {
		origin = account.Domain.String
	}
	return &api.User{
		Id:             int(account.ID),
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

func StatusToComment(status *db.GetCommentsByPostIdRow) *api.Comment {
	postID := 0
	if status.InReplyToID.Valid {
		postID = int(status.InReplyToID.Int32)
	}
	username := status.Username
	return &api.Comment{
		Id:        int(status.ID),
		PostID:    postID,
		UserID:    int(status.AccountID),
		Content:   status.Content,
		ParentID:  postID, // Comments have parentID same as postID
		CreatedAt: status.CreatedAt,
		UpdatedAt: status.UpdatedAt,
		Username:  &username,
	}
}
