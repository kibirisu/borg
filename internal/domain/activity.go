package domain

type Actor struct {
	Context           any    `json:"@context"`
	ID                string `json:"id"`
	Type              string `json:"type"`
	PreferredUsername string `json:"preferredUsername"`
	Inbox             string `json:"inbox"`
	Outbox            string `json:"outbox"`
	Following         string `json:"following"`
	Followers         string `json:"followers"`
}
