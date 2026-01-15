-- +goose Up
ALTER TABLE favourites
  DROP CONSTRAINT IF EXISTS favourites_status_id_fkey,
  ADD CONSTRAINT favourites_status_id_fkey
    FOREIGN KEY (status_id) REFERENCES statuses (id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE favourites
  DROP CONSTRAINT IF EXISTS favourites_status_id_fkey,
  ADD CONSTRAINT favourites_status_id_fkey
    FOREIGN KEY (status_id) REFERENCES statuses (id);
