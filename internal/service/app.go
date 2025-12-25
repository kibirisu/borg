package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/config"
	"github.com/kibirisu/borg/internal/db"
	repo "github.com/kibirisu/borg/internal/repository"
)

type AppService interface {
	Register(context.Context, api.AuthForm) error
	Login(context.Context, api.AuthForm) (string, error)
	GetAccountFollowers(context.Context, int) ([]db.Account, error)
	GetLocalAccount(context.Context, string) (*db.Account, error)
	AddRemoteAccount(ctx context.Context, remote *db.CreateActorParams) (*db.Account, error)
	CreateFollow(ctx context.Context, follow *db.CreateFollowParams) (*db.Follow, error) 
	AddNote(context.Context, db.CreateStatusParams) (db.Status, error)
	FollowAccount(context.Context, int, int) (*db.Follow, error)
	GetAccountById(context.Context, int) (db.Account, error)
	GetAccount(context.Context, db.GetAccountParams) (*db.Account, error)
}

type appService struct {
	store repo.Store
	conf  *config.Config
}

var _ AppService = (*appService)(nil)

func NewAppService(
	store repo.Store,
	conf *config.Config,
) AppService {
	return &appService{store, conf}
}

// Register implements AppService.
func (s *appService) Register(ctx context.Context, form api.AuthForm) error {
	uri := fmt.Sprintf("http://%s/users/%s", s.conf.ListenHost, form.Username)
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
		return err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	if err = s.store.Users().Create(ctx, db.CreateUserParams{
		AccountID:    actor.ID,
		PasswordHash: string(hash),
	}); err != nil {
		return err
	}
	return nil
}

// Login implements AppService.
func (s *appService) Login(ctx context.Context, form api.AuthForm) (string, error) {
	auth, err := s.store.Users().GetByUsername(ctx, form.Username)
	if err != nil {
		return "", err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(form.Password), []byte(auth.PasswordHash)); err != nil {
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

func (s *appService) CreateFollow(ctx context.Context, follow *db.CreateFollowParams) (*db.Follow, error) {
	return s.store.Follows().Create(ctx, *follow)
}

// AddNote implements AppService.
func (s *appService) AddNote(ctx context.Context, note db.CreateStatusParams) (db.Status, error) {
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
	return s.store.Accounts().GetById(ctx, accountID);
}

// GetAccountFollowers implements AppService.
func (s *appService) GetAccountFollowers(
	ctx context.Context, accountID int,
) ([]db.Account, error) {
	return s.store.Accounts().GetFollowers(ctx, accountID);
}

// FollowAccount implements AppService.
func (s *appService) FollowAccount(ctx context.Context, follower int, followee int) (*db.Follow, error) {
	createParams := db.CreateFollowParams {
		Uri: "", //TODO
		AccountID: int32(follower),
		TargetAccountID: int32(followee),
	}
	return s.store.Follows().Create(ctx, createParams);
}
