-- +goose Up
CREATE TABLE accounts (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    username TEXT NOT NULL,
    uri TEXT UNIQUE NOT NULL,
    display_name TEXT,
    domain TEXT,
    inbox_uri TEXT NOT NULL,
    outbox_uri TEXT NOT NULL,
    followers_uri TEXT NOT NULL,
    following_uri TEXT NOT NULL,
    url TEXT NOT NULL,
    UNIQUE (username, domain)
);

CREATE TABLE statuses (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    uri TEXT UNIQUE NOT NULL,
    url TEXT NOT NULL,
    local BOOLEAN DEFAULT FALSE,
    content TEXT NOT NULL,
    account_id INTEGER NOT NULL REFERENCES accounts (id),
    in_reply_to_id INTEGER REFERENCES statuses,
    reblog_of_id INTEGER REFERENCES statuses
);

CREATE TABLE follows (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    uri TEXT UNIQUE NOT NULL,
    account_id INTEGER NOT NULL REFERENCES accounts (id),
    target_account_id INTEGER NOT NULL REFERENCES accounts (id),
    UNIQUE (account_id, target_account_id),
    CHECK (account_id != target_account_id)
);

CREATE TABLE favourites (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    uri TEXT UNIQUE NOT NULL,
    account_id INTEGER NOT NULL REFERENCES accounts (id),
    status_id INTEGER NOT NULL REFERENCES statuses (id),
    UNIQUE (account_id, status_id)
);

CREATE TABLE users (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    account_id INTEGER NOT NULL UNIQUE REFERENCES accounts (id),
    password_hash TEXT NOT NULL
);

-- +goose Down
DROP TABLE accounts;
DROP TABLE statuses;
DROP TABLE follows;
DROP TABLE favourites;
DROP TABLE users;
