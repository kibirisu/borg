package mapper

import (
	"database/sql"
	"encoding/json"
	"strconv"

	"github.com/kibirisu/borg/internal/db"
	"github.com/kibirisu/borg/internal/domain"
)

func AccountToActor(account *db.Account) *domain.ActorOld {
	return &domain.ActorOld{
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

func ActorToAccountCreate(account *domain.ActorOld, domain string) *db.CreateActorParams {
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
	followerURI, err := json.Marshal(follower.Uri)
	if err != nil {
		return nil
	}
	followedURI, err := json.Marshal(followee.Uri)
	if err != nil {
		return nil
	}
	return &domain.Follow{
		ID:     follow.Uri,
		Type:   "Follow",
		Actor:  json.RawMessage(followerURI),
		Object: json.RawMessage(followedURI),
	}
}

func DBToFavourite(fav *db.Favourite, liker *db.Account, post *db.Status) *domain.Like {
	likerURI, err := json.Marshal(liker.Uri)
	if err != nil {
		return nil
	}
	postURI, err := json.Marshal(post.Uri)
	if err != nil {
		return nil
	}

	return &domain.Like{
		ID:     fav.Uri,
		Type:   "Like",
		Actor:  json.RawMessage(likerURI),
		Object: json.RawMessage(postURI),
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

func PostToCreateNote(post *db.Status, poster *db.Account, followersURI string) *domain.Create {
	note := domain.NoteOld{
		ID:           post.Uri,
		Type:         "Note",
		Published:    post.CreatedAt,
		AttributedTo: poster.Uri,
		Content:      post.Content,
		To:           []string{followersURI},
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
	return &activity
}

func PostToUpdateNote(post *db.Status, poster *db.Account, followersUri string) *domain.Update {
	note := domain.Note{
		ID:           post.Uri,
		Type:         "Note",
		Published:    post.CreatedAt,
		AttributedTo: poster.Uri,
		Content:      post.Content,
		To:           []string{followersUri},
	}
	noteBytes, err := json.Marshal(note)
	if err != nil {
		return nil
	}
	actorBytes, err := json.Marshal(poster.Uri)
	if err != nil {
		return nil
	}

	activity := domain.Update{
		ID:     poster.Uri + "/posts/" + strconv.Itoa(int(post.ID)) + "/update",
		Type:   "Update",
		Actor:  json.RawMessage(actorBytes),
		Object: json.RawMessage(noteBytes),
	}
	return &activity
}