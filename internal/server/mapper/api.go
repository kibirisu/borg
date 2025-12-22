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
