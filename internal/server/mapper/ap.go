package mapper

import (
	"github.com/kibirisu/borg/internal/db"
	"github.com/kibirisu/borg/internal/domain"
)

func AccountToActor(account *db.Account) *domain.Actor {
	return &domain.Actor{
		Context:           "https://www.w3.org/ns/activitystreams",
		ID:                account.Uri,
		Type:              "Person",
		PreferredUsername: account.Username,
		Inbox:             account.InboxUri,
		Outbox:            account.OutboxUri,
		Following:         account.FollowingUri,
		Followers:         account.FollowersUri,
	}
}
