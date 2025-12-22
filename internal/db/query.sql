-- name: GetActor :one
SELECT * FROM accounts WHERE username = $1 AND domain IS NULL;

-- name: AuthData :one
SELECT a.id, u.password_hash FROM accounts a JOIN users u ON a.id = u.account_id WHERE a.username = $1;

-- name: CreateActor :one
INSERT INTO accounts (
    username, uri, display_name, domain, inbox_uri, outbox_uri, url
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: CreateUser :exec
INSERT INTO users (
    account_id, password_hash
) VALUES (
    $1, $2
);

-- name: GetAccount :one
SELECT * FROM accounts WHERE username = $1 AND domain = $2;

-- name: GetAllStatuses :many
SELECT s.id, s.created_at, s.updated_at, s.uri, s.url, s.local, s.content, s.account_id, s.in_reply_to_id, s.reblog_of_id, a.username 
FROM statuses s 
JOIN accounts a ON s.account_id = a.id 
ORDER BY s.created_at DESC;
