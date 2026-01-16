package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"

	"github.com/kibirisu/borg/internal/ap"
	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/config"
	"github.com/kibirisu/borg/internal/db"
	proc "github.com/kibirisu/borg/internal/processing"
	repo "github.com/kibirisu/borg/internal/repository"
	"github.com/kibirisu/borg/internal/server/auth"
	"github.com/kibirisu/borg/internal/util"
	"github.com/kibirisu/borg/internal/worker"
)

type AppService interface {
	Register(context.Context, api.AuthForm) error
	Login(context.Context, api.AuthForm) (string, error)
	CreateStatus(context.Context, api.PostApiStatusesJSONBody) (worker.Job, error)
	ViewStatus(context.Context, string) (*api.Status, error)
	GetAccountFollowers(context.Context, string) ([]db.Account, error)
	GetAccountFollowing(context.Context, string) ([]db.Account, error)
	GetLocalAccount(context.Context, string) (*db.Account, error)
	CreateFollow(ctx context.Context, follow *db.CreateFollowParams) (*db.Follow, error)
	AddNote(context.Context, db.CreateStatusParams) (db.Status, error)
	AddFavourite(context.Context, string, string) (db.Favourite, error)
	FollowAccount(context.Context, string, string) (*db.Follow, error)
	GetAccountByID(context.Context, string) (db.Account, error)
	GetAccount(context.Context, db.GetAccountParams) (*db.Account, error)
	GetLocalPosts(context.Context) ([]db.GetLocalStatusesRow, error)
	GetPostByAccountID(context.Context, string) ([]db.GetStatusesByAccountIdRow, error)
	GetPostByID(context.Context, string) (*db.Status, error)
	GetPostLikes(context.Context, string) ([]db.Favourite, error)
	GetPostShares(context.Context, string) ([]db.Status, error)
	GetPostByIDWithMetadata(context.Context, string) (*db.GetStatusByIdWithMetadataRow, error)
	GetLikedPostsByAccountId(context.Context, string) ([]db.GetLikedPostsByAccountIdRow, error)
	GetSharedPostsByAccountId(context.Context, string) ([]db.GetSharedPostsByAccountIdRow, error)
	GetTimelinePostsByAccountId(
		context.Context,
		string,
	) ([]db.GetTimelinePostsByAccountIdRow, error)
	// EW, idk if this should stay here
	DeliverToFollowers(http.ResponseWriter, *http.Request, string, func(recipientURI string) any)
}

type appService struct {
	store    repo.Store
	prcessor proc.Processor
	conf     *config.Config
	builder  util.URIBuilder
}

var _ AppService = (*appService)(nil)

// Register implements AppService.
func (s *appService) Register(ctx context.Context, form api.AuthForm) error {
	uri := fmt.Sprintf("http://%s/users/%s", s.conf.ListenHost, form.Username)
	log.Printf("register: creating actor username=%s uri=%s", form.Username, uri)
	actor, err := s.store.Accounts().Create(ctx, db.CreateActorParams{
		ID:          xid.New(),
		Username:    form.Username,
		Uri:         uri,
		DisplayName: sql.NullString{}, // hassle to maintain that, gonna abandon display name
		Domain:      sql.NullString{},
		InboxUri:    uri + "/inbox",
		OutboxUri:   uri + "/outbox",
		Url:         fmt.Sprintf("http://%s/profiles/%s", s.conf.ListenHost, form.Username),
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
		ID:           xid.New(),
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
	token, err := issueToken(auth.ID.String(), form.Username, s.conf.ListenHost, s.conf.JWTSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}

// CreateStatus implements AppService.
func (s *appService) CreateStatus(
	ctx context.Context,
	status api.PostApiStatusesJSONBody,
) (worker.Job, error) {
	token, ok := ctx.Value(auth.TokenContextKey).(*auth.TokenData)
	if !ok {
		return nil, errors.New("auth failure")
	}

	statusID := xid.New()
	statusURIs := s.builder.StatusURIs(token.ID, statusID.String())
	accountID, err := xid.FromString(token.ID)
	if err != nil {
		return nil, err
	}

	var inReplyToID *xid.ID
	if status.InReplyToID != nil {
		id, err := xid.FromString(*status.InReplyToID)
		if err != nil {
			return nil, err
		}
		inReplyToID = &id
	}

	createdStatus, err := s.store.Statuses().CreateNew(ctx, db.CreateStatusNewParams{
		ID:          statusID,
		Uri:         statusURIs.Status,
		Url:         "not needed rn",
		Content:     status.Status,
		AccountID:   accountID,
		AccountUri:  token.URI,
		InReplyToID: inReplyToID,
	})
	if err != nil {
		return nil, err
	}

	inReplyTo := ap.NewEmptyNote()
	if createdStatus.InReplyToUri.Valid {
		inReplyTo.SetLink(createdStatus.InReplyToUri.String)
	}

	actor := ap.NewEmptyActor().WithLink(token.URI)
	create := ap.NewEmptyCreateActivity().WithObject(ap.Activity[ap.Note]{
		ID:    statusURIs.Create,
		Type:  "Create",
		Actor: actor,
		Object: ap.NewEmptyNote().WithObject(ap.Note{
			ID:           statusURIs.Status,
			Type:         "Note",
			Content:      status.Status,
			InReplyTo:    inReplyTo,
			Published:    time.Now(),
			AttributedTo: actor,
			To:           []string{},
			CC:           []string{},
			Replies: ap.NewEmptyNoteCollection().WithObject(ap.Collection[ap.Note]{
				ID:   statusURIs.Replies,
				Type: "Collection",
				First: ap.NewEmptyNoteCollectionPage().WithObject(ap.CollectionPage[ap.Note]{
					ID:     "None",
					Type:   "CollectionPage",
					Next:   ap.NewEmptyNoteCollectionPage(),
					PartOf: ap.NewEmptyNoteCollection().WithLink(statusURIs.Replies),
					Items:  []ap.Objecter[ap.Note]{},
				}),
			}),
		}),
	})

	return func(ctx context.Context) error {
		return s.prcessor.DistributeObject(ctx, create.GetRaw().Object, accountID)
	}, nil
}

// ViewStatus implements AppService.
func (s *appService) ViewStatus(ctx context.Context, id string) (*api.Status, error) {
	token, ok := ctx.Value(auth.TokenContextKey).(*auth.TokenData)
	if !ok {
		return nil, errors.New("auth failure")
	}

	statusID, err := xid.FromString(id)
	if err != nil {
		return nil, err
	}
	accountID, err := xid.FromString(token.ID)
	if err != nil {
		return nil, err
	}

	status, err := s.store.Statuses().GetByIDNew(ctx, db.GetStatusByIDNewParams{
		ID:        statusID,
		AccountID: accountID,
	})
	if err != nil {
		return nil, err
	}

	var inReplyToID, inReplyToAccountID *string
	if status.Status.InReplyToID != nil {
		id := status.Status.InReplyToID.String()
		inReplyToID = &id
	}
	if status.Status.InReplyToAccountID.Valid {
		inReplyToAccountID = &status.Status.InReplyToAccountID.String
	}

	res := api.Status{
		Content:            status.Status.Content,
		Favourited:         status.Favourited,
		FavouritesCount:    int(status.FavouritesCount),
		Id:                 status.Status.ID.String(),
		InReplyToAccountId: inReplyToAccountID,
		InReplyToId:        inReplyToID,
		Reblogged:          status.Reblogged,
		ReblogsCount:       int(status.ReblogsCount),
		RepliesCount:       int(status.RepliesCount),
		Uri:                status.Status.Uri,
	}
	return &res, nil
}

// GetLocalAccount implements AppService.
func (s *appService) GetLocalAccount(ctx context.Context, username string) (*db.Account, error) {
	user, err := s.store.Accounts().GetLocalByUsername(ctx, username)
	return &user, err
}

func (s *appService) CreateFollow(
	ctx context.Context,
	follow *db.CreateFollowParams,
) (*db.Follow, error) {
	if follow.Uri == "" {
		follow.Uri = fmt.Sprintf("http://%s/follows/%s", s.conf.ListenHost, uuid.NewString())
	}
	return s.store.Follows().Create(ctx, *follow)
}

// AddNote implements AppService.
func (s *appService) AddNote(ctx context.Context, note db.CreateStatusParams) (db.Status, error) {
	if note.Uri == "" {
		note.Uri = fmt.Sprintf("http://%s/statuses/%s", s.conf.ListenHost, uuid.NewString())
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
	ctx context.Context, accountID string,
) (db.Account, error) {
	actorID, err := xid.FromString(accountID)
	if err != nil {
		return db.Account{}, err
	}
	return s.store.Accounts().GetByID(ctx, actorID)
}

// AddFavourite implements AppService.
func (s *appService) AddFavourite(
	ctx context.Context, accountID string, postID string,
) (db.Favourite, error) {
	actorID, err := xid.FromString(accountID)
	if err != nil {
		return db.Favourite{}, err
	}
	statusID, err := xid.FromString(postID)
	if err != nil {
		return db.Favourite{}, err
	}
	params := db.CreateFavouriteParams{
		AccountID: actorID,
		StatusID:  statusID,
		Uri:       fmt.Sprintf("http://%s/likes/%s", s.conf.ListenHost, uuid.NewString()),
	}
	return s.store.Favourites().Create(ctx, params)
}

// GetAccountFollowers implements AppService.
func (s *appService) GetAccountFollowers(
	ctx context.Context, accountID string,
) ([]db.Account, error) {
	id, err := xid.FromString(accountID)
	if err != nil {
		return []db.Account{}, err
	}
	return s.store.Accounts().GetFollowers(ctx, id)
}

// GetAccountFollowing implements AppService.
func (s *appService) GetAccountFollowing(
	ctx context.Context, accountID string,
) ([]db.Account, error) {
	id, err := xid.FromString(accountID)
	if err != nil {
		return []db.Account{}, err
	}
	return s.store.Accounts().GetFollowing(ctx, id)
}

func (s *appService) GetPostByIDWithMetadata(
	ctx context.Context,
	id string,
) (*db.GetStatusByIdWithMetadataRow, error) {
	statusID, err := xid.FromString(id)
	if err != nil {
		return nil, err
	}
	status, err := s.store.Statuses().GetByIDWithMetadata(ctx, statusID)
	return &status, err
}

func (s *appService) GetPostByID(ctx context.Context, id string) (*db.Status, error) {
	statusID, err := xid.FromString(id)
	if err != nil {
		return nil, err
	}
	status, err := s.store.Statuses().GetByID(ctx, statusID)
	return &status, err
}

func (s *appService) GetPostLikes(ctx context.Context, id string) ([]db.Favourite, error) {
	statusID, err := xid.FromString(id)
	if err != nil {
		return []db.Favourite{}, err
	}
	return s.store.Favourites().GetByPost(ctx, statusID)
}

func (s *appService) GetPostShares(ctx context.Context, id string) ([]db.Status, error) {
	statusID, err := xid.FromString(id)
	if err != nil {
		return []db.Status{}, err
	}
	return s.store.Statuses().GetShares(ctx, statusID)
}

// FollowAccount implements AppService.
func (s *appService) FollowAccount(
	ctx context.Context,
	follower string,
	followee string,
) (*db.Follow, error) {
	followerID, err := xid.FromString(follower)
	if err != nil {
		return nil, err
	}
	followedID, err := xid.FromString(followee)
	if err != nil {
		return nil, err
	}
	createParams := db.CreateFollowParams{
		Uri:             fmt.Sprintf("http://%s/follows/%s", s.conf.ListenHost, uuid.NewString()),
		AccountID:       followerID,
		TargetAccountID: followedID,
	}
	return s.store.Follows().Create(ctx, createParams)
}

func (s *appService) DeliverToFollowers(
	w http.ResponseWriter, r *http.Request, userID string,
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
	id string,
) ([]db.GetStatusesByAccountIdRow, error) {
	actorID, err := xid.FromString(id)
	if err != nil {
		return []db.GetStatusesByAccountIdRow{}, err
	}
	return s.store.Accounts().GetPosts(ctx, actorID)
}

func (s *appService) GetLocalPosts(ctx context.Context) ([]db.GetLocalStatusesRow, error) {
	return s.store.Statuses().GetLocalStatuses(ctx)
}

func (s *appService) GetLikedPostsByAccountId(
	ctx context.Context,
	accountID string,
) ([]db.GetLikedPostsByAccountIdRow, error) {
	actorID, err := xid.FromString(accountID)
	if err != nil {
		return []db.GetLikedPostsByAccountIdRow{}, err
	}
	return s.store.Favourites().GetLikedPostsByAccountId(ctx, actorID)
}

func (s *appService) GetSharedPostsByAccountId(
	ctx context.Context,
	accountID string,
) ([]db.GetSharedPostsByAccountIdRow, error) {
	actorID, err := xid.FromString(accountID)
	if err != nil {
		return []db.GetSharedPostsByAccountIdRow{}, err
	}
	return s.store.Statuses().GetSharedPostsByAccountId(ctx, actorID)
}

func (s *appService) GetTimelinePostsByAccountId(
	ctx context.Context,
	accountID string,
) ([]db.GetTimelinePostsByAccountIdRow, error) {
	actorID, err := xid.FromString(accountID)
	if err != nil {
		return []db.GetTimelinePostsByAccountIdRow{}, err
	}
	return s.store.Statuses().GetTimelinePostsByAccountId(ctx, actorID)
}
