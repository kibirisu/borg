package mapper

import (
	"database/sql"
	"encoding/json"

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
func ActorToAccountCreate(account *domain.Actor, domain string) *db.CreateActorParams {
	return &db.CreateActorParams{
		Username: account.PreferredUsername,
		Uri: account.ID,
		DisplayName: sql.NullString{
			String: account.PreferredUsername, 
			Valid:  true,
		},
		Domain: sql.NullString{
			String: domain, 
			Valid:  true,
		},
		InboxUri: account.Inbox,
		OutboxUri: account.Outbox,
		FollowersUri: account.Followers,
		FollowingUri: account.Following,
		Url: "", //TODO
	}
}
func ToFollow(data []byte) (*domain.Follow, error) {
	var f domain.Follow
	err := json.Unmarshal(data, &f)
	return &f, err
}
func ToCreate(data []byte) (*domain.Create, error) {
	var f domain.Create
	err := json.Unmarshal(data, &f)
	return &f, err
}

func PostToNote(post *db.Status, senderURI string, receiverURIs []string) *domain.Note {
    if len(receiverURIs) == 0 {
        receiverURIs = []string{"https://www.w3.org/ns/activitystreams#Public"}
    }

    return &domain.Note{
        ID:           post.Uri,
        Type:         "Note",
        Published:    post.CreatedAt,
        AttributedTo: senderURI,     
        Content:      post.Content,
        To:           receiverURIs,
    }
}
