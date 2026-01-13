-- +goose Up
CREATE UNIQUE INDEX IF NOT EXISTS statuses_unique_account_reblog_idx
  ON statuses (account_id, reblog_of_id)
  WHERE reblog_of_id IS NOT NULL;

-- +goose Down
DROP INDEX IF EXISTS statuses_unique_account_reblog_idx;
