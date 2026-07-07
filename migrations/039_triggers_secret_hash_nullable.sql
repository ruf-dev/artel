-- +goose Up
ALTER TABLE triggers ALTER COLUMN secret_hash DROP NOT NULL;

-- +goose Down
-- Not reversible safely (would require backfilling secret_hash for provider-linked triggers);
-- intentionally left as a no-op.
