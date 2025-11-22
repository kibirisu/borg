-- +goose Up
ALTER TABLE users
ADD COLUMN origin VARCHAR(255);

UPDATE users
SET origin = '127.0.0.1'
WHERE origin IS NULL;

-- Increment followers
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION increment_follow_counters()
RETURNS trigger AS
'
BEGIN
    UPDATE users
        SET followers_count = followers_count + 1
        WHERE id = NEW.following_id;

    UPDATE users
        SET following_count = following_count + 1
        WHERE id = NEW.follower_id;
    RETURN NEW;
END;
' LANGUAGE plpgsql;
-- +goose StatementEnd
CREATE TRIGGER trg_followers_insert
AFTER INSERT ON followers
FOR EACH ROW
EXECUTE FUNCTION increment_follow_counters();

-- Increment followers
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION decrement_follow_counters()
RETURNS trigger AS
'
BEGIN
    UPDATE users
        SET followers_count = followers_count - 1
        WHERE id = OLD.following_id;

    UPDATE users
        SET following_count = following_count - 1
        WHERE id = OLD.follower_id;

    RETURN OLD;
END;
' LANGUAGE plpgsql;
-- +goose StatementEnd
CREATE TRIGGER trg_followers_delete
AFTER DELETE ON followers
FOR EACH ROW
EXECUTE FUNCTION decrement_follow_counters();

-- Increment likes
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION increment_like_count()
RETURNS trigger AS
'
BEGIN
    UPDATE posts
        SET like_count = like_count + 1
        WHERE id = NEW.post_id;
    RETURN NEW;
END;
' LANGUAGE plpgsql;
-- +goose StatementEnd
CREATE TRIGGER trg_likes_insert
AFTER INSERT ON likes
FOR EACH ROW
EXECUTE FUNCTION increment_like_count();

-- Decrement likes
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION decrement_like_count()
RETURNS trigger AS
'
BEGIN
    UPDATE posts
        SET like_count = like_count - 1
        WHERE id = OLD.post_id;
    RETURN OLD;
END;
' LANGUAGE plpgsql;
-- +goose StatementEnd
CREATE TRIGGER trg_likes_delete
AFTER DELETE ON likes
FOR EACH ROW
EXECUTE FUNCTION decrement_like_count();


-- Increment shares
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION increment_share_count()
RETURNS trigger AS
'
BEGIN
    UPDATE posts
        SET share_count = share_count + 1
        WHERE id = NEW.post_id;
    RETURN NEW;
END;
' LANGUAGE plpgsql;
-- +goose StatementEnd
CREATE TRIGGER trg_shares_insert
AFTER INSERT ON shares
FOR EACH ROW
EXECUTE FUNCTION increment_share_count();


-- Decrement shares
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION decrement_share_count()
RETURNS trigger AS
'
BEGIN
    UPDATE posts
        SET share_count = share_count - 1
        WHERE id = OLD.post_id;
    RETURN OLD;
END;
' LANGUAGE plpgsql;
-- +goose StatementEnd
CREATE TRIGGER trg_shares_delete
AFTER DELETE ON shares
FOR EACH ROW
EXECUTE FUNCTION decrement_share_count();

-- Increment comments
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION increment_comment_count()
RETURNS trigger AS
'
BEGIN
    UPDATE posts
        SET comment_count = comment_count + 1
        WHERE id = NEW.post_id;
    RETURN NEW;
END;
' LANGUAGE plpgsql;
-- +goose StatementEnd
CREATE TRIGGER trg_comments_insert
AFTER INSERT ON comments
FOR EACH ROW
EXECUTE FUNCTION increment_comment_count();

-- Decrement comments
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION decrement_comment_count()
RETURNS trigger AS
'
BEGIN
    UPDATE posts
        SET comment_count = comment_count - 1
        WHERE id = OLD.post_id;
    RETURN OLD;
END;
' LANGUAGE plpgsql;
-- +goose StatementEnd
CREATE TRIGGER trg_comments_delete
AFTER DELETE ON comments
FOR EACH ROW
EXECUTE FUNCTION decrement_comment_count();

BEGIN;
UPDATE users
SET followers_count = 0,
    following_count = 0;

UPDATE posts
SET like_count = 0,
    share_count = 0,
    comment_count = 0;

-- +goose StatementBegin
DO
'
BEGIN
    -- Likes
    CREATE TEMP TABLE tmp_likes AS SELECT * FROM likes;
    TRUNCATE likes;
    INSERT INTO likes (post_id, user_id, created_at)
    SELECT post_id, user_id, created_at FROM tmp_likes;

    -- Shares
    CREATE TEMP TABLE tmp_shares AS SELECT * FROM shares;
    TRUNCATE shares;
    INSERT INTO shares (post_id, user_id, created_at)
    SELECT post_id, user_id, created_at FROM tmp_shares;

    -- Comments
    CREATE TEMP TABLE tmp_comments AS SELECT * FROM comments;
    TRUNCATE comments;
    INSERT INTO comments (post_id, user_id, content, parent_id, created_at, updated_at)
    SELECT post_id, user_id, content, parent_id, created_at, updated_at FROM tmp_comments;
END
';
-- +goose StatementEnd
COMMIT;

-- +goose Down
ALTER TABLE users
DROP COLUMN origin;

DROP TRIGGER IF EXISTS trg_followers_delete ON followers;
DROP TRIGGER IF EXISTS trg_followers_insert ON followers;
DROP FUNCTION IF EXISTS decrement_follow_counters();
DROP FUNCTION IF EXISTS increment_follow_counters();

DROP TRIGGER IF EXISTS trg_comments_delete ON comments;
DROP TRIGGER IF EXISTS trg_comments_insert ON comments;
DROP FUNCTION IF EXISTS decrement_comment_count();
DROP FUNCTION IF EXISTS increment_comment_count();

DROP TRIGGER IF EXISTS trg_shares_delete ON shares;
DROP TRIGGER IF EXISTS trg_shares_insert ON shares;
DROP FUNCTION IF EXISTS decrement_share_count();
DROP FUNCTION IF EXISTS increment_share_count();

DROP TRIGGER IF EXISTS trg_likes_delete ON likes;
DROP TRIGGER IF EXISTS trg_likes_insert ON likes;
DROP FUNCTION IF EXISTS decrement_like_count();
DROP FUNCTION IF EXISTS increment_like_count();
