package domain

import (
	"context"

	"github.com/kibirisu/borg/internal/api"
	"github.com/kibirisu/borg/internal/db"
)

type UserRepository interface {
	Repository[*api.User, *api.NewUser, *api.UpdateUser]
	GetFollowed(context.Context, int32) ([]*api.User, error)
	GetFollowers(context.Context, int32) ([]*api.User, error)
}

type userRepository struct {
	*db.Queries
}

var _ UserRepository = (*userRepository)(nil)

func (r *userRepository) Create(ctx context.Context, user *api.NewUser) error {
	panic("unimplemented")
}

func (r *userRepository) GetByID(ctx context.Context, id int32) (*api.User, error) {
	panic("unimplemented")
}

func (r *userRepository) Update(ctx context.Context, user *api.UpdateUser) error {
	panic("unimplemented")
}

func (r *userRepository) Delete(ctx context.Context, id int32) error {
	panic("unimplemented")
}

func (r *userRepository) GetFollowed(ctx context.Context, id int32) ([]*api.User, error) {
	panic("unimplemented")
}

func (r *userRepository) GetFollowers(ctx context.Context, id int32) ([]*api.User, error) {
	panic("unimplemented")
}

func userToAPI(u *db.User) *api.User {
	return &api.User{}
}

// func addUserToDB(u *api.NewUser) *db.AddUserParams {
// 	return &db.AddUserParams{
// 		Username:       u.Username,
// 		PasswordHash:   "",
// 		Bio:            sql.NullString{},
// 		FollowersCount: sql.NullInt32{},
// 		FollowingCount: sql.NullInt32{},
// 		IsAdmin:        sql.NullBool{},
// 	}
// }

// func updateUserToDB(u *api.UpdateUser) *db.UpdateUserParams {
// 	var bio sql.NullString
// 	var isAdmin sql.NullBool
// 	if u.Bio != nil {
// 		bio = sql.NullString{
// 			String: *u.Bio,
// 			Valid:  true,
// 		}
// 	}
// 	if u.IsAdmin != nil {
// 		isAdmin = sql.NullBool{
// 			Bool:  *u.IsAdmin,
// 			Valid: true,
// 		}
// 	}
// 	return &db.UpdateUserParams{
// 		ID:             0,
// 		PasswordHash:   "",
// 		Bio:            bio,
// 		FollowersCount: sql.NullInt32{},
// 		FollowingCount: sql.NullInt32{},
// 		IsAdmin:        isAdmin,
// 	}
// }
