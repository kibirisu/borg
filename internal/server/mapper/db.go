package mapper

import (
	"database/sql"

	"github.com/kibirisu/borg/internal/db"
	"github.com/kibirisu/borg/internal/domain"
)

func ActorToDB(actor *domain.Actor, domain string) *db.CreateActorParams {
	return &db.CreateActorParams{
		Username:    actor.PreferredUsername,
		Uri:         actor.ID,
		DisplayName: sql.NullString{},
		Domain:      sql.NullString{String: domain, Valid: true},
		InboxUri:    actor.Inbox,
		OutboxUri:   actor.Outbox,
		Url:         "TODO",
	}
}
