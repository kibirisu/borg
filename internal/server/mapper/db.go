package mapper

import (
	"database/sql"

	"github.com/rs/xid"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/db"
)

func NewPostToDB(newPost *api.NewPost, isLocal bool) *db.CreateStatusParams {
	actorID, err := xid.FromString(newPost.UserID)
	if err != nil {
		return nil
	}
	return &db.CreateStatusParams{
		Url:         "TODO",
		Local:       sql.NullBool{Bool: isLocal, Valid: true},
		Content:     newPost.Content,
		AccountID:   actorID,
		InReplyToID: nil,
		ReblogOfID:  nil,
	}
}

func NewCommentToDB(comment *api.NewComment) *db.CreateStatusParams {
	actorID, err := xid.FromString(comment.UserID)
	if err != nil {
		return nil
	}
	replyOfID, err := xid.FromString(comment.PostID)
	if err != nil {
		return nil
	}
	return &db.CreateStatusParams{
		Url:         "TODO",
		Local:       sql.NullBool{Bool: true, Valid: true},
		Content:     comment.Content,
		AccountID:   actorID,
		InReplyToID: &replyOfID,
		ReblogOfID:  nil,
	}
}

func NewShareToDB(share *api.NewShare) *db.CreateStatusParams {
	actorID, err := xid.FromString(share.UserID)
	if err != nil {
		return nil
	}
	reblogOfID, err := xid.FromString(share.PostID)
	if err != nil {
		return nil
	}
	return &db.CreateStatusParams{
		Url:         "TODO",
		Local:       sql.NullBool{Bool: true, Valid: true},
		Content:     "",
		AccountID:   actorID,
		InReplyToID: nil,
		ReblogOfID:  &reblogOfID,
	}
}
