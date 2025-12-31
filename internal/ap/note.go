package ap

import (
	"time"

	"github.com/kibirisu/borg/internal/domain"
)

type Noter interface {
	Objecter[Note]
}

type Note struct {
	ID           string
	Type         string
	Content      string
	InReplyTo    Noter
	Published    time.Time
	AttributedTo Actorer
	To           []string
	CC           []string
	Replies      Collectioner
}

type note struct {
	object
}

var _ Noter = (*note)(nil)

// GetObject implements Noter.
// Subtle: this method shadows the method (object).GetObject of note.object.
func (n *note) GetObject() Note {
	obj := n.raw.Object
	return Note{
		ID:           obj.ID,
		Type:         obj.Type,
		Content:      obj.Note.Content,
		InReplyTo:    &note{object{obj.Note.InReplyTo}},
		Published:    obj.Publication.Published,
		AttributedTo: &actor{object{obj.Publication.AttributedTo}},
		To:           obj.Publication.To,
		CC:           obj.Publication.CC,
		Replies:      nil, // TODO
	}
}

// SetObject implements Noter.
// Subtle: this method shadows the method (object).SetObject of note.object.
func (n *note) SetObject(note Note) {
	n.raw = &domain.ObjectOrLink{
		Object: &domain.Object{
			ID:   note.ID,
			Type: note.Type,
			Publication: &domain.Publication{
				Published:    note.Published,
				AttributedTo: note.AttributedTo.GetRaw(),
				To:           note.To,
				CC:           note.CC,
			},
			Note: &domain.Note{
				Content:   note.Content,
				InReplyTo: note.InReplyTo.GetRaw(),
				Replies:   *note.Replies.GetRaw(),
			},
		},
	}
}
