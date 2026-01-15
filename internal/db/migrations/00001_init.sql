-- +goose Up
CREATE TABLE accounts (
    id VARCHAR(20) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
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
    id VARCHAR(20) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    uri TEXT UNIQUE NOT NULL,
    url TEXT NOT NULL,
    local BOOLEAN DEFAULT FALSE,
    content TEXT NOT NULL,
    account_id VARCHAR(20) NOT NULL REFERENCES accounts (id),
    account_uri TEXT NOT NULL REFERENCES accounts (uri),
    in_reply_to_id VARCHAR(20) REFERENCES statuses,
    in_reply_to_uri TEXT REFERENCES statuses (uri),
    reblog_of_id VARCHAR(20) REFERENCES statuses
);

CREATE TABLE follows (
    id VARCHAR(20) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    uri TEXT UNIQUE NOT NULL,
    account_id VARCHAR(20) NOT NULL REFERENCES accounts (id),
    target_account_id VARCHAR(20) NOT NULL REFERENCES accounts (id),
    UNIQUE (account_id, target_account_id),
    CHECK (account_id != target_account_id)
);

CREATE TABLE follow_requests (
    id VARCHAR(20) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    uri TEXT UNIQUE NOT NULL,
    account_id VARCHAR(20) NOT NULL REFERENCES accounts (id),
    target_account_id VARCHAR(20) NOT NULL REFERENCES accounts (id),
    UNIQUE (account_id, target_account_id),
    CHECK (account_id != target_account_id)
);

CREATE TABLE favourites (
    id VARCHAR(20) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    uri TEXT UNIQUE NOT NULL,
    account_id VARCHAR(20) NOT NULL REFERENCES accounts (id),
    status_id VARCHAR(20) NOT NULL REFERENCES statuses (id),
    UNIQUE (account_id, status_id)
);

CREATE TABLE users (
    id VARCHAR(20) PRIMARY KEY,
    account_id VARCHAR(20) NOT NULL UNIQUE REFERENCES accounts (id),
    password_hash TEXT NOT NULL
);

-- +goose Down
DROP TABLE follows;
DROP TABLE follow_requests;
DROP TABLE favourites;
DROP TABLE users;
DROP TABLE statuses;
DROP TABLE accounts;
