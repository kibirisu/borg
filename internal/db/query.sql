-- name: GetActor :one
SELECT * FROM accounts WHERE username = $1;

-- name: AddUser :exec
INSERT INTO users (
  username,
  password_hash,
  bio,
  followers_count,
  following_count,
  is_admin,
  origin
) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: GetFollowedUsers :many
SELECT u.* FROM users u JOIN followers f ON u.id = f.following_id WHERE f.follower_id = $1;

-- name: GetFollowingUsers :many
SELECT u.* FROM users u JOIN followers f ON u.id = f.follower_id WHERE f.following_id = $1;

-- name: UpdateUser :exec
UPDATE users SET password_hash = $2, bio = $3, followers_count = $4, following_count = $5, is_admin = $6 WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;


-- name: AddPost :exec
INSERT INTO posts (user_id, content) VALUES ($1, $2);

-- name: GetPostsByUserID :many
SELECT * FROM posts WHERE user_id = $1;

-- name: GetPost :one
SELECT * FROM posts WHERE id = $1;

-- name: UpdatePost :exec
UPDATE posts SET content = $2, like_count = $3, share_count = $4, comment_count = $5 WHERE id = $1;

-- name: DeletePost :exec
DELETE FROM posts WHERE id = $1;


-- name: AddComment :exec
INSERT INTO comments (post_id, user_id, content, parent_id) VALUES ($1, $2, $3, $4);

SELECT * FROM comments WHERE id = $1;

-- name: GetPostComments :many
SELECT c.* FROM comments c JOIN posts p ON c.post_id = p.id WHERE p.id = $1;

-- name: GetUserComments :many
SELECT c.* FROM comments c JOIN users u ON c.user_id = u.id WHERE u.id = $1;

-- name: DeleteComment :exec
DELETE FROM comments WHERE id = $1;


-- name: AddLike :exec
INSERT INTO likes (post_id, user_id) VALUES ($1, $2);

-- name: GetLikeByID :one
SELECT * FROM likes WHERE id = $1;

-- name: GetLikesByPostID :many
SELECT * FROM likes WHERE post_id = $1;

-- name: GetLikesByUserID :many
SELECT * FROM likes WHERE user_id = $1;

-- name: DeleteLike :exec
DELETE FROM likes WHERE id = $1;


-- name: AddShare :exec
INSERT INTO shares (post_id, user_id) VALUES ($1, $2);

-- name: GetShareByID :one
SELECT * FROM shares WHERE id = $1;

-- name: GetSharesByPostID :many
SELECT * FROM shares WHERE post_id = $1;

-- name: GetShareByUserID :many
SELECT * FROM shares WHERE user_id = $1;

-- name: DeleteShare :exec
DELETE FROM shares WHERE id = $1;


-- queries that are needed for frontend
-- name: GetPostsByOrigin :many
SELECT p.* FROM posts p JOIN users u ON p.user_id = u.id WHERE u.origin = $1;

-- name: GetAllUsers :many
SELECT * FROM users;

-- name: GetAllPosts :many
SELECT p.*, u.username FROM posts p JOIN users u ON p.user_id = u.id ORDER BY p.created_at DESC;
