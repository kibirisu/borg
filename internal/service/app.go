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
	"github.com/kibirisu/borg/internal/server/mapper"
	"github.com/kibirisu/borg/internal/util"
	"github.com/kibirisu/borg/internal/worker"
)

type AppService interface {
	WebfingerAccount(context.Context, api.GetWellKnownWebfingerParams) (*api.Webfinger, error)
	Register(context.Context, api.AuthForm) error
	Login(context.Context, api.AuthForm) (string, error)
	GetAccount(context.Context, string) (*api.Account, error)
	GetAccountStatuses(context.Context, string) ([]api.Status, error)
	GetAccountFollowers(context.Context, string) ([]api.Account, error)
	GetAccountFollowing(context.Context, string) ([]api.Account, error)
	FollowAccount(context.Context, string) (worker.Job, error)
	UnfollowAccount(context.Context, string) (worker.Job, error)
	LookupAccount(context.Context, string) (*api.Account, error)
	CreateStatus(context.Context, api.PostApiStatusesJSONBody) (worker.Job, error)
	ViewStatus(context.Context, string) (*api.Status, error)
	GetStatusReplies(context.Context, string) ([]api.Status, error)
	FavouriteStatus(context.Context, string) (worker.Job, error)
	UnfavouriteStatus(context.Context, string) (worker.Job, error)
	ReblogStatus(context.Context, string) (worker.Job, error)
	UnreblogStatus(context.Context, string) (worker.Job, error)
	ViewHomeTimeline(context.Context) ([]api.Status, error)
	ViewFavouriteTimeline(context.Context) ([]api.Status, error)
	ViewRebloggedTimeline(context.Context) ([]api.Status, error)
}

type appService struct {
	store    repo.Store
	prcessor proc.Processor
	conf     *config.Config
	builder  util.URIBuilder
}

var _ AppService = (*appService)(nil)

// WebfingerAccount implements AppService.
func (s *appService) WebfingerAccount(
	ctx context.Context,
	resource api.GetWellKnownWebfingerParams,
) (*api.Webfinger, error) {
	webfinger, err := s.store.Accounts().
		GetWebfinger(ctx, util.ExtractUsernameFromAcct(resource.Resource))
	if err != nil {
		return nil, err
	}
	links := api.WebfingerLinks(webfinger)
	return &api.Webfinger{
		Subject: resource.Resource,
		Links:   []api.WebfingerLinks{links},
	}, nil
}

// Register implements AppService.
func (s *appService) Register(ctx context.Context, form api.AuthForm) error {
	id := xid.New()
	actorURIs := s.builder.ActorURIs(id.String())

	log.Printf("register: creating actor username=%s uri=%s", form.Username, actorURIs.Actor)

	_, err := s.store.WithTX(ctx, func(ctx context.Context, s repo.Store) (any, error) {
		actor, err := s.Accounts().Create(ctx, db.CreateActorParams{
			ID:       id,
			Username: form.Username,
			Uri:      actorURIs.Actor,
			DisplayName: sql.NullString{
				String: form.Username,
				Valid:  true,
			},
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
			"register: user and actor created username=%s account_id=%s",
			form.Username,
			actor.ID.String(),
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
	if err = bcrypt.CompareHashAndPassword(
		[]byte(auth.PasswordHash),
		[]byte(form.Password),
	); err != nil {
		return
	}
	return issueToken(auth.ID.String(), auth.Uri, s.conf.JWTSecret)
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
	return mapper.ToAPIAccount(&account), nil
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
		// _, _ = status, idx
		s := db.GetStatusByIDRow(status)
		res[idx] = *mapper.ToAPIStatus(&s)
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
		f := db.GetAccountByIDRow(follower)
		res[idx] = *mapper.ToAPIAccount(&f)
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
		f := db.GetAccountByIDRow(followed)
		res[idx] = *mapper.ToAPIAccount(&f)
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

	var follow ap.FollowActivitier

	followReq, err := s.store.WithTX(ctx, func(ctx context.Context, store repo.Store) (any, error) {
		req, err := store.FollowRequests().Create(ctx, db.CreateFollowRequestParams{
			ID:              id,
			Uri:             s.builder.FollowRequestURI(token.ID, id.String()),
			AccountID:       followerID,
			TargetAccountID: targetAccountID,
		})
		if err != nil {
			return nil, err
		}
		if !req.Local {
			follow = ap.NewEmptyFollowActivity().WithObject(ap.Activity[ap.Actor]{
				ID:     req.Uri,
				Type:   "Follow",
				Actor:  ap.NewEmptyActor().WithLink(token.URI),
				Object: ap.NewEmptyActor().WithLink(req.TargetAccountUri),
			})
			return req, nil
		}
		err = store.Follows().Create(ctx, db.CreateFollowNewParams{
			ID:              id,
			Uri:             s.builder.FollowURI(token.ID, id.String()),
			AccountID:       req.AccountID,
			TargetAccountID: req.TargetAccountID,
		})
		return nil, err
	})
	if err != nil {
		return nil, err
	}

	if follow == nil {
		return worker.EmptyJob, nil
	}

	req := followReq.(db.CreateFollowRequestRow)

	return func(ctx context.Context) error {
		return s.prcessor.SendObject(ctx, follow.GetRaw().Object, req.TargetAccountID)
	}, nil
}

// UnfollowAccount implements AppService.
func (s *appService) UnfollowAccount(ctx context.Context, accountID string) (worker.Job, error) {
	token, ok := ctx.Value(auth.TokenContextKey).(*auth.TokenData)
	if !ok {
		return nil, errors.New("auth failure")
	}

	targetAccountID, err := xid.FromString(accountID)
	if err != nil {
		return nil, err
	}
	followerID, err := xid.FromString(token.ID)
	if err != nil {
		return nil, err
	}

	var undo ap.UndoFollowActiviter

	followReq, err := s.store.WithTX(ctx, func(ctx context.Context, store repo.Store) (any, error) {
		req, err := store.FollowRequests().
			DeleteByTargetAccountID(ctx, db.DeleteFollowRequestByAccountIDParams{
				TargetAccountID: targetAccountID,
				AccountID:       followerID,
			})
		if err != nil {
			return nil, err
		}
		if !req.Local {
			actor := ap.NewEmptyActor().WithLink(token.URI)
			undo = ap.NewEmptyUndoFollowActivity().WithObject(ap.Activity[ap.Activity[ap.Actor]]{
				ID:    "doesn't matter",
				Type:  "Undo",
				Actor: actor,
				Object: ap.NewEmptyFollowActivity().WithObject(ap.Activity[ap.Actor]{
					ID:     req.Uri,
					Type:   "Follow",
					Actor:  actor,
					Object: ap.NewEmptyActor().WithLink(req.TargetAccountUri),
				}),
			})
		}
		return req, store.Follows().DeleteByID(ctx, req.ID)
	})
	if err != nil {
		return nil, err
	}

	if undo == nil {
		return worker.EmptyJob, nil
	}

	req := followReq.(db.DeleteFollowRequestByAccountIDRow)

	return func(ctx context.Context) error {
		return s.prcessor.SendObject(ctx, undo.GetRaw().Object, req.TargetAccountID)
	}, nil
}

// LookupAccount implements AppService.
func (s *appService) LookupAccount(ctx context.Context, acct string) (*api.Account, error) {
	handle := util.ExtractHandleParts(acct)
	if s.conf.Address == handle.Domain.String {
		handle.Domain.Valid = false
	}

	queryParams := db.GetAccountByUsernameAndDomainParams(handle)
	account, err := s.store.Accounts().GetByUsernameAndDomain(ctx, queryParams)
	if err != nil {
		if !handle.Domain.Valid {
			return nil, err
		}
		account, err := s.prcessor.FetchAndStoreAccount(ctx, handle.Username, handle.Domain.String)
		if err != nil {
			return nil, err
		}
		a := &db.GetAccountByIDRow{
			Account:        account,
			Acct:           acct,
			FollowersCount: 0, // we don't check for followers/following collections of freshly fetched actors
			FollowingCount: 0,
		}
		return mapper.ToAPIAccount(a), nil
	}

	a := db.GetAccountByIDRow(account)
	return mapper.ToAPIAccount(&a), nil
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
	if status.InReplyToId != nil {
		id, err := xid.FromString(*status.InReplyToId)
		if err != nil {
			return nil, err
		}
		inReplyToID = &id
	}

	createdStatus, err := s.store.Statuses().Create(ctx, db.CreateStatusParams{
		ID:  statusID,
		Uri: statusURIs.Status,
		Url: "not needed rn",
		Content: sql.NullString{
			String: status.Status,
			Valid:  true,
		},
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
			Replies: ap.NewEmptyNoteCollection().WithObject(ap.Collection[ap.Note]{
				ID:   statusURIs.Replies,
				Type: "Collection",
				First: ap.NewEmptyNoteCollectionPage().WithObject(ap.CollectionPage[ap.Note]{
					ID:     "None",
					Type:   "CollectionPage",
					Next:   ap.NewEmptyNoteCollectionPage(),
					PartOf: ap.NewEmptyNoteCollection().WithLink(statusURIs.Replies),
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

	status, err := s.store.Statuses().GetByID(ctx, db.GetStatusByIDParams{
		ID:        statusID,
		AccountID: accountID,
	})
	if err != nil {
		return nil, err
	}

	return mapper.ToAPIStatus(&status), nil
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

	favourite, err := s.store.Favourites().Create(ctx, db.CreateFavouriteParams{
		ID:         id,
		AccountID:  accountID,
		AccountUri: token.URI,
		StatusID:   statusID,
		Uri:        s.builder.LikeRequestURI(token.ID, id.String()),
	})
	if err != nil {
		return nil, err
	}

	if favourite.Local.Bool {
		return worker.EmptyJob, nil
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
func (s *appService) ReblogStatus(ctx context.Context, statusID string) (worker.Job, error) {
	token, ok := ctx.Value(auth.TokenContextKey).(*auth.TokenData)
	if !ok {
		return nil, errors.New("auth failure")
	}

	id := xid.New()
	reblogOfID, err := xid.FromString(statusID)
	if err != nil {
		return nil, err
	}
	accountID, err := xid.FromString(token.ID)
	if err != nil {
		return nil, err
	}

	reblog, err := s.store.Statuses().ReblogStatus(ctx, db.CreateReblogParams{
		ID:         id,
		Uri:        s.builder.AnnounceURI(token.ID, id.String()),
		Url:        "no",
		AccountID:  accountID,
		AccountUri: token.URI,
		ReblogOfID: &reblogOfID,
	})
	if err != nil {
		return nil, err
	}

	announce := ap.NewEmptyAnnounceActivity().WithObject(ap.Activity[ap.Note]{
		ID:     reblog.Uri,
		Type:   "Announce",
		Actor:  ap.NewEmptyActor().WithLink(reblog.AccountUri),
		Object: ap.NewEmptyNote().WithLink(reblog.ReblogOfUri.String),
	})

	return func(ctx context.Context) error {
		return s.prcessor.DistributeObject(ctx, announce.GetRaw().Object, accountID)
	}, nil
}

// UnfavouriteStatus implements AppService.
func (s *appService) UnfavouriteStatus(ctx context.Context, id string) (worker.Job, error) {
	token, ok := ctx.Value(auth.TokenContextKey).(*auth.TokenData)
	if !ok {
		return nil, errors.New("auth failure")
	}
	loggedInID, err := xid.FromString(token.ID)
	if err != nil {
		return nil, err
	}
	statusID, err := xid.FromString(id)
	if err != nil {
		return nil, err
	}
	fav, err := s.store.Favourites().DeleteByStatusID(ctx, loggedInID, statusID)
	if err != nil {
		return nil, err
	}

	actor := ap.NewEmptyActor().WithLink(fav.AccountUri)
	undo := ap.NewEmptyUndoActivity().WithObject(ap.Activity[ap.Activity[ap.Note]]{
		ID:    "nope",
		Type:  "Undo",
		Actor: actor,
		Object: ap.NewEmptyLikeActivity().WithObject(ap.Activity[ap.Note]{
			ID:     fav.Uri,
			Type:   "Like",
			Actor:  actor,
			Object: ap.NewEmptyNote().WithLink(fav.StatusUri),
		}),
	})

	return func(ctx context.Context) error {
		return s.prcessor.SendObject(ctx, undo.GetRaw().Object, fav.TargetAccountID)
	}, nil
}

// UnreblogStatus implements AppService.
func (s *appService) UnreblogStatus(ctx context.Context, id string) (worker.Job, error) {
	token, ok := ctx.Value(auth.TokenContextKey).(*auth.TokenData)
	if !ok {
		return nil, errors.New("auth failure")
	}

	accountID, err := xid.FromString(token.ID)
	if err != nil {
		return nil, err
	}
	statusID, err := xid.FromString(id)
	if err != nil {
		return nil, err
	}

	status, err := s.store.Statuses().DeleteReblogByStatusID(ctx, &statusID, accountID)
	if err != nil {
		return nil, err
	}

	actor := ap.NewEmptyActor().WithLink(status.AccountUri)
	undo := ap.NewEmptyUndoActivity().WithObject(ap.Activity[ap.Activity[ap.Note]]{
		ID:    "does it matter?",
		Type:  "Undo",
		Actor: actor,
		Object: ap.NewEmptyAnnounceActivity().WithObject(ap.Activity[ap.Note]{
			ID:     status.Uri,
			Type:   "Announce",
			Actor:  actor,
			Object: ap.NewEmptyNote().WithLink(status.ReblogOfUri.String),
		}),
	})

	return func(ctx context.Context) error {
		return s.prcessor.SendObject(ctx, undo.GetRaw().Object, *status.ReblogOfAccountID)
	}, nil
}

// GetStatusReplies implements AppService.
func (s *appService) GetStatusReplies(ctx context.Context, id string) ([]api.Status, error) {
	token, ok := ctx.Value(auth.TokenContextKey).(*auth.TokenData)
	if !ok {
		return nil, errors.New("auth failure")
	}
	loggedInID, err := xid.FromString(token.ID)
	if err != nil {
		return nil, err
	}
	statusID, err := xid.FromString(id)
	if err != nil {
		return nil, err
	}
	statuses, err := s.store.Statuses().GetReplies(ctx, loggedInID, statusID)
	if err != nil {
		return nil, err
	}
	res := make([]api.Status, len(statuses))
	for idx, status := range statuses {
		s := db.GetStatusByIDRow(status)
		res[idx] = *mapper.ToAPIStatus(&s)
	}
	return res, nil
}

// ViewHomeTimeline implements AppService.
func (s *appService) ViewHomeTimeline(ctx context.Context) ([]api.Status, error) {
	token, ok := ctx.Value(auth.TokenContextKey).(*auth.TokenData)
	if !ok {
		return nil, errors.New("auth failure")
	}
	loggedInID, err := xid.FromString(token.ID)
	if err != nil {
		return nil, err
	}
	statuses, err := s.store.Statuses().GetHomeTimelineByAccountID(ctx, loggedInID)
	if err != nil {
		return nil, err
	}

	res := make([]api.Status, len(statuses))
	for idx, status := range statuses {
		s := db.GetStatusByIDRow(status)
		res[idx] = *mapper.ToAPIStatus(&s)
	}
	return res, nil
}

// ViewFavouriteTimeline implements AppService.
func (s *appService) ViewFavouriteTimeline(ctx context.Context) ([]api.Status, error) {
	token, ok := ctx.Value(auth.TokenContextKey).(*auth.TokenData)
	if !ok {
		return nil, errors.New("auth failure")
	}
	loggedInID, err := xid.FromString(token.ID)
	if err != nil {
		return nil, err
	}
	statuses, err := s.store.Statuses().GetFavouriteByAccountID(ctx, loggedInID)
	if err != nil {
		return nil, err
	}

	res := make([]api.Status, len(statuses))
	for idx, status := range statuses {
		s := db.GetStatusByIDRow(status)
		res[idx] = *mapper.ToAPIStatus(&s)
	}
	return res, nil
}

// ViewRebloggedTimeline implements AppService.
func (s *appService) ViewRebloggedTimeline(ctx context.Context) ([]api.Status, error) {
	token, ok := ctx.Value(auth.TokenContextKey).(*auth.TokenData)
	if !ok {
		return nil, errors.New("auth failure")
	}
	loggedInID, err := xid.FromString(token.ID)
	if err != nil {
		return nil, err
	}
	statuses, err := s.store.Statuses().GetRebloggedByAccountID(ctx, loggedInID)
	if err != nil {
		return nil, err
	}

	res := make([]api.Status, len(statuses))
	for idx, status := range statuses {
		s := db.GetStatusByIDRow(status)
		res[idx] = *mapper.ToAPIStatus(&s)
	}
	return res, nil
}
