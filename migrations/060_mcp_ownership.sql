-- +goose Up
ALTER TABLE mcps
    ADD COLUMN owner_user_id UUID NULL REFERENCES users(id) ON DELETE SET NULL,
    ADD COLUMN is_community  BOOLEAN NOT NULL DEFAULT false;

CREATE INDEX ON mcps(is_community);

-- +goose Down
DROP INDEX IF EXISTS mcps_is_community_idx;
ALTER TABLE mcps DROP COLUMN owner_user_id, DROP COLUMN is_community;
