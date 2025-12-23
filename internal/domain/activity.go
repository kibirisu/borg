package domain

import "time"

type ActivityType string

const (
	ActivityTypeAccept ActivityType = "Accept"
	ActivityTypeCreate ActivityType = "Create"
	ActivityTypeFollow ActivityType = "Follow"
)

type Object struct {
	Context any    `json:"@context"`
	Type    string `json:"type"`
	ID      string `json:"id"`
}

type Actor struct {
	Base              Object `json:",inline"`
	PreferredUsername string `json:"preferredUsername"`
	Inbox             string `json:"inbox"`
	Outbox            string `json:"outbox"`
	Following         string `json:"following"`
	Followers         string `json:"followers"`
}

type Activity struct {
	Base   Object         `json:",inline"`
	Actor  Actor          `json:"actor"`
	Object ActivityObject `json:"object"`
}

type ActivityObject struct {
	Type   ActivityType
	Note   *Note
	Actor  *Actor
	Follow *Activity
}

type NoteProperties struct {
	Published    time.Time `json:"published"`
	AttributedTo string    `json:"attributedTo"`
	To           []string  `json:"to"`
	CC           []string  `json:"cc"`
	Content      string    `json:"content"`
}

type Note struct {
	Base       Object         `json:",inline"`
	Properties NoteProperties `json:",inline"`
	InReplyTo  struct {
		Base       Object         `json:",inline"`
		Properties NoteProperties `json:",inline"`
	} `json:"inReplyTo"`
}
