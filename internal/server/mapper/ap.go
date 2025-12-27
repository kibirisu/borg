package mapper

import (
	"database/sql"
	"encoding/json"
	"strconv"

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
		Uri:      account.ID,
		DisplayName: sql.NullString{
			String: account.PreferredUsername,
			Valid:  true,
		},
		Domain: sql.NullString{
			String: domain,
			Valid:  true,
		},
		InboxUri:     account.Inbox,
		OutboxUri:    account.Outbox,
		FollowersUri: account.Followers,
		FollowingUri: account.Following,
		Url:          "", // TODO
	}
}
func DBToFollow(follow *db.Follow, follower *db.Account, followee *db.Account) *domain.Follow {
	followerURI, _ := json.Marshal(follower.Uri)
	followedURI, _ := json.Marshal(followee.Uri)
	return &domain.Follow{
		ID: follow.Uri,
		Type: "Follow",
		Actor:  json.RawMessage(followerURI),
		Object: json.RawMessage(followedURI),
	};
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

func PostToCreateNote(post *db.Status, poster *db.Account, receiverURIs []string) *domain.Create {
	if len(receiverURIs) == 0 {
		receiverURIs = []string{"https://www.w3.org/ns/activitystreams#Public"}
	}

	note := domain.Note{
		ID:           post.Uri,
		Type:         "Note",
		Published:    post.CreatedAt,
		AttributedTo: poster.Uri,
		Content:      post.Content,
		To:           receiverURIs,
	}
	noteBytes, err := json.Marshal(note)
	if err != nil {
		return nil
	}
	actorBytes, err := json.Marshal(poster.Uri)
	if err != nil {
		return nil
	}

	activity := domain.Create{
		ID:     poster.Uri + "/posts/" + strconv.Itoa(int(post.ID)),
		Type:   "Create",
		Actor:  json.RawMessage(actorBytes),
		Object: json.RawMessage(noteBytes),
	}
	return &activity;
}
