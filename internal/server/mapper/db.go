package mapper

import (
	"database/sql"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/db"
)

func NewPostToDB(newPost *api.NewPost, isLocal bool) *db.CreateStatusParams {
	return &db.CreateStatusParams{
		Uri:         "",
		Url:         "TODO",
		Local:       sql.NullBool{Bool: isLocal, Valid: true},
		Content:     newPost.Content,
		AccountID:   int32(newPost.UserID),
		InReplyToID: sql.NullInt32{},
		ReblogOfID:  sql.NullInt32{},
	}
}
