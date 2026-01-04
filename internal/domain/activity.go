package domain

import (
	"encoding/json"
	"time"
)

type ActorOld struct {
	Context           any    `json:"@context"`
	ID                string `json:"id"`
	Type              string `json:"type"`
	PreferredUsername string `json:"preferredUsername"`
	Inbox             string `json:"inbox"`
	Outbox            string `json:"outbox"`
	Following         string `json:"following"`
	Followers         string `json:"followers"`
}

type Activity struct {
	Context any             `json:"@context"`
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Actor   json.RawMessage `json:"actor"`
	Object  json.RawMessage `json:"object"`
}

type Create struct {
	ID     string          `json:"id"`
	Type   string          `json:"type"`
	Actor  json.RawMessage `json:"actor"`
	Object json.RawMessage `json:"object"`
}

type Follow struct {
	ID     string          `json:"id"`
	Type   string          `json:"type"`
	Actor  json.RawMessage `json:"actor"`
	Object json.RawMessage `json:"object"`
}

type NoteOld struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Published    time.Time `json:"published"`
	AttributedTo string    `json:"attributedTo"`
	Content      string    `json:"content"`
	To           []string  `json:"to"`
}

type Accept struct {
	Context any    `json:"@context"`
	ID      string `json:"id"`
	Type    string `json:"type"`
	Actor   string `json:"actor"`
	Object  any    `json:"object"`
}
type Like struct {
	Context any             `json:"@context"`
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Actor   json.RawMessage `json:"actor"`
	Object  json.RawMessage `json:"object"`
}
