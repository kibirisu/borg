-- +goose Up
-- Insert sample accounts
INSERT INTO accounts (username, uri, display_name, domain, inbox_uri, outbox_uri, followers_uri, following_uri, url) VALUES
('alice', 'http://localhost:8080/users/alice', 'Alice Smith', NULL, 'http://localhost:8080/users/alice/inbox', 'http://localhost:8080/users/alice/outbox', 'http://localhost:8080/users/alice/followers', 'http://localhost:8080/users/alice/following', 'http://localhost:8080/profiles/alice'),
('bob', 'http://localhost:8080/users/bob', 'Bob Johnson', NULL, 'http://localhost:8080/users/bob/inbox', 'http://localhost:8080/users/bob/outbox', 'http://localhost:8080/users/bob/followers', 'http://localhost:8080/users/bob/following', 'http://localhost:8080/profiles/bob'),
('charlie', 'http://localhost:8080/users/charlie', 'Charlie Brown', NULL, 'http://localhost:8080/users/charlie/inbox', 'http://localhost:8080/users/charlie/outbox', 'http://localhost:8080/users/charlie/followers', 'http://localhost:8080/users/charlie/following', 'http://localhost:8080/profiles/charlie');

-- Insert sample users (linked to accounts)
INSERT INTO users (account_id, password_hash) VALUES
((SELECT id FROM accounts WHERE username = 'alice'), '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'), -- password: password123
((SELECT id FROM accounts WHERE username = 'bob'), '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'), -- password: password123
((SELECT id FROM accounts WHERE username = 'charlie'), '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy'); -- password: password123

-- Insert sample statuses (posts)
INSERT INTO statuses (uri, url, local, content, account_id) VALUES
('http://localhost:8080/statuses/1', 'http://localhost:8080/posts/1', TRUE, 'Hello world! This is my first post on this platform. #excited', (SELECT id FROM accounts WHERE username = 'alice')),
('http://localhost:8080/statuses/2', 'http://localhost:8080/posts/2', TRUE, 'Just finished reading an amazing book about distributed systems. Highly recommend!', (SELECT id FROM accounts WHERE username = 'bob')),
('http://localhost:8080/statuses/3', 'http://localhost:8080/posts/3', TRUE, 'Beautiful sunset today! 🌅', (SELECT id FROM accounts WHERE username = 'alice')),
('http://localhost:8080/statuses/4', 'http://localhost:8080/posts/4', TRUE, 'Working on a new project. Stay tuned for updates!', (SELECT id FROM accounts WHERE username = 'charlie')),
('http://localhost:8080/statuses/5', 'http://localhost:8080/posts/5', TRUE, 'Coffee and code - the perfect combination ☕💻', (SELECT id FROM accounts WHERE username = 'bob')),
('http://localhost:8080/statuses/6', 'http://localhost:8080/posts/6', TRUE, 'Learning new technologies is always exciting. What are you learning today?', (SELECT id FROM accounts WHERE username = 'alice'));

-- Insert sample follows
INSERT INTO follows (uri, account_id, target_account_id) VALUES
('http://localhost:8080/follows/1', (SELECT id FROM accounts WHERE username = 'alice'), (SELECT id FROM accounts WHERE username = 'bob')),
('http://localhost:8080/follows/2', (SELECT id FROM accounts WHERE username = 'bob'), (SELECT id FROM accounts WHERE username = 'charlie')),
('http://localhost:8080/follows/3', (SELECT id FROM accounts WHERE username = 'charlie'), (SELECT id FROM accounts WHERE username = 'alice'));

-- Insert sample favourites (likes)
INSERT INTO favourites (uri, account_id, status_id) VALUES
('http://localhost:8080/favourites/1', (SELECT id FROM accounts WHERE username = 'bob'), (SELECT id FROM statuses WHERE uri = 'http://localhost:8080/statuses/1')),
('http://localhost:8080/favourites/2', (SELECT id FROM accounts WHERE username = 'charlie'), (SELECT id FROM statuses WHERE uri = 'http://localhost:8080/statuses/1')),
('http://localhost:8080/favourites/3', (SELECT id FROM accounts WHERE username = 'alice'), (SELECT id FROM statuses WHERE uri = 'http://localhost:8080/statuses/2')),
('http://localhost:8080/favourites/4', (SELECT id FROM accounts WHERE username = 'charlie'), (SELECT id FROM statuses WHERE uri = 'http://localhost:8080/statuses/3'));

-- +goose Down
DELETE FROM favourites;
DELETE FROM follows;
DELETE FROM statuses;
DELETE FROM users;
DELETE FROM accounts;
