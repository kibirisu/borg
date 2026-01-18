package util

import (
	"database/sql"
	"fmt"
	"strings"
)

type URIBuilder struct {
	base string
}

type ActorURIs struct {
	Actor     string
	Inbox     string
	Outbox    string
	Followers string
	Following string
}

type StatusURIs struct {
	Status  string
	Replies string
	Create  string
}

type HandleParts struct {
	Username string
	Domain   sql.NullString
}

func NewURIBuilder(addr string) URIBuilder {
	return URIBuilder{fmt.Sprintf("http://%s", addr)}
}

func (b URIBuilder) ActorURIs(id string) ActorURIs {
	baseURI := fmt.Sprintf("%s/users/%s", b.base, id)
	return ActorURIs{
		Actor:     baseURI,
		Inbox:     fmt.Sprintf("%s/inbox", baseURI),
		Outbox:    fmt.Sprintf("%s/outbox", baseURI),
		Followers: fmt.Sprintf("%s/followers", baseURI),
		Following: fmt.Sprintf("%s/following", baseURI),
	}
}

func (b URIBuilder) StatusURIs(actorID, statusID string) StatusURIs {
	baseURI := fmt.Sprintf("%s/users/%s/statuses/%s", b.base, actorID, statusID)
	return StatusURIs{
		Status:  baseURI,
		Replies: fmt.Sprintf("%s/replies", baseURI),
		Create:  fmt.Sprintf("%s/create", baseURI),
	}
}

func (b URIBuilder) FollowURI(actorID, followID string) string {
	return fmt.Sprintf("%s/users/%s/follows/%s", b.base, actorID, followID)
}

func (b URIBuilder) FollowRequestURI(actorID, requestID string) string {
	return fmt.Sprintf("%s/users/%s/requests/%s", b.base, actorID, requestID)
}

func (b URIBuilder) LikeRequestURI(actorID, requestID string) string {
	return fmt.Sprintf("%s/users/%s/likes/%s", b.base, actorID, requestID)
}

func (b URIBuilder) AnnounceURI(actorID, announceID string) string {
	return fmt.Sprintf("%s/users/%s/reblogs/%s", b.base, actorID, announceID)
}

func ExtractDomainFromURI(uri string) string {
	res := strings.SplitN(uri, "/", 4)
	return res[2]
}

func ExtractUsernameFromAcct(acct string) string {
	handle := strings.TrimPrefix(acct, "acct:")
	res := strings.SplitN(handle, "@", 2)
	return res[0]
}

func ExtractHandleParts(acct string) HandleParts {
	res := strings.SplitN(acct, "@", 2)
	if len(res) == 2 {
		return HandleParts{
			Username: res[0],
			Domain: sql.NullString{
				String: res[1],
				Valid:  true,
			},
		}
	}
	return HandleParts{res[0], sql.NullString{}}
}
