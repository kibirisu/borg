package mapper

import (
	"database/sql"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/db"
)

func NewPostToDB(newPost *api.NewPost, isLocal bool) *db.CreateStatusParams {
	return &db.CreateStatusParams{
		Url:         "TODO",
		Local:       sql.NullBool{Bool: isLocal, Valid: true},
		Content:     newPost.Content,
		AccountID:   int32(newPost.UserID),
		InReplyToID: sql.NullInt32{},
		ReblogOfID:  sql.NullInt32{},
	}
}

func NewCommentToDB(comment *api.NewComment) *db.CreateStatusParams {
	return &db.CreateStatusParams{
		Url:         "TODO",
		Local:       sql.NullBool{Bool: true, Valid: true},
		Content:     comment.Content,
		AccountID:   int32(comment.UserID),
		InReplyToID: sql.NullInt32{Int32: int32(comment.PostID), Valid: true},
		ReblogOfID:  sql.NullInt32{Valid: false},
	}
}

func NewShareToDB(share *api.NewShare) *db.CreateStatusParams {
	return &db.CreateStatusParams{
		Url:         "TODO",
		Local:       sql.NullBool{Bool: true, Valid: true},
		Content:     "",
		AccountID:   int32(share.UserID),
		InReplyToID: sql.NullInt32{Valid: false},
		ReblogOfID:  sql.NullInt32{Valid: true, Int32: int32(share.PostID)},
	}
}
