package domain

import (
	"context"
	"database/sql"
	"log"

	"borg/internal/api"
	"borg/internal/db"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Repository[*api.User, *api.NewUser, *api.UpdateUser]
	GetFollowed(context.Context, int32) ([]*api.User, error)
	GetFollowers(context.Context, int32) ([]*api.User, error)
	RegisterUser(context.Context, *api.Login) error
	ValidateCredentials(context.Context, *api.Login) error
}

type userRepository struct {
	*db.Queries
}

var _ UserRepository = (*userRepository)(nil)

func newUserRepository(q *db.Queries) UserRepository {
	return &userRepository{q}
}

func (r *userRepository) Create(ctx context.Context, user *api.NewUser) error {
	return r.AddUser(ctx, *addUserToDB(user))
}

func (r *userRepository) GetByID(ctx context.Context, id int32) (*api.User, error) {
	u, err := r.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return userToAPI(&u), nil
}

func (r *userRepository) Update(ctx context.Context, user *api.UpdateUser) error {
	return r.UpdateUser(ctx, *updateUserToDB(user))
}

func (r *userRepository) Delete(ctx context.Context, id int32) error {
	return r.DeleteUser(ctx, id)
}

func (r *userRepository) GetFollowed(ctx context.Context, id int32) ([]*api.User, error) {
	users, err := r.GetFollowedUsers(ctx, id)
	if err != nil {
		return nil, err
	}
	var res []*api.User
	for _, i := range users {
		res = append(res, userToAPI(&i))
	}
	return res, nil
}

func (r *userRepository) GetFollowers(ctx context.Context, id int32) ([]*api.User, error) {
	users, err := r.GetFollowingUsers(ctx, id)
	if err != nil {
		return nil, err
	}
	var res []*api.User
	for _, i := range users {
		res = append(res, userToAPI(&i))
	}
	return res, nil
}

func (r *userRepository) RegisterUser(ctx context.Context, credentials *api.Login) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	log.Println(string(hash))
	user := db.AddUserParams{Username: credentials.Username, PasswordHash: string(hash)}
	if err := r.AddUser(ctx, user); err != nil {
		return err
	}
	return nil
}

func (r *userRepository) ValidateCredentials(ctx context.Context, credentials *api.Login) error {
	user, err := r.GetUserByUsername(ctx, credentials.Username)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(credentials.Password)); err != nil {
		return err
	}
	return nil
}

func userToAPI(u *db.User) *api.User {
	return &api.User{
		Bio:            u.Bio.String,
		CreatedAt:      u.CreatedAt.Time,
		FollowersCount: int(u.FollowersCount.Int32),
		FollowingCount: int(u.FollowingCount.Int32),
		Id:             int(u.ID),
		IsAdmin:        u.IsAdmin.Bool,
		Origin:         u.Origin.String,
		UpdatedAt:      u.UpdatedAt.Time,
		Username:       u.Username,
	}
}

func addUserToDB(u *api.NewUser) *db.AddUserParams {
	return &db.AddUserParams{
		Username:       u.Username,
		PasswordHash:   "",
		Bio:            sql.NullString{},
		FollowersCount: sql.NullInt32{},
		FollowingCount: sql.NullInt32{},
		IsAdmin:        sql.NullBool{},
	}
}

func updateUserToDB(u *api.UpdateUser) *db.UpdateUserParams {
	var bio sql.NullString
	var isAdmin sql.NullBool
	if u.Bio != nil {
		bio = sql.NullString{
			String: *u.Bio,
			Valid:  true,
		}
	}
	if u.IsAdmin != nil {
		isAdmin = sql.NullBool{
			Bool:  *u.IsAdmin,
			Valid: true,
		}
	}
	return &db.UpdateUserParams{
		ID:             0,
		PasswordHash:   "",
		Bio:            bio,
		FollowersCount: sql.NullInt32{},
		FollowingCount: sql.NullInt32{},
		IsAdmin:        isAdmin,
	}
}
