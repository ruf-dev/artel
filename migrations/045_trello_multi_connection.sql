-- +goose Up
-- Allows a user to hold multiple Trello connections at once, while every other provider keeps
-- today's "one connection per provider" constraint. A table-wide UNIQUE(user_id, provider) can't
-- express that, so it's replaced by a partial unique index that simply excludes trello.
ALTER TABLE external_connections DROP CONSTRAINT external_connections_user_id_provider_key;

CREATE UNIQUE INDEX external_connections_user_provider_single_uidx
    ON external_connections (user_id, provider)
    WHERE provider <> 'trello';

-- +goose Down
DROP INDEX external_connections_user_provider_single_uidx;

ALTER TABLE external_connections ADD CONSTRAINT external_connections_user_id_provider_key
    UNIQUE (user_id, provider);
