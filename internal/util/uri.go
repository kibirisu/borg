package util

import (
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

func NewURIBuilder(host, port string) URIBuilder {
	if port == "80" {
		return URIBuilder{fmt.Sprintf("http://%s", host)}
	}
	return URIBuilder{fmt.Sprintf("http://%s:%s", host, port)}
}

func (b URIBuilder) ActorURIs(id string) ActorURIs {
	baseURI := fmt.Sprintf("%s/user/%s", b.base, id)
	return ActorURIs{
		Actor:     baseURI,
		Inbox:     fmt.Sprintf("%s/inbox", baseURI),
		Outbox:    fmt.Sprintf("%s/outbox", baseURI),
		Followers: fmt.Sprintf("%s/followers", baseURI),
		Following: fmt.Sprintf("%s/following", baseURI),
	}
}

func (b URIBuilder) StatusURIs(actorID, statusID string) StatusURIs {
	baseURI := fmt.Sprintf("%s/user/%s/statuses/%s", b.base, actorID, statusID)
	return StatusURIs{
		Status:  baseURI,
		Replies: fmt.Sprintf("%s/replies", baseURI),
		Create:  fmt.Sprintf("%s/create", baseURI),
	}
}

func ExtractDomainFromURI(uri string) string {
	res := strings.SplitN(uri, "/", 4)
	return res[2]
}
