package domain

import "time"

type ActivityType string

const (
	ActivityTypeAccept        ActivityType = "Accept"
	ActivityTypeAnnounce      ActivityType = "Announce"
	ActivityTypeCreate        ActivityType = "Create"
	ActivityTypeFollow        ActivityType = "Follow"
	ActivityTypeLike          ActivityType = "Like"
	ActivityTypeUnimplemented ActivityType = "Unimplemented"
)

type Object struct {
	Context        any             `json:"@context,omitempty"`
	ID             string          `json:"id"`
	Type           string          `json:"type"`
	Actor          *Actor          `json:",inline"`
	Publication    *Publication    `json:",inline"`
	Note           *Note           `json:",inline"`
	Collection     *Collection     `json:",inline"`
	CollectionPage *CollectionPage `json:",inline"`
	ActivityActor  *ObjectOrLink   `json:"actor,omitempty"`
	ActivityObject *ObjectOrLink   `json:"object,omitempty"`
}

type Actor struct {
	PreferredUsername string `json:"preferredUsername"`
	Inbox             string `json:"inbox"`
	Outbox            string `json:"outbox"`
	Following         string `json:"following"`
	Followers         string `json:"followers"`
}

type Publication struct {
	Published    time.Time     `json:"published"`
	AttributedTo *ObjectOrLink `json:"attributedTo,omitempty"`
	To           []string      `json:"to"`
	CC           []string      `json:"cc"`
}

type Note struct {
	Content   string        `json:"content"`
	InReplyTo *ObjectOrLink `json:"inReplyTo"`
	Replies   *ObjectOrLink  `json:"replies,omitempty"`
}

type Collection struct {
	First ObjectOrLink `json:"first"`
}

type CollectionPage struct {
	Next   *ObjectOrLink  `json:"next,omitempty"`
	PartOf ObjectOrLink   `json:"partOf"`
	Items  []ObjectOrLink `json:"items"`
}
