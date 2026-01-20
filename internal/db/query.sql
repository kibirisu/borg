-- name: GetActorByURI :one
SELECT * FROM accounts WHERE uri LIKE '%' || $1::text;

-- name: AuthData :one
SELECT a.id, a.uri, u.password_hash FROM accounts a JOIN users u ON a.id = u.account_id WHERE a.username = $1;

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

-- name: AddStatusByActorURI :exec
INSERT INTO statuses (
    id, uri, url, content, account_id, account_uri
) SELECT @id, @uri, @url, @content, a.id, @account_uri FROM accounts a WHERE a.uri = @account_uri;

-- name: AddFollowByActorURI :one
WITH follower AS (
    SELECT a.id, a.inbox_uri FROM accounts a WHERE a.uri = @account_uri
), follow AS (
    INSERT INTO follows (
        id, uri, account_id, target_account_id
    ) SELECT @id, @uri, follower.id, @target_account_id FROM follower
) SELECT inbox_uri FROM follower;

-- name: AddAccount :one
INSERT INTO accounts (
    id, username, uri, domain, inbox_uri, outbox_uri, followers_uri, following_uri, url
) VALUES (
    @id, @username, @uri, @domain, @inbox_uri, @outbox_uri, @followers_uri, @following_uri, @url
) RETURNING *;

-- name: AddStatus :exec
INSERT INTO statuses (
    id, uri, url, content, account_id, account_uri, in_reply_to_id, in_reply_to_uri, in_reply_to_account_id
) VALUES (
    @id, @uri, @url, @content, @account_id, @account_uri, @in_reply_to_id, @in_reply_to_uri, @in_reply_to_account_id
);

-- name: GetLocalActorByID :one
SELECT * FROM accounts WHERE id = $1 AND domain IS NULL;

-- name: GetAccountWebfinger :one
SELECT uri AS href, 'self' AS rel, 'application/activity+json' AS type FROM accounts WHERE username = $1 AND domain IS NULL;

-- name: GetAccountByID :one
SELECT 
    sqlc.embed(a),
    CONCAT(a.username, '@' || a.domain)::TEXT AS acct,
    (SELECT COUNT(*) FROM follows f WHERE f.target_account_id = a.id) AS followers_count,
    (SELECT COUNT(*) FROM follows f WHERE f.account_id = a.id) AS following_count
FROM accounts a WHERE a.id = $1;

-- name: GetFollowersByAccountID :many
SELECT 
    sqlc.embed(a),
    CONCAT(a.username, '@' || a.domain)::TEXT AS acct,
    (SELECT COUNT(*) FROM follows f WHERE f.target_account_id = a.id) AS followers_count,
    (SELECT COUNT(*) FROM follows f WHERE f.account_id = a.id) AS following_count
FROM accounts a JOIN follows f ON a.id = f.account_id WHERE f.target_account_id = $1;

-- name: GetFollowingByAccountID :many
SELECT 
    sqlc.embed(a),
    CONCAT(a.username, '@' || a.domain)::TEXT AS acct,
    (SELECT COUNT(*) FROM follows f WHERE f.target_account_id = a.id) AS followers_count,
    (SELECT COUNT(*) FROM follows f WHERE f.account_id = a.id) AS following_count
FROM accounts a JOIN follows f ON a.id = f.target_account_id WHERE f.account_id = $1;

-- name: GetAccountByUsernameAndDomain :one
SELECT
    sqlc.embed(a),
    CONCAT(a.username, '@' || a.domain)::TEXT AS acct,
    (SELECT COUNT(*) FROM follows f WHERE f.target_account_id = a.id) AS followers_count,
    (SELECT COUNT(*) FROM follows f WHERE f.account_id = a.id) AS following_count
FROM accounts a WHERE a.username = $1 AND (a.domain = $2 OR (a.domain IS NULL AND $2 IS NULL));

-- name: GetAccountRemoteFollowersInboxes :many
SELECT inbox_uri FROM accounts a JOIN follows f ON a.id = f.account_id WHERE f.target_account_id = $1 AND a.domain IS NOT NULL;

-- name: GetAccountInbox :one
SELECT inbox_uri FROM accounts WHERE id = $1;

-- name: GetLocalStatusByID :one
SELECT * FROM statuses WHERE id = $1 AND local;

-- name: GetStatusByID :one
SELECT 
    sqlc.embed(s),
    sqlc.embed(a),
    reblogged.content AS reblogged_status_content,
    reblogged.in_reply_to_id AS reblogged_reply_to_id,
    reblogged.in_reply_to_account_id AS reblogged_reply_to_account_id,
    reblogged_author.username AS reblogged_username,
    reblogged_author.display_name AS reblogged_display_name,
    CONCAT(reblogged_author.username, '@' || a.domain)::TEXT AS reblogged_acct,
    CONCAT(a.username, '@' || a.domain)::TEXT AS acct,
    (SELECT COUNT(*) FROM follows f WHERE f.target_account_id = a.id) AS followers_count,
    (SELECT COUNT(*) FROM follows f WHERE f.target_account_id = reblogged_author.id) AS reblogged_followers_count,
    (SELECT COUNT(*) FROM follows f WHERE f.account_id = a.id) AS following_count,
    (SELECT COUNT(*) FROM follows f WHERE f.account_id = reblogged_author.id) AS reblogged_following_count,
    (SELECT COUNT(*) FROM statuses r WHERE r.in_reply_to_id = COALESCE(s.reblog_of_id, s.id)) AS replies_count,
    (SELECT COUNT(*) FROM favourites f WHERE f.status_id = COALESCE(s.reblog_of_id, s.id)) AS favourites_count,
    (SELECT COUNT(*) FROM statuses r WHERE r.reblog_of_id = COALESCE(s.reblog_of_id, s.id)) AS reblogs_count,
    EXISTS(SELECT 1 FROM favourites f WHERE f.status_id = COALESCE(s.reblog_of_id, s.id) AND f.account_id = $2) AS favourited,
    EXISTS(SELECT 1 FROM statuses r WHERE r.reblog_of_id = COALESCE(s.reblog_of_id, s.id) AND r.account_id = $2) AS reblogged
FROM statuses s
JOIN accounts a ON s.account_id = a.id
LEFT JOIN statuses reblogged ON s.reblog_of_id = reblogged.id
LEFT JOIN accounts reblogged_author ON reblogged.account_id = reblogged_author.id
WHERE s.id = $1;

-- name: GetStatusReplies :many
SELECT 
    sqlc.embed(s),
    sqlc.embed(a),
    reblogged.content AS reblogged_status_content,
    reblogged.in_reply_to_id AS reblogged_reply_to_id,
    reblogged.in_reply_to_account_id AS reblogged_reply_to_account_id,
    reblogged_author.username AS reblogged_username,
    reblogged_author.display_name AS reblogged_display_name,
    CONCAT(reblogged_author.username, '@', reblogged_author.domain)::TEXT AS reblogged_acct,
    CONCAT(a.username, '@', a.domain)::TEXT AS acct,
    (SELECT COUNT(*) FROM follows f WHERE f.target_account_id = a.id) AS followers_count,
    (SELECT COUNT(*) FROM follows f WHERE f.target_account_id = reblogged_author.id) AS reblogged_followers_count,
    (SELECT COUNT(*) FROM follows f WHERE f.account_id = a.id) AS following_count,
    (SELECT COUNT(*) FROM follows f WHERE f.account_id = reblogged_author.id) AS reblogged_following_count,
    (SELECT COUNT(*) FROM statuses r WHERE r.in_reply_to_id = COALESCE(s.reblog_of_id, s.id)) AS replies_count,
    (SELECT COUNT(*) FROM favourites f WHERE f.status_id = COALESCE(s.reblog_of_id, s.id)) AS favourites_count,
    (SELECT COUNT(*) FROM statuses r WHERE r.reblog_of_id = COALESCE(s.reblog_of_id, s.id)) AS reblogs_count,
    EXISTS(SELECT 1 FROM favourites f WHERE f.status_id = COALESCE(s.reblog_of_id, s.id) AND f.account_id = $2) AS favourited,
    EXISTS(SELECT 1 FROM statuses r WHERE r.reblog_of_id = COALESCE(s.reblog_of_id, s.id) AND r.account_id = $2) AS reblogged
FROM statuses s
JOIN accounts a ON s.account_id = a.id
LEFT JOIN statuses reblogged ON s.reblog_of_id = reblogged.id
LEFT JOIN accounts reblogged_author ON reblogged.account_id = reblogged_author.id
WHERE s.in_reply_to_id = $1
ORDER BY s.created_at ASC;

-- name: GetStatusByURI :one
SELECT * FROM statuses WHERE uri = $1;

-- name: GetStatusesByAccountID :many
SELECT 
    sqlc.embed(s),
    sqlc.embed(a),
    reblogged.content AS reblogged_status_content,
    reblogged.in_reply_to_id AS reblogged_reply_to_id,
    reblogged.in_reply_to_account_id AS reblogged_reply_to_account_id,
    reblogged_author.username AS reblogged_username,
    reblogged_author.display_name AS reblogged_display_name,
    CONCAT(reblogged_author.username, '@' || a.domain)::TEXT AS reblogged_acct,
    CONCAT(a.username, '@' || a.domain)::TEXT AS acct,
    (SELECT COUNT(*) FROM follows f WHERE f.target_account_id = a.id) AS followers_count,
    (SELECT COUNT(*) FROM follows f WHERE f.target_account_id = reblogged_author.id) AS reblogged_followers_count,
    (SELECT COUNT(*) FROM follows f WHERE f.account_id = a.id) AS following_count,
    (SELECT COUNT(*) FROM follows f WHERE f.account_id = reblogged_author.id) AS reblogged_following_count,
    (SELECT COUNT(*) FROM statuses r WHERE r.in_reply_to_id = COALESCE(s.reblog_of_id, s.id)) AS replies_count,
    (SELECT COUNT(*) FROM favourites f WHERE f.status_id = COALESCE(s.reblog_of_id, s.id)) AS favourites_count,
    (SELECT COUNT(*) FROM statuses r WHERE r.reblog_of_id = COALESCE(s.reblog_of_id, s.id)) AS reblogs_count,
    EXISTS(SELECT 1 FROM favourites f WHERE f.status_id = COALESCE(s.reblog_of_id, s.id) AND f.account_id = @logged_in_id) AS favourited,
    EXISTS(SELECT 1 FROM statuses r WHERE r.reblog_of_id = COALESCE(s.reblog_of_id, s.id) AND r.account_id = @logged_in_id) AS reblogged
FROM statuses s
JOIN accounts a ON s.account_id = a.id
LEFT JOIN statuses reblogged ON s.reblog_of_id = reblogged.id
LEFT JOIN accounts reblogged_author ON reblogged.account_id = reblogged_author.id
WHERE s.account_id = @account_id;

-- name: DeleteStatusByIDNew :one
DELETE FROM statuses WHERE id = $1 RETURNING *;

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

-- name: CreateFollowNew :exec
INSERT INTO follows (
  id, uri, account_id, target_account_id
) VALUES (
  @id, @uri, @account_id, @target_account_id
);

-- name: DeleteFollow :exec
DELETE FROM follows WHERE id = $1;

-- name: GetFollowerCollection :one
SELECT 
    (SELECT followers_uri FROM accounts a WHERE a.username = $1),
    (SELECT COUNT(*) FROM follows f JOIN accounts a ON f.target_account_id = a.id WHERE a.username = $1);

-- name: GetFollowingCollection :one
SELECT 
    (SELECT following_uri FROM accounts a WHERE a.username = $1),
    (SELECT COUNT(*) FROM follows f JOIN accounts a ON f.account_id = a.id WHERE a.username = $1);

-- name: CreateFollowRequest :one
WITH account AS (
  SELECT a.uri, (a.domain IS NULL)::BOOLEAN AS local FROM accounts a WHERE a.id = @target_account_id
), request AS (
  INSERT INTO follow_requests (
    id, uri, account_id, target_account_id, target_account_uri
  ) SELECT @id, @uri, @account_id, @target_account_id, uri FROM account RETURNING *
) SELECT r.*, account.local FROM request r, account;

-- name: DeleteFollowRequestByAccountID :one
WITH account AS (
  SELECT a.id, a.uri, (a.domain IS NULL)::BOOLEAN AS local FROM accounts a WHERE a.id = @target_account_id
), request AS (
    DELETE FROM follow_requests WHERE account_id = @account_id AND target_account_id = (SELECT id FROM account) RETURNING *
) SELECT r.*, account.local FROM request r, account;

-- name: CreateStatus :one
INSERT INTO statuses (
    id, url, local, content, account_id, account_uri, in_reply_to_id, reblog_of_id, uri
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: CreateStatusNew :one
WITH parent AS (
    SELECT uri, account_id FROM statuses WHERE id = @in_reply_to_id
) INSERT INTO statuses (
    id, uri, url, local, content, account_id, account_uri, 
    in_reply_to_id, in_reply_to_uri, in_reply_to_account_id
) VALUES (
    @id, @uri, @url, true, @content, @account_id, @account_uri, @in_reply_to_id,
    (SELECT uri FROM parent),
    (SELECT account_id FROM parent)
) RETURNING *;

-- name: DeleteStatusByID :exec
DELETE FROM statuses WHERE id = $1;

-- name: CreateReblog :one
WITH parent AS (
    SELECT s.uri, s.account_id FROM statuses s WHERE s.id = @reblog_of_id
) INSERT INTO statuses (
    id, uri, url, local, account_id, account_uri, 
    reblog_of_id, reblog_of_uri, reblog_of_account_id
) VALUES (
    @id, @uri, @url, true, @account_id, @account_uri, @reblog_of_id,
    (SELECT uri FROM parent),
    (SELECT account_id FROM parent)
) RETURNING *;

-- name: GetLocalLikeByID :one
SELECT 
    sqlc.embed(f),
    sqlc.embed(a),
    sqlc.embed(s)
FROM favourites f JOIN accounts a ON f.account_id = a.id
JOIN statuses s ON f.status_id = s.id
WHERE f.id = $1 AND a.domain IS NULL;

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

-- name: CreateFavouriteNew :one
WITH favourited AS (
    SELECT account_id, account_uri, uri FROM statuses WHERE id = $5
) INSERT INTO favourites (
    id, uri, account_id, account_uri, target_account_id, status_id, status_uri
) VALUES (
    $1, $2, $3, $4,
    (SELECT account_id FROM favourited),
    $5,
    (SELECT uri FROM favourited)
) RETURNING *;

-- name: GetFavouriteByURI :one
SELECT * FROM favourites WHERE uri LIKE '%' || $1::text;

-- name: GetLocalFollowByID :one
SELECT
    f.uri AS follow_uri,
    a1.uri AS following_uri,
    a2.uri AS followed_uri
FROM follows f JOIN accounts a1 ON f.account_id = a1.id
JOIN accounts a2 ON f.target_account_id = a2.id
WHERE f.id = $1 AND a1.domain IS NULL;

-- name: DeleteFavouriteByID :exec
DELETE FROM favourites WHERE id = $1;

-- name: DeleteFavouriteByIDNew :one
DELETE FROM favourites WHERE id = $1 RETURNING *;

-- name: DeleteFavouriteByStatusID :one
DELETE FROM favourites WHERE account_id = $1 AND status_id = $2 RETURNING *;

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

-- name: GetTimelinePostsByAccountId :many
SELECT 
    sqlc.embed(s),
    sqlc.embed(a),
    reblogged.content AS reblogged_status_content,
    reblogged.in_reply_to_id AS reblogged_reply_to_id,
    reblogged.in_reply_to_account_id AS reblogged_reply_to_account_id,
    reblogged_author.username AS reblogged_username,
    reblogged_author.display_name AS reblogged_display_name,
    CONCAT(reblogged_author.username, '@', reblogged_author.domain)::TEXT AS reblogged_acct,
    CONCAT(a.username, '@', a.domain)::TEXT AS acct,
    (SELECT COUNT(*) FROM follows f WHERE f.target_account_id = a.id) AS followers_count,
    (SELECT COUNT(*) FROM follows f WHERE f.target_account_id = reblogged_author.id) AS reblogged_followers_count,
    (SELECT COUNT(*) FROM follows f WHERE f.account_id = a.id) AS following_count,
    (SELECT COUNT(*) FROM follows f WHERE f.account_id = reblogged_author.id) AS reblogged_following_count,
    (SELECT COUNT(*) FROM statuses r WHERE r.in_reply_to_id = COALESCE(s.reblog_of_id, s.id)) AS replies_count,
    (SELECT COUNT(*) FROM favourites f WHERE f.status_id = COALESCE(s.reblog_of_id, s.id)) AS favourites_count,
    (SELECT COUNT(*) FROM statuses r WHERE r.reblog_of_id = COALESCE(s.reblog_of_id, s.id)) AS reblogs_count,
    EXISTS(SELECT 1 FROM favourites f WHERE f.status_id = COALESCE(s.reblog_of_id, s.id) AND f.account_id = $1) AS favourited,
    EXISTS(SELECT 1 FROM statuses r WHERE r.reblog_of_id = COALESCE(s.reblog_of_id, s.id) AND r.account_id = $1) AS reblogged
FROM statuses s
JOIN accounts a ON s.account_id = a.id
JOIN follows f_logic ON a.id = f_logic.target_account_id
LEFT JOIN statuses reblogged ON s.reblog_of_id = reblogged.id
LEFT JOIN accounts reblogged_author ON reblogged.account_id = reblogged_author.id
WHERE f_logic.account_id = $1 
  AND s.in_reply_to_id IS NULL
ORDER BY s.created_at DESC;
