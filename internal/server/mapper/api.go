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
