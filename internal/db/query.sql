-- name: GetActor :one
SELECT * FROM accounts WHERE username = $1 AND domain IS NULL;

-- name: AuthData :one
SELECT a.id, u.password_hash FROM accounts a JOIN users u ON a.id = u.account_id WHERE a.username = $1;

-- name: CreateActor :one
INSERT INTO accounts (
    username, uri, display_name, domain, inbox_uri, outbox_uri, url, followers_uri, following_uri
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: CreateUser :exec
INSERT INTO users (
    account_id, password_hash
) VALUES (
    $1, $2
);

-- name: GetAccount :one
SELECT * FROM accounts WHERE username = $1 AND domain = $2;

-- name: GetAccountById :one
SELECT * FROM accounts WHERE id = $1;

-- name: GetLocalStatuses :many
SELECT 
    sqlc.embed(s),
    sqlc.embed(a),
    (SELECT COUNT(*) FROM favourites f WHERE f.status_id = s.id) AS like_count,
    (SELECT COUNT(*) FROM statuses r WHERE r.in_reply_to_id = s.id) AS comment_count,
    (SELECT COUNT(*) FROM statuses b WHERE b.reblog_of_id = s.id) AS share_count
FROM statuses s
JOIN accounts a ON s.account_id = a.id
WHERE a.domain is null and s.in_reply_to_id is null;

-- name: GetStatusById :one
SELECT * FROM statuses WHERE id = $1;

-- name: GetStatusByIdWithMetadata :one
SELECT 
    sqlc.embed(s),
    sqlc.embed(a),
    (SELECT COUNT(*) FROM favourites f WHERE f.status_id = s.id) AS like_count,
    (SELECT COUNT(*) FROM statuses r WHERE r.in_reply_to_id = s.id) AS comment_count,
    (SELECT COUNT(*) FROM statuses b WHERE b.reblog_of_id = s.id) AS share_count
FROM statuses s
JOIN accounts a ON s.account_id = a.id
WHERE s.id = $1;

-- name: GetStatusFavourites :many
SELECT *
FROM favourites
WHERE status_id = $1;

-- name: GetStatusShares :many
SELECT *
FROM statuses 
WHERE reblog_of_id = $1;

-- name: GetStatusesByAccountId :many
SELECT 
    sqlc.embed(s),
    sqlc.embed(a),
    (SELECT COUNT(*) FROM favourites f WHERE f.status_id = s.id) AS like_count,
    (SELECT COUNT(*) FROM statuses r WHERE r.in_reply_to_id = s.id) AS comment_count,
    (SELECT COUNT(*) FROM statuses b WHERE b.reblog_of_id = s.id) AS share_count
FROM statuses s
JOIN accounts a ON s.account_id = a.id
WHERE s.account_id = $1;

-- name: CreateFollow :one
INSERT INTO follows (
    uri, account_id, target_account_id
) VALUES (
    $1, $2, $3
) ON CONFLICT (account_id, target_account_id) 
DO UPDATE SET 
    uri = EXCLUDED.uri,
    updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: CreateStatus :one
INSERT INTO statuses (
    uri, url, local, content, account_id, in_reply_to_id, reblog_of_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: CreateFavourite :one
INSERT INTO favourites (
    uri, 
    account_id, 
    status_id
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetAccountFollowers :many
SELECT a.* FROM accounts a
JOIN follows f ON a.id = f.account_id
WHERE f.target_account_id = $1;

-- name: GetAccountFollowing :many
SELECT a.* FROM accounts a
JOIN follows f ON a.id = f.account_id
WHERE f.account_id = $1;

-- name: GetStatusComments :many
SELECT 
    sqlc.embed(s),
    sqlc.embed(a)
FROM statuses s
JOIN accounts a ON s.account_id = a.id
WHERE s.in_reply_to_id = $1;

-- name: UpdateStatus :one
UPDATE statuses
SET 
    content = COALESCE($1, content),
    updated_at = CURRENT_TIMESTAMP
WHERE id = $2
RETURNING *;