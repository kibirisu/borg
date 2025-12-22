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

func PostToAPI(status db.GetAllStatusesRow) api.Post {
	username := status.Username
	return api.Post{
		Id:           int(status.ID),
		UserID:       int(status.AccountID),
		Content:      status.Content,
		LikeCount:    0, // TODO: Calculate from favourites table
		ShareCount:   0, // TODO: Calculate from reblogs
		CommentCount: 0, // TODO: Calculate from replies
		CreatedAt:    status.CreatedAt,
		UpdatedAt:    status.UpdatedAt,
		Username:     &username,
	}
}
