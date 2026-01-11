package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/config"
	"github.com/kibirisu/borg/internal/db"
	repo "github.com/kibirisu/borg/internal/repository"
	"github.com/kibirisu/borg/internal/util"
)

type AppService interface {
	Register(context.Context, api.AuthForm) error
	Login(context.Context, api.AuthForm) (string, error)
	GetAccountFollowers(context.Context, int) ([]db.Account, error)
	GetAccountFollowing(context.Context, int) ([]db.Account, error)
	GetLocalAccount(context.Context, string) (*db.Account, error)
	AddRemoteAccount(ctx context.Context, remote *db.CreateActorParams) (*db.Account, error)
	CreateFollow(ctx context.Context, follow *db.CreateFollowParams) (*db.Follow, error)
	AddNote(context.Context, db.CreateStatusParams) (db.Status, error)
	AddFavourite(context.Context, int, int) (db.Favourite, error)
	FollowAccount(context.Context, int, int) (*db.Follow, error)
	GetAccountById(context.Context, int) (db.Account, error)
	GetAccount(context.Context, db.GetAccountParams) (*db.Account, error)
	GetLocalPosts(context.Context) ([]db.GetLocalStatusesRow, error)
	GetPostByAccountId(context.Context, int) ([]db.GetStatusesByAccountIdRow, error)
	GetPostById(context.Context, int) (*db.Status, error)
	GetPostLikes(context.Context, int) ([]db.Favourite, error)
	GetPostShares(context.Context, int) ([]db.Status, error)
	GetPostByIdWithMetadata(context.Context, int) (*db.GetStatusByIdWithMetadataRow, error)
	UpdatePost(context.Context, db.UpdateStatusParams) (db.Status, error)
	GetPostComments(context.Context, int) ([]db.Status, error)
	UpdateAccount(ctx context.Context, params db.UpdateAccountParams) (db.Account, error)
	DeletePost(context.Context, int) error
	GetLikedPostsByUser(context.Context, int) ([]db.GetLikedPostsByAccountIdRow, error)
	GetSharedPostsByUser(context.Context, int) ([]db.GetSharedPostsByAccountIdRow, error)
	GetTimelinePosts(context.Context, int) ([]db.GetTimelinePostsByAccountIdRow, error)
	UnfollowAccount(context.Context, int, int) error
	// EW, idk if this should stay here
	DeliverToFollowers(http.ResponseWriter, *http.Request, int, func(recipientURI string) any)
}

type appService struct {
	store repo.Store
	conf  *config.Config
}

var _ AppService = (*appService)(nil)

// Register implements AppService.
func (s *appService) Register(ctx context.Context, form api.AuthForm) error {
	uri := fmt.Sprintf("http://%s/users/%s", s.conf.ListenHost, form.Username)
	log.Printf("register: creating actor username=%s uri=%s", form.Username, uri)
	actor, err := s.store.Accounts().Create(ctx, db.CreateActorParams{
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

// GetAccountById implements AppService.
func (s *appService) GetAccountById(
	ctx context.Context, accountID int,
) (db.Account, error) {
	return s.store.Accounts().GetById(ctx, accountID)
}

// AddFavourite implements AppService.
func (s *appService) AddFavourite(
	ctx context.Context, accountID int, postID int,
) (db.Favourite, error) {
	params := db.CreateFavouriteParams{
		AccountID: int32(accountID),
		StatusID:  int32(postID),
		Uri:       fmt.Sprintf("http://%s/likes/%s", s.conf.ListenHost, uuid.NewString()),
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

func (s *appService) GetPostByIdWithMetadata(
	ctx context.Context,
	id int,
) (*db.GetStatusByIdWithMetadataRow, error) {
	status, err := s.store.Statuses().GetByIdWithMetadata(ctx, id)
	if err != nil {
		return nil, err
	} else {
		return &status, nil
	}
}

func (s *appService) GetPostById(ctx context.Context, id int) (*db.Status, error) {
	status, err := s.store.Statuses().GetById(ctx, id)
	if err != nil {
		return nil, err
	} else {
		return &status, nil
	}
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
		Uri:             fmt.Sprintf("http://%s/follows/%s", s.conf.ListenHost, uuid.NewString()),
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

func (s *appService) GetPostByAccountId(
	ctx context.Context,
	id int,
) ([]db.GetStatusesByAccountIdRow, error) {
	return s.store.Accounts().GetPosts(ctx, id)
}

func (s *appService) GetLocalPosts(ctx context.Context) ([]db.GetLocalStatusesRow, error) {
	return s.store.Statuses().GetLocalStatuses(ctx)
}

// GetPostComments implements AppService.
func (s *appService) GetPostComments(ctx context.Context, id int) ([]db.Status, error) {
	return s.store.Statuses().GetPostComments(ctx, id)
}

// UpdatePost implements AppService.
func (s *appService) UpdatePost(
	ctx context.Context,
	params db.UpdateStatusParams,
) (db.Status, error) {
	return s.store.Statuses().Update(ctx, params)
}

// UpdateAccount implements AppService.
func (s *appService) UpdateAccount(
	ctx context.Context,
	params db.UpdateAccountParams,
) (db.Account, error) {
	return s.store.Accounts().UpdateAccount(ctx, params)
}

// DeletePost implements AppService.
func (s *appService) DeletePost(ctx context.Context, id int) error {
	return s.store.Statuses().Delete(ctx, id)
}

// GetLikedPostsByUser implements AppService.
func (s *appService) GetLikedPostsByUser(
	ctx context.Context,
	accountID int,
) ([]db.GetLikedPostsByAccountIdRow, error) {
	return s.store.Favourites().GetLikedPostsByUser(ctx, accountID)
}

// GetSharedPostsByUser implements AppService.
func (s *appService) GetSharedPostsByUser(
	ctx context.Context,
	accountID int,
) ([]db.GetSharedPostsByAccountIdRow, error) {
	return s.store.Statuses().GetSharedPostsByUser(ctx, accountID)
}

// GetTimelinePosts implements AppService.
func (s *appService) GetTimelinePosts(
	ctx context.Context,
	accountID int,
) ([]db.GetTimelinePostsByAccountIdRow, error) {
	return s.store.Accounts().GetTimelinePosts(ctx, accountID)
}
// UnfollowAccount implements AppService.
func (s *appService) UnfollowAccount(
	ctx context.Context,
	follower int,
	followee int,
) error {
	return s.store.Follows().Delete(ctx, db.DeleteFollowParams{
		AccountID:       int32(follower),
		TargetAccountID: int32(followee),
	})
}

