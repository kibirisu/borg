package mapper

import (
	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/db"
)

func ToAPIAccount(account *db.GetAccountByIDRow) *api.Account {
	return &api.Account{
		Acct:           account.Acct,
		DisplayName:    account.Account.DisplayName.String,
		FollowersCount: int(account.FollowersCount),
		FollowingCount: int(account.FollowingCount),
		Id:             account.Account.ID.String(),
		Url:            account.Account.Url,
		Username:       account.Account.Username,
	}
}

func ToAPIStatus(status *db.GetStatusByIDRow) *api.Status {
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

		return &api.Status{
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
	}

	if status.RebloggedReplyToID != nil {
		id := status.RebloggedReplyToID.String()
		inReplyToID = &id
	}
	if status.RebloggedReplyToAccountID != nil {
		id := status.RebloggedReplyToAccountID.String()
		inReplyToAccountID = &id
	}

	return &api.Status{
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
			Uri:                status.RebloggedUri.String,
		},
		Reblogged:    &status.Reblogged,
		ReblogsCount: int(status.ReblogsCount),
		RepliesCount: int(status.RepliesCount),
		Uri:          status.Status.Uri,
	}
}
