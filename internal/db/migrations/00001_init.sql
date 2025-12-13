-- in implementations i looked up more fields are nullable
-- +goose Up
CREATE TABLE accounts (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    username TEXT NOT NULL,
    uri TEXT NOT NULL, -- AP identifier
    display_name TEXT,
    domain TEXT,
    inbox_uri TEXT NOT NULL,
    outbox_uri TEXT NOT NULL,
    followers_uri TEXT NOT NULL,
    following_uri TEXT NOT NULL,
    liked_uri TEXT NOT NULL, -- ??? according to the specification, not present in mastodon, gotosocial
    url TEXT NOT NULL
);

CREATE TABLE statuses (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    uri TEXT NOT NULL,
    url TEXT NOT NULL,
    local BOOLEAN DEFAULT FALSE,
    content TEXT NOT NULL,
    account_id INTEGER NOT NULL REFERENCES accounts (id),
    account_uri TEXT NOT NULL REFERENCES accounts (uri)
);

CREATE TABLE follows (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    uri TEXT NOT NULL,
    account_id INTEGER NOT NULL REFERENCES accounts (id),
    target_account_id INTEGER NOT NULL REFERENCES accounts (id)
);

CREATE TABLE users (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    account_id INTEGER NOT NULL REFERENCES accounts (id),
    password_hash TEXT NOT NULL
);

-- +goose Down
DROP TABLE accounts;
DROP TABLE statuses;
DROP TABLE follows;
DROP TABLE users;
