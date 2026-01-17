-- name: GetActor :one
SELECT * FROM accounts WHERE username = $1 AND domain IS NULL;

-- name: GetActorByURI :one
SELECT * FROM accounts WHERE uri LIKE '%' || $1::text;

-- name: AuthData :one
SELECT a.id, u.password_hash FROM accounts a JOIN users u ON a.id = u.account_id WHERE a.username = $1;

-- name: CreateActor :one
INSERT INTO accounts (
    id, username, uri, display_name, domain, inbox_uri, outbox_uri, url, followers_uri, following_uri
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: CreateUser :exec
INSERT INTO users (
    id, account_id, password_hash
) VALUES (
    $1, $2, $3
);

-- name: GetAccountByID :one
SELECT 
    sqlc.embed(a),
    (a.username || COALESCE('@' || a.domain, ''))::text AS acct,
    (SELECT COUNT(*) FROM statuses s WHERE s.account_id = a.id) AS statuses_count,
    (SELECT COUNT(*) FROM follows f WHERE f.target_account_id = a.id) AS followers_count,
    (SELECT COUNT(*) FROM follows f WHERE f.account_id = a.id) AS following_count
FROM accounts a WHERE a.id = $1;

-- name: GetFollowersByAccountID :many
SELECT 
    sqlc.embed(a),
    (a.username || COALESCE('@' || a.domain, ''))::text AS acct,
    (SELECT COUNT(*) FROM statuses s WHERE s.account_id = a.id) AS statuses_count,
    (SELECT COUNT(*) FROM follows f WHERE f.target_account_id = a.id) AS followers_count,
    (SELECT COUNT(*) FROM follows f WHERE f.account_id = a.id) AS following_count
FROM accounts a JOIN follows f ON a.id = f.account_id WHERE f.target_account_id = $1;

-- name: GetFollowingByAccountID :many
SELECT 
    sqlc.embed(a),
    (a.username || COALESCE('@' || a.domain, ''))::text AS acct,
    (SELECT COUNT(*) FROM statuses s WHERE s.account_id = a.id) AS statuses_count,
    (SELECT COUNT(*) FROM follows f WHERE f.target_account_id = a.id) AS followers_count,
    (SELECT COUNT(*) FROM follows f WHERE f.account_id = a.id) AS following_count
FROM accounts a JOIN follows f ON a.id = f.target_account_id WHERE f.account_id = $1;

-- name: GetAccountRemoteFollowersInboxes :many
SELECT inbox_uri FROM accounts a JOIN follows f ON a.id = f.account_id WHERE f.target_account_id = $1 AND a.domain IS NOT NULL;

-- name: GetAccountInbox :one
SELECT inbox_uri FROM accounts WHERE id = $1;

-- name: GetStatusById :one
SELECT * FROM statuses WHERE id = $1;

-- name: GetStatusByIDNew :one
SELECT 
    sqlc.embed(s),
    (SELECT COUNT(*) FROM statuses r WHERE r.in_reply_to_id = s.id) AS replies_count,
    (SELECT COUNT(*) FROM favourites f WHERE f.status_id = s.id) AS favourites_count,
    (SELECT COUNT(*) FROM statuses r WHERE r.reblog_of_id = s.id) AS reblogs_count,
    EXISTS(SELECT 1 FROM favourites f WHERE f.status_id = s.id AND f.account_id = $2) AS favourited,
    EXISTS(SELECT 1 FROM statuses r WHERE r.reblog_of_id = s.id AND r.account_id = $2) AS reblogged
FROM statuses s WHERE s.id = $1;

-- name: GetStatusByURI :one
SELECT * FROM statuses WHERE uri = $1;

-- name: GetStatusesByAccountID :many
SELECT 
    sqlc.embed(s),
    (SELECT COUNT(*) FROM statuses r WHERE r.in_reply_to_id = s.id) AS replies_count,
    (SELECT COUNT(*) FROM favourites f WHERE f.status_id = s.id) AS favourites_count,
    (SELECT COUNT(*) FROM statuses r WHERE r.reblog_of_id = s.id) AS reblogs_count,
    EXISTS(SELECT 1 FROM favourites f WHERE f.status_id = s.id AND f.account_id = sqlc.arg(logged_in_id)) AS favourited,
    EXISTS(SELECT 1 FROM statuses r WHERE r.reblog_of_id = s.id AND r.account_id = sqlc.arg(logged_in_id)) AS reblogged
FROM statuses s WHERE s.account_id = sqlc.arg(account_id);

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

-- name: CreateFollow :one
INSERT INTO follows (
    id, uri, account_id, target_account_id
) VALUES (
    $1, $2, $3, $4
) ON CONFLICT (account_id, target_account_id) 
DO UPDATE SET 
    uri = EXCLUDED.uri,
    updated_at = CURRENT_TIMESTAMP
RETURNING *;

-- name: GetFollowerCollection :one
SELECT 
    (SELECT followers_uri FROM accounts a WHERE a.username = $1),
    (SELECT COUNT(*) FROM follows f JOIN accounts a ON f.target_account_id = a.id WHERE a.username = $1);

-- name: GetFollowingCollection :one
SELECT 
    (SELECT following_uri FROM accounts a WHERE a.username = $1),
    (SELECT COUNT(*) FROM follows f JOIN accounts a ON f.account_id = a.id WHERE a.username = $1);

-- name: CreateFollowRequest :one
INSERT INTO follow_requests (
    id, uri, account_id, target_account_id, target_account_uri
) VALUES (
    $1, $2, $3, $4, (SELECT uri FROM accounts WHERE id = $4)
) RETURNING *;

-- name: CreateStatus :one
INSERT INTO statuses (
    id, url, local, content, account_id, account_uri, in_reply_to_id, reblog_of_id, uri
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: CreateStatusNew :one
WITH parent AS (
    SELECT uri, account_id FROM statuses WHERE id = $7
) INSERT INTO statuses (
    id, uri, url, local, content, account_id, account_uri, 
    in_reply_to_id, in_reply_to_uri, in_reply_to_account_id
) VALUES (
    $1, $2, $3, true, $4, $5, $6, $7,
    (SELECT uri FROM parent),
    (SELECT account_id FROM parent)
) RETURNING *;

-- name: DeleteStatusByID :exec
DELETE FROM statuses WHERE id = $1;

-- name: CreateFavourite :one
INSERT INTO favourites (
    id,
    account_id, 
    status_id,
    uri
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetFavouriteByURI :one
SELECT * FROM favourites WHERE uri LIKE '%' || $1::text;

-- name: GetFollowByURI :one
SELECT * FROM follows WHERE uri LIKE '%' || $1::text;

-- name: DeleteFavouriteByID :exec
DELETE FROM favourites WHERE id = $1;

-- name: GetAccountFollowers :many
SELECT a.* FROM accounts a
JOIN follows f ON a.id = f.account_id
WHERE f.target_account_id = $1;

-- name: GetAccountFollowing :many
SELECT a.* FROM accounts a
JOIN follows f ON a.id = f.target_account_id
WHERE f.account_id = $1;

-- name: GetLikedPostsByAccountId :many
SELECT 
    sqlc.embed(s),
    sqlc.embed(a),
    (SELECT COUNT(*) FROM favourites f WHERE f.status_id = s.id) AS like_count,
    (SELECT COUNT(*) FROM statuses r WHERE r.in_reply_to_id = s.id) AS comment_count,
    (SELECT COUNT(*) FROM statuses b WHERE b.reblog_of_id = s.id) AS share_count
FROM favourites f
JOIN statuses s ON f.status_id = s.id
JOIN accounts a ON s.account_id = a.id
WHERE f.account_id = $1;

-- name: GetSharedPostsByAccountId :many
SELECT 
    sqlc.embed(s),
    sqlc.embed(a),
    (SELECT COUNT(*) FROM favourites f WHERE f.status_id = s.id) AS like_count,
    (SELECT COUNT(*) FROM statuses r WHERE r.in_reply_to_id = s.id) AS comment_count,
    (SELECT COUNT(*) FROM statuses b WHERE b.reblog_of_id = s.id) AS share_count
FROM statuses s
JOIN accounts a ON s.account_id = a.id
WHERE s.account_id = $1 AND s.reblog_of_id IS NOT NULL;

-- name: GetTimelinePostsByAccountId :many
SELECT 
    sqlc.embed(s),
    sqlc.embed(a),
    (SELECT COUNT(*) FROM favourites f WHERE f.status_id = s.id) AS like_count,
    (SELECT COUNT(*) FROM statuses r WHERE r.in_reply_to_id = s.id) AS comment_count,
    (SELECT COUNT(*) FROM statuses b WHERE b.reblog_of_id = s.id) AS share_count
FROM statuses s
JOIN accounts a ON s.account_id = a.id
JOIN follows f ON a.id = f.target_account_id
WHERE f.account_id = $1 AND s.in_reply_to_id IS NULL
ORDER BY s.created_at DESC;
