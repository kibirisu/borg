package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/config"
	"github.com/kibirisu/borg/internal/db"
	proc "github.com/kibirisu/borg/internal/processing"
	repo "github.com/kibirisu/borg/internal/repository"
	"github.com/kibirisu/borg/internal/util"
	"github.com/kibirisu/borg/internal/worker"
)

type AppService interface {
	Register(context.Context, api.AuthForm) error
	Login(context.Context, api.AuthForm) (string, error)
	CreateStatus(context.Context, api.NewPost, LoginData) (worker.Job, error)
	GetAccountFollowers(context.Context, int) ([]db.Account, error)
	GetAccountFollowing(context.Context, int) ([]db.Account, error)
	GetLocalAccount(context.Context, string) (*db.Account, error)
	AddRemoteAccount(ctx context.Context, remote *db.CreateActorParams) (*db.Account, error)
	CreateFollow(ctx context.Context, follow *db.CreateFollowParams) (*db.Follow, error)
	AddNote(context.Context, db.CreateStatusParams) (db.Status, error)
	AddFavourite(context.Context, int, int) (db.Favourite, error)
	FollowAccount(context.Context, int, int) (*db.Follow, error)
	GetAccountByID(context.Context, int) (db.Account, error)
	UpdateAccount(context.Context, int, *string) (db.Account, error)
	GetAccount(context.Context, db.GetAccountParams) (*db.Account, error)
	GetLocalPosts(context.Context) ([]db.GetLocalStatusesRow, error)
	GetPostByAccountID(context.Context, int) ([]db.GetStatusesByAccountIdRow, error)
	GetPostByID(context.Context, int) (*db.Status, error)
	UpdatePost(context.Context, int, string) (*db.Status, error)
	DeletePost(context.Context, int) error
	GetPostLikes(context.Context, int) ([]db.Favourite, error)
	GetPostShares(context.Context, int) ([]db.Status, error)
	GetPostByIDWithMetadata(context.Context, int) (*db.GetStatusByIdWithMetadataRow, error)
	GetLikedPostsByAccountID(context.Context, int) ([]db.GetLikedPostsByAccountIdRow, error)
	GetSharedPostsByAccountID(context.Context, int) ([]db.GetSharedPostsByAccountIdRow, error)
	GetTimelinePostsByAccountID(context.Context, int) ([]db.GetTimelinePostsByAccountIdRow, error)
	GetCommentsByPostID(context.Context, int) ([]db.GetCommentsByPostIdRow, error)
	// EW, idk if this should stay here
	DeliverToFollowers(http.ResponseWriter, *http.Request, int, func(recipientURI string) any)
}

type appService struct {
	store    repo.Store
	prcessor proc.Processor
	conf     *config.Config
}

type LoginData struct {
	ID       int
	Username string
}

var _ AppService = (*appService)(nil)

// Register implements AppService.
func (s *appService) Register(ctx context.Context, form api.AuthForm) error {
	uri := fmt.Sprintf("http://%s:%s/user/%s", s.conf.ListenHost, s.conf.ListenPort, form.Username)
	log.Printf("register: creating actor username=%s uri=%s", form.Username, uri)
	actor, err := s.store.Accounts().Create(ctx, db.CreateActorParams{
		Username:    form.Username,
		Uri:         uri,
		DisplayName: sql.NullString{}, // hassle to maintain that, gonna abandon display name
		Domain:      sql.NullString{},
		InboxUri:    uri + "/inbox",
		OutboxUri:   uri + "/outbox",
		Url: fmt.Sprintf(
			"http://%s:%s/profiles/%s",
			s.conf.ListenHost,
			s.conf.ListenPort,
			form.Username,
		),
		FollowersUri: fmt.Sprintf(
			"http://%s:%s/user/%s/followers",
			s.conf.ListenHost,
			s.conf.ListenPort,
			form.Username,
		),
		FollowingUri: fmt.Sprintf(
			"http://%s:%s/user/%s/following",
			s.conf.ListenHost,
			s.conf.ListenPort,
			form.Username,
		),
	})
	if err != nil {
		log.Printf("register: failed to create actor username=%s err=%v", form.Username, err)
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("register: failed to hash password username=%s err=%v", form.Username, err)
		return err
	}
	if err = s.store.Users().Create(ctx, db.CreateUserParams{
		AccountID:    actor.ID,
		PasswordHash: string(hash),
	}); err != nil {
		log.Printf("register: failed to create user username=%s err=%v", form.Username, err)
		return err
	}
	log.Printf(
		"register: user and actor created username=%s account_id=%d",
		form.Username,
		actor.ID,
	)
	return nil
}

// Login implements AppService.
func (s *appService) Login(ctx context.Context, form api.AuthForm) (string, error) {
	auth, err := s.store.Users().GetByUsername(ctx, form.Username)
	if err != nil {
		return "", err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(auth.PasswordHash), []byte(form.Password)); err != nil {
		return "", err
	}
	token, err := issueToken(auth.ID, form.Username, s.conf.ListenHost, s.conf.JWTSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}

// CreateStatus implements AppService.
func (s *appService) CreateStatus(
	ctx context.Context,
	status api.NewPost,
	login LoginData,
) (worker.Job, error) {
	actorURI := fmt.Sprintf(
		"http://%s:%s/user/%s",
		s.conf.ListenHost,
		s.conf.ListenPort,
		login.Username,
	)
	statusURI := fmt.Sprintf("%s/statuses/%s", actorURI, uuid.New())
	createdStatus, err := s.store.Statuses().Create(ctx, db.CreateStatusParams{
		Local: sql.NullBool{
			Bool:  true,
			Valid: true,
		},
		Content:   status.Content,
		AccountID: int32(login.ID),
		Uri:       statusURI,
	})
	if err != nil {
		return nil, err
	}
	return func(ctx context.Context) error {
		actor := ap.NewActor(nil)
		actor.SetLink(actorURI)
		status := ap.NewNote(nil)
		replies := ap.NewNoteCollection(nil)
		page := ap.NewNoteCollectionPage(nil)
		page.SetObject(ap.CollectionPage[ap.Note]{
			ID:     "None",
			Type:   "CollectionPage",
			Next:   ap.NewNoteCollectionPage(nil),
			PartOf: replies,
			Items:  []ap.Objecter[ap.Note]{},
		})
		replies.SetObject(ap.Collection[ap.Note]{
			ID:    fmt.Sprintf("%s/replies", statusURI),
			Type:  "Collection",
			First: nil,
		})
		status.SetObject(ap.Note{
			ID:           createdStatus.Uri,
			Type:         "Note",
			Content:      createdStatus.Content,
			InReplyTo:    ap.NewNote(nil),
			Published:    createdStatus.CreatedAt,
			AttributedTo: actor,
			To:           []string{},
			CC:           []string{},
			Replies:      replies,
		})
		return s.prcessor.PropagateStatus(ctx, status)
	}, nil
}

// GetLocalAccount implements AppService.
func (s *appService) GetLocalAccount(ctx context.Context, username string) (*db.Account, error) {
	user, err := s.store.Accounts().GetLocalByUsername(ctx, username)
	return &user, err
}

func (s *appService) AddRemoteAccount(
	ctx context.Context,
	remote *db.CreateActorParams,
) (*db.Account, error) {
	if !remote.Domain.Valid {
		return nil, errors.New("domain must be a remote server")
	}
	acc, err := s.store.Accounts().Create(ctx, *remote)
	return &acc, err
}

func (s *appService) CreateFollow(
	ctx context.Context,
	follow *db.CreateFollowParams,
) (*db.Follow, error) {
	if follow.Uri == "" {
		follow.Uri = fmt.Sprintf(
			"http://%s:%s/follows/%s",
			s.conf.ListenHost,
			s.conf.ListenPort,
			uuid.NewString(),
		)
	}
	return s.store.Follows().Create(ctx, *follow)
}

// AddNote implements AppService.
func (s *appService) AddNote(ctx context.Context, note db.CreateStatusParams) (db.Status, error) {
	if note.Uri == "" {
		note.Uri = fmt.Sprintf(
			"http://%s:%s/statuses/%s",
			s.conf.ListenHost,
			s.conf.ListenPort,
			uuid.NewString(),
		)
	}
	if note.Url == "" {
		note.Url = note.Uri
	}
	return s.store.Statuses().Create(ctx, note)
}

// GetAccount implements AppService.
func (s *appService) GetAccount(
	ctx context.Context,
	account db.GetAccountParams,
) (*db.Account, error) {
	res, err := s.store.Accounts().Get(ctx, account)
	return &res, err
}

// GetAccountByID implements AppService.
func (s *appService) GetAccountByID(
	ctx context.Context, accountID int,
) (db.Account, error) {
	return s.store.Accounts().GetByID(ctx, accountID)
}

// UpdateAccount implements AppService.
func (s *appService) UpdateAccount(
	ctx context.Context, accountID int, bio *string,
) (db.Account, error) {
	return s.store.Accounts().Update(ctx, accountID, bio)
}

// AddFavourite implements AppService.
func (s *appService) AddFavourite(
	ctx context.Context, accountID int, postID int,
) (db.Favourite, error) {
	params := db.CreateFavouriteParams{
		AccountID: int32(accountID),
		StatusID:  int32(postID),
		Uri: fmt.Sprintf(
			"http://%s:%s/likes/%s",
			s.conf.ListenHost,
			s.conf.ListenPort,
			uuid.NewString(),
		),
	}
	return s.store.Favourites().Create(ctx, params)
}

// GetAccountFollowers implements AppService.
func (s *appService) GetAccountFollowers(
	ctx context.Context, accountID int,
) ([]db.Account, error) {
	return s.store.Accounts().GetFollowers(ctx, accountID)
}

// GetAccountFollowing implements AppService.
func (s *appService) GetAccountFollowing(
	ctx context.Context, accountID int,
) ([]db.Account, error) {
	return s.store.Accounts().GetFollowing(ctx, accountID)
}

func (s *appService) GetPostByIDWithMetadata(
	ctx context.Context,
	id int,
) (*db.GetStatusByIdWithMetadataRow, error) {
	status, err := s.store.Statuses().GetByIDWithMetadata(ctx, id)
	if err != nil {
		return nil, err
	} else {
		return &status, nil
	}
}

func (s *appService) GetPostByID(ctx context.Context, id int) (*db.Status, error) {
	status, err := s.store.Statuses().GetByID(ctx, id)
	if err != nil {
		return nil, err
	} else {
		return &status, nil
	}
}

func (s *appService) UpdatePost(ctx context.Context, id int, content string) (*db.Status, error) {
	status, err := s.store.Statuses().Update(ctx, id, content)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (s *appService) DeletePost(ctx context.Context, id int) error {
	return s.store.Statuses().DeleteByID(ctx, int32(id))
}

func (s *appService) GetPostLikes(ctx context.Context, id int) ([]db.Favourite, error) {
	return s.store.Favourites().GetByPost(ctx, id)
}

func (s *appService) GetPostShares(ctx context.Context, id int) ([]db.Status, error) {
	return s.store.Statuses().GetShares(ctx, id)
}

// FollowAccount implements AppService.
func (s *appService) FollowAccount(
	ctx context.Context,
	follower int,
	followee int,
) (*db.Follow, error) {
	createParams := db.CreateFollowParams{
		Uri: fmt.Sprintf(
			"http://%s:%s/follows/%s",
			s.conf.ListenHost,
			s.conf.ListenPort,
			uuid.NewString(),
		),
		AccountID:       int32(follower),
		TargetAccountID: int32(followee),
	}
	return s.store.Follows().Create(ctx, createParams)
}

func (s *appService) DeliverToFollowers(
	w http.ResponseWriter, r *http.Request, userID int,
	build func(recipientURI string) any,
) {
	followers, err := s.GetAccountFollowers(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	for _, follower := range followers {
		if !follower.Domain.Valid {
			continue
		}
		payload := build(follower.Uri)
		util.DeliverToEndpoint(follower.InboxUri, payload)
	}
}

func (s *appService) GetPostByAccountID(
	ctx context.Context,
	id int,
) ([]db.GetStatusesByAccountIdRow, error) {
	return s.store.Accounts().GetPosts(ctx, id)
}

func (s *appService) GetLocalPosts(ctx context.Context) ([]db.GetLocalStatusesRow, error) {
	return s.store.Statuses().GetLocalStatuses(ctx)
}

func (s *appService) GetLikedPostsByAccountID(
	ctx context.Context,
	accountID int,
) ([]db.GetLikedPostsByAccountIdRow, error) {
	return s.store.Favourites().GetLikedPostsByAccountID(ctx, accountID)
}

func (s *appService) GetSharedPostsByAccountID(
	ctx context.Context,
	accountID int,
) ([]db.GetSharedPostsByAccountIdRow, error) {
	return s.store.Statuses().GetSharedPostsByAccountID(ctx, accountID)
}

func (s *appService) GetTimelinePostsByAccountID(
	ctx context.Context,
	accountID int,
) ([]db.GetTimelinePostsByAccountIdRow, error) {
	return s.store.Statuses().GetTimelinePostsByAccountID(ctx, accountID)
}

func (s *appService) GetCommentsByPostID(
	ctx context.Context,
	postID int,
) ([]db.GetCommentsByPostIdRow, error) {
	return s.store.Statuses().GetCommentsByPostID(ctx, postID)
}
