package service

import (
	"context"
	"database/sql"
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
}

type appService struct {
	accounts repo.AccountRepository
	users    repo.UserRepository
	conf     *config.Config
}

var _ AppService = (*appService)(nil)

func NewAppService(
	accounts repo.AccountRepository,
	users repo.UserRepository,
	conf *config.Config,
) AppService {
	return &appService{accounts, users, conf}
}

// Register implements AppService.
func (s *appService) Register(ctx context.Context, form api.AuthForm) error {
	// uri := "http://" + s.conf.ListenHost + "/users/" + form.Username
	uri := fmt.Sprintf("http://%s/users/%s", s.conf.ListenHost, form.Username)
	actor, err := s.accounts.Create(ctx, db.CreateActorParams{
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
	if err = s.users.Create(ctx, db.CreateUserParams{
		AccountID:    actor.ID,
		PasswordHash: string(hash),
	}); err != nil {
		return err
	}
	return nil
}

// Login implements AppService.
func (s *appService) Login(ctx context.Context, form api.AuthForm) (string, error) {
	auth, err := s.users.GetByUsername(ctx, form.Username)
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
