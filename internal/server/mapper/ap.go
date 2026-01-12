package mapper

import (

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/db"
)

func StatusToCreateActivity(status db.Status, author db.Account, statusParent *db.Status) ap.CreateActivitier {
	actor := ap.NewActor(nil)
	actor.SetLink(author.Uri)
	note := ap.NewNote(nil)

	noteParent := ap.NewNote(nil)
	if statusParent != nil {
		noteParent.SetLink(statusParent.Uri)
	}
	replies := ap.NewNoteCollection(nil)
	note.SetObject(ap.Note{ 
		ID: status.Uri,
		Type: "Note",
		Content: status.Content,
		InReplyTo: noteParent,
		Published: status.CreatedAt,
		AttributedTo: actor,
		To:          []string {},
		CC:          []string {},
		Replies:     replies,
	})

	activity := ap.NewCreateActivity(nil)
	activity.SetObject(ap.Activity[ap.Note]{
		ID: "TODO",
		Type: "Create",
		Actor: actor,
		Object: note,
	})
	return activity
}
