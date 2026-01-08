package mapper

import (
	"database/sql"

	"github.com/kibirisu/borg/internal/api"
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
func NewPostToDB(newPost *api.NewPost, isLocal bool) *db.CreateStatusParams {
	return &db.CreateStatusParams{
		Uri:			"",
		Url: 			"TODO",
		Local:			sql.NullBool{Bool: isLocal, Valid: true},
		Content:		newPost.Content,
		AccountID:		int32(newPost.UserID),
		InReplyToID:	sql.NullInt32{},
		ReblogOfID:		sql.NullInt32{},
	}
}
func NewCommentToDB(comment *api.NewComment) *db.CreateStatusParams {
	return &db.CreateStatusParams{
		Uri: "TODO",
		Url: "TODO",
		Local: sql.NullBool{Bool: true, Valid: true},
		Content:    comment.Content,
		AccountID:  int32(comment.UserID),
		InReplyToID: sql.NullInt32{Int32: int32(comment.PostID), Valid: true},
		ReblogOfID:  sql.NullInt32{ Valid: false},
	}
}

func NewShareToDB(share *api.NewShare) *db.CreateStatusParams {
	return &db.CreateStatusParams{
		Uri: "TODO",
		Url: "TODO",
		Local: sql.NullBool{Bool: true, Valid: true},
		Content:    "",
		AccountID:  int32(share.UserID),
		InReplyToID: sql.NullInt32{ Valid: false},
		ReblogOfID:  sql.NullInt32{ Valid: true, Int32: int32(share.PostID)},
	}
}

func UpdatePostToDB(update *api.UpdatePost, id int) *db.UpdateStatusParams {
    var content string
    if update.Content != nil {
        content = *update.Content
    }
    
    return &db.UpdateStatusParams{
        Content: content,
        ID:      int32(id),
    }
}

func UpdateUserToDB(update *api.UpdateUser, id int) *db.UpdateAccountParams {
	var displayName sql.NullString
	if update.Bio != nil {
		displayName = sql.NullString{String: *update.Bio, Valid: true}
	}
	return &db.UpdateAccountParams{
		DisplayName: displayName,
		ID:          int32(id),
	}
}