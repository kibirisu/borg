package domain

import (
	"context"
	"database/sql"
	"time"

	"borg/internal/api"
	"borg/internal/db"
)

type UserRepository interface {
	Repository[*api.User, *api.NewUser, *api.UpdateUser]
	GetFollowed(context.Context, int32) ([]*api.User, error)
	GetFollowers(context.Context, int32) ([]*api.User, error)
	GetByUsername(context.Context, string) (*api.User, error)
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

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*api.User, error) {
	u, err := r.GetUserByUsername(ctx, username)
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

func userToAPI(u *db.User) *api.User {
	bio := ""
	if u.Bio.Valid {
		bio = u.Bio.String
	}
	origin := ""
	if u.Origin.Valid {
		origin = u.Origin.String
	}
	followersCount := 0
	if u.FollowersCount.Valid {
		followersCount = int(u.FollowersCount.Int32)
	}
	followingCount := 0
	if u.FollowingCount.Valid {
		followingCount = int(u.FollowingCount.Int32)
	}
	isAdmin := false
	if u.IsAdmin.Valid {
		isAdmin = u.IsAdmin.Bool
	}
	createdAt := u.CreatedAt.Time
	if !u.CreatedAt.Valid {
		createdAt = time.Time{}
	}
	updatedAt := u.UpdatedAt.Time
	if !u.UpdatedAt.Valid {
		updatedAt = time.Time{}
	}
	return &api.User{
		Bio:            bio,
		CreatedAt:      createdAt,
		FollowersCount: followersCount,
		FollowingCount: followingCount,
		Id:             int(u.ID),
		IsAdmin:        isAdmin,
		Origin:         origin,
		UpdatedAt:      updatedAt,
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
