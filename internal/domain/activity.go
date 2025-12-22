package domain

import "time"

type Object struct {
	Context any    `json:"@context"`
	Type    string `json:"type"`
	ID      string `json:"id"`
}

type Actor struct {
	Object
	PreferredUsername string `json:"preferredUsername"`
	Inbox             string `json:"inbox"`
	Outbox            string `json:"outbox"`
	Following         string `json:"following"`
	Followers         string `json:"followers"`
}

type Activity struct {
	Object
	Actor Actor `json:"actor"`
}

type Create struct {
	Activity
	Object Note `json:"object"`
}

type Follow struct {
	Activity
	Object Actor `json:"object"`
}

type Accept struct {
	Activity
	Object Follow `json:"object"`
}

type Note struct {
	Object
	Published    time.Time `json:"published"`
	AttributedTo string    `json:"attributedTo"`
	Content      string    `json:"content"`
	To           []string  `json:"to"`
}
