package service

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

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
	GetAccount(context.Context, string) (*api.Account, error)
	GetAccountStatuses(context.Context, string) ([]api.Status, error)
	GetAccountFollowers(context.Context, string) ([]api.Account, error)
	GetAccountFollowing(context.Context, string) ([]api.Account, error)
	FollowAccount(context.Context, string) (worker.Job, error)
	UnfollowAccount(context.Context, string) (worker.Job, error)
	CreateStatus(context.Context, api.PostApiStatusesJSONBody) (worker.Job, error)
	ViewStatus(context.Context, string) (*api.Status, error)
	FavouriteStatus(context.Context, string) (worker.Job, error)
	UnfavouriteStatus(context.Context, string) (worker.Job, error)
	ReblogStatus(context.Context, string) (worker.Job, error)
	UnreblogStatus(context.Context, string) (worker.Job, error)
	GetPostByIDWithMetadata(context.Context, string) (*db.GetStatusByIdWithMetadataRow, error)
	GetLikedPostsByAccountId(context.Context, string) ([]db.GetLikedPostsByAccountIdRow, error)
	GetSharedPostsByAccountId(context.Context, string) ([]db.GetSharedPostsByAccountIdRow, error)
	GetTimelinePostsByAccountId(
		context.Context,
		string,
	) ([]db.GetTimelinePostsByAccountIdRow, error)
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
	id := xid.New()
	actorURIs := s.builder.ActorURIs(id.String())

	log.Printf("register: creating actor username=%s uri=%s", form.Username, actorURIs.Actor)

	_, err := s.store.WithTX(ctx, func(ctx context.Context, s repo.Store) (any, error) {
		actor, err := s.Accounts().Create(ctx, db.CreateActorParams{
			ID:           id,
			Username:     form.Username,
			Uri:          actorURIs.Actor,
			DisplayName:  sql.NullString{},
			InboxUri:     actorURIs.Inbox,
			OutboxUri:    actorURIs.Outbox,
			Url:          "not gonna use that rn",
			Domain:       sql.NullString{},
			FollowersUri: actorURIs.Followers,
			FollowingUri: actorURIs.Following,
		})
		if err != nil {
			return nil, err
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		if err = s.Users().Create(ctx, db.CreateUserParams{
			ID:           xid.New(),
			AccountID:    actor.ID,
			PasswordHash: string(hash),
		}); err != nil {
			return nil, err
		}

		log.Printf(
			"register: user and actor created username=%s account_id=%d",
			form.Username,
			actor.ID,
		)
		return nil, nil
	})
	return err
}

// Login implements AppService.
func (s *appService) Login(ctx context.Context, form api.AuthForm) (token string, err error) {
	auth, err := s.store.Users().GetByUsername(ctx, form.Username)
	if err != nil {
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(auth.PasswordHash), []byte(form.Password)); err != nil {
		return
	}
	token, err = issueToken(auth.ID.String(), form.Username, s.conf.JWTSecret)
	return
}

// GetAccount implements AppService.
func (s *appService) GetAccount(ctx context.Context, id string) (*api.Account, error) {
	accountID, err := xid.FromString(id)
	if err != nil {
		return nil, err
	}
	account, err := s.store.Accounts().GetByID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	res := api.Account{
		Acct:           account.Acct,
		DisplayName:    account.Account.DisplayName.String,
		FollowersCount: int(account.FollowersCount),
		FollowingCount: int(account.FollowingCount),
		Id:             id,
		Url:            account.Account.Url,
		Username:       account.Account.Username,
	}
	return &res, nil
}

// GetAccountStatuses implements AppService.
func (s *appService) GetAccountStatuses(ctx context.Context, id string) ([]api.Status, error) {
	token, ok := ctx.Value(auth.TokenContextKey).(*auth.TokenData)
	if !ok {
		return nil, errors.New("auth failure")
	}
	accountID, err := xid.FromString(id)
	if err != nil {
		return nil, err
	}
	loggedInID, err := xid.FromString(token.ID)
	if err != nil {
		return nil, err
	}
	statuses, err := s.store.Statuses().GetByAccountID(ctx, db.GetStatusesByAccountIDParams{
		LoggedInID: loggedInID,
		AccountID:  accountID,
	})
	if err != nil {
		return nil, err
	}

	res := make([]api.Status, len(statuses))
	for idx, status := range statuses {
		var inReplyToID, inReplyToAccountID *string
		if status.Status.InReplyToID != nil {
			id := status.Status.InReplyToID.String()
			inReplyToID = &id
		}
		if status.Status.InReplyToAccountID != nil {
			id := status.Status.InReplyToID.String()
			inReplyToAccountID = &id
		}
		res[idx] = api.Status{
			Content:            status.Status.Content,
			Favourited:         status.Favourited,
			FavouritesCount:    int(status.FavouritesCount),
			Id:                 status.Status.ID.String(),
			InReplyToAccountId: inReplyToAccountID,
			InReplyToId:        inReplyToID,
			Reblogged:          status.Reblogged,
			ReblogsCount:       int(status.FavouritesCount),
			RepliesCount:       int(status.FavouritesCount),
			Uri:                status.Status.Uri,
		}
	}
	return res, nil
}

// GetAccountFollowers implements AppService.
func (s *appService) GetAccountFollowers(ctx context.Context, id string) ([]api.Account, error) {
	accountID, err := xid.FromString(id)
	if err != nil {
		return nil, err
	}
	followers, err := s.store.Accounts().GetFollowersByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	res := make([]api.Account, len(followers))
	for idx, follower := range followers {
		res[idx] = api.Account{
			Acct:           follower.Acct,
			DisplayName:    follower.Account.DisplayName.String,
			FollowersCount: int(follower.FollowersCount),
			FollowingCount: int(follower.FollowingCount),
			Id:             follower.Account.ID.String(),
			Url:            follower.Account.Url,
			Username:       follower.Account.Username,
		}
	}
	return res, nil
}

// GetAccountFollowing implements AppService.
func (s *appService) GetAccountFollowing(ctx context.Context, id string) ([]api.Account, error) {
	accountID, err := xid.FromString(id)
	if err != nil {
		return nil, err
	}
	following, err := s.store.Accounts().GetFollowingByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}
	res := make([]api.Account, len(following))
	for idx, followed := range following {
		res[idx] = api.Account{
			Acct:           followed.Acct,
			DisplayName:    followed.Account.DisplayName.String,
			FollowersCount: int(followed.FollowersCount),
			FollowingCount: int(followed.FollowingCount),
			Id:             followed.Account.ID.String(),
			Url:            followed.Account.Url,
			Username:       followed.Account.Username,
		}
	}
	return res, nil
}

// FollowAccount implements AppService.
func (s *appService) FollowAccount(ctx context.Context, accountID string) (worker.Job, error) {
	token, ok := ctx.Value(auth.TokenContextKey).(*auth.TokenData)
	if !ok {
		return nil, errors.New("auth failure")
	}
	id := xid.New()
	followerID, err := xid.FromString(token.ID)
	if err != nil {
		return nil, err
	}
	targetAccountID, err := xid.FromString(accountID)
	if err != nil {
		return nil, err
	}

	req, err := s.store.FollowRequests().Create(ctx, db.CreateFollowRequestParams{
		ID:              id,
		Uri:             s.builder.FollowRequestURI(token.ID, id.String()),
		AccountID:       followerID,
		TargetAccountID: targetAccountID,
	})
	if err != nil {
		return nil, err
	}

	follow := ap.NewEmptyFollowActivity().WithObject(ap.Activity[ap.Actor]{
		ID:     req.Uri,
		Type:   "Follow",
		Actor:  ap.NewEmptyActor().WithLink(token.URI),
		Object: ap.NewEmptyActor().WithLink(req.TargetAccountUri),
	})

	return func(ctx context.Context) error {
		return s.prcessor.SendObject(ctx, follow.GetRaw().Object, req.AccountID)
	}, nil
}

// UnfollowAccount implements AppService.
func (s *appService) UnfollowAccount(ctx context.Context, accountID string) (worker.Job, error) {
	token, ok := ctx.Value(auth.TokenContextKey).(*auth.TokenData)
	if !ok {
		return nil, errors.New("auth failure")
	}
	_ = token
	panic("unimplemented")
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
	if status.Status.InReplyToAccountID != nil {
		id := status.Status.InReplyToID.String()
		inReplyToAccountID = &id
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

// FavouriteStatus implements AppService.
func (s *appService) FavouriteStatus(ctx context.Context, favouritedID string) (worker.Job, error) {
	token, ok := ctx.Value(auth.TokenContextKey).(*auth.TokenData)
	if !ok {
		return nil, errors.New("auth failure")
	}

	id := xid.New()
	statusID, err := xid.FromString(favouritedID)
	if err != nil {
		return nil, err
	}
	accountID, err := xid.FromString(token.ID)
	if err != nil {
		return nil, err
	}

	favourite, err := s.store.Favourites().CreateNew(ctx, db.CreateFavouriteNewParams{
		ID:        id,
		AccountID: accountID,
		StatusID:  statusID,
		Uri:       s.builder.LikeRequestURI(token.ID, id.String()),
	})
	if err != nil {
		return nil, err
	}

	like := ap.NewEmptyLikeActivity().WithObject(ap.Activity[ap.Note]{
		ID:     favourite.Uri,
		Type:   "Like",
		Actor:  ap.NewEmptyActor().WithLink(token.URI),
		Object: ap.NewEmptyNote().WithLink(favourite.StatusUri),
	})

	return func(ctx context.Context) error {
		return s.prcessor.SendObject(ctx, like.GetRaw().Object, favourite.TargetAccountID)
	}, nil
}

// ReblogStatus implements AppService.
func (s *appService) ReblogStatus(context.Context, string) (worker.Job, error) {
	panic("unimplemented")
}

// UnfavouriteStatus implements AppService.
func (s *appService) UnfavouriteStatus(context.Context, string) (worker.Job, error) {
	panic("unimplemented")
}

// UnreblogStatus implements AppService.
func (s *appService) UnreblogStatus(context.Context, string) (worker.Job, error) {
	panic("unimplemented")
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
