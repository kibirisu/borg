-- +goose Up
ALTER TABLE statuses
  DROP CONSTRAINT IF EXISTS statuses_reblog_of_id_fkey,
  ADD CONSTRAINT statuses_reblog_of_id_fkey
    FOREIGN KEY (reblog_of_id) REFERENCES statuses (id) ON DELETE CASCADE;

ALTER TABLE statuses
  DROP CONSTRAINT IF EXISTS statuses_in_reply_to_id_fkey,
  ADD CONSTRAINT statuses_in_reply_to_id_fkey
    FOREIGN KEY (in_reply_to_id) REFERENCES statuses (id) ON DELETE CASCADE;

ALTER TABLE favourites
  DROP CONSTRAINT IF EXISTS favourites_status_id_fkey,
  ADD CONSTRAINT favourites_status_id_fkey
    FOREIGN KEY (status_id) REFERENCES statuses (id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE statuses
  DROP CONSTRAINT IF EXISTS statuses_reblog_of_id_fkey,
  ADD CONSTRAINT statuses_reblog_of_id_fkey
    FOREIGN KEY (reblog_of_id) REFERENCES statuses (id);

ALTER TABLE statuses
  DROP CONSTRAINT IF EXISTS statuses_in_reply_to_id_fkey,
  ADD CONSTRAINT statuses_in_reply_to_id_fkey
    FOREIGN KEY (in_reply_to_id) REFERENCES statuses (id);

ALTER TABLE favourites
  DROP CONSTRAINT IF EXISTS favourites_status_id_fkey,
  ADD CONSTRAINT favourites_status_id_fkey
    FOREIGN KEY (status_id) REFERENCES statuses (id);
