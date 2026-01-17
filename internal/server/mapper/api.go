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
		Content:      post.Content.String,
		CreatedAt:    post.CreatedAt,
		Id:           post.ID.String(),
		LikeCount:    likeCount,
		ShareCount:   shareCount,
		UpdatedAt:    post.UpdatedAt,
		UserID:       post.AccountID.String(),
		Username:     &acc.Username,
	}
}

func ToAPIStatus(status *db.GetStatusByIDNewRow) *api.Status {
	var inReplyToID, inReplyToAccountID *string

	if status.Status.ReblogOfID == nil {
		if status.Status.InReplyToID != nil {
			id := status.Status.InReplyToID.String()
			inReplyToID = &id
		}
		if status.Status.InReplyToAccountID != nil {
			id := status.Status.InReplyToAccountID.String()
			inReplyToAccountID = &id
		}

		res := &api.Status{
			Account: api.Account{
				Acct:           status.Acct,
				DisplayName:    status.Account.DisplayName.String,
				FollowersCount: int(status.FollowersCount),
				FollowingCount: int(status.FollowingCount),
				Id:             status.Account.ID.String(),
				Url:            status.Account.Url,
				Username:       status.Account.Username,
			},
			Content:            status.Status.Content.String,
			Favourited:         &status.Favourited,
			FavouritesCount:    int(status.FavouritesCount),
			Id:                 status.Status.ID.String(),
			InReplyToAccountId: inReplyToAccountID,
			InReplyToId:        inReplyToID,
			Reblogged:          &status.Reblogged,
			ReblogsCount:       int(status.ReblogsCount),
			RepliesCount:       int(status.RepliesCount),
			Uri:                status.Status.Uri,
		}
		return res
	}

	if status.RebloggedReplyToID != nil {
		id := status.RebloggedReplyToID.String()
		inReplyToID = &id
	}
	if status.RebloggedReplyToAccountID != nil {
		id := status.RebloggedReplyToAccountID.String()
		inReplyToAccountID = &id
	}

	res := &api.Status{
		Account: api.Account{
			Acct:           status.Acct,
			DisplayName:    status.Account.DisplayName.String,
			FollowersCount: int(status.FollowersCount),
			FollowingCount: int(status.FollowingCount),
			Id:             status.Account.ID.String(),
			Url:            status.Account.Url,
			Username:       status.Account.Username,
		},
		Content:            status.Status.Content.String,
		Favourited:         &status.Favourited,
		FavouritesCount:    int(status.FavouritesCount),
		Id:                 status.Status.ID.String(),
		InReplyToAccountId: inReplyToAccountID,
		InReplyToId:        inReplyToID,
		Reblog: &api.Status{
			Account: api.Account{
				Acct:           status.RebloggedAcct,
				DisplayName:    status.RebloggedDisplayName.String,
				FollowersCount: int(status.RebloggedFollowersCount),
				FollowingCount: int(status.RebloggedFollowingCount),
				Id:             status.Status.ReblogOfID.String(),
				Url:            ":3",
				Username:       status.RebloggedUsername.String,
			},
			Content:            status.RebloggedStatusContent.String,
			Favourited:         &status.Favourited,
			FavouritesCount:    int(status.FavouritesCount),
			Id:                 status.Status.ReblogOfID.String(),
			InReplyToAccountId: inReplyToAccountID,
			InReplyToId:        inReplyToID,
			Reblogged:          &status.Reblogged,
			ReblogsCount:       int(status.ReblogsCount),
			RepliesCount:       int(status.RepliesCount),
			Uri:                status.Status.ReblogOfUri.String,
		},
		Reblogged:    &status.Reblogged,
		ReblogsCount: int(status.ReblogsCount),
		RepliesCount: int(status.RepliesCount),
		Uri:          status.Status.Uri,
	}
	return res
}
