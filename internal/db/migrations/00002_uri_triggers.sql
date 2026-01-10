-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fn_generate_status_uri()
RETURNS TRIGGER AS $$
BEGIN
    SELECT '/users/' || username || '/statuses/' || NEW.id
    INTO NEW.uri
    FROM accounts
    WHERE id = NEW.account_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_create_status_uri
BEFORE INSERT ON statuses
FOR EACH ROW
EXECUTE FUNCTION fn_generate_status_uri();
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION fn_generate_favourite_uri()
RETURNS TRIGGER AS $$
BEGIN
    SELECT '/users/' || username || '/favourites/' || NEW.id
    INTO NEW.uri
    FROM accounts
    WHERE id = NEW.account_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_create_favourite_uri
BEFORE INSERT ON favourites
FOR EACH ROW
EXECUTE FUNCTION fn_generate_favourite_uri();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_create_status_uri ON statuses;
DROP FUNCTION IF EXISTS fn_generate_status_uri();

DROP TRIGGER IF EXISTS trg_create_favourite_uri ON favourites;
DROP FUNCTION IF EXISTS fn_generate_favourite_uri();
-- +goose StatementEnd
