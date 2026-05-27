-- +goose Up

-- Restructure user_permissions: user_id becomes PK, drop id/timestamps
ALTER TABLE user_permissions DROP COLUMN IF EXISTS created_at;
ALTER TABLE user_permissions DROP COLUMN IF EXISTS updated_at;
ALTER TABLE user_permissions DROP CONSTRAINT IF EXISTS user_permissions_pkey;
ALTER TABLE user_permissions DROP CONSTRAINT IF EXISTS user_permissions_user_id_key;
ALTER TABLE user_permissions DROP COLUMN IF EXISTS id;
ALTER TABLE user_permissions ADD PRIMARY KEY (user_id);

-- Restructure subscriptions: user_id becomes PK, drop id/timestamps
ALTER TABLE subscriptions DROP COLUMN IF EXISTS created_at;
ALTER TABLE subscriptions DROP COLUMN IF EXISTS updated_at;
ALTER TABLE subscriptions DROP CONSTRAINT subscriptions_pkey;
ALTER TABLE subscriptions DROP COLUMN IF EXISTS id;
ALTER TABLE subscriptions ADD PRIMARY KEY (user_id);

-- Backfill missing records for existing users
INSERT INTO user_permissions (user_id, is_administrator)
SELECT id, FALSE FROM users
WHERE id NOT IN (SELECT user_id FROM user_permissions);

INSERT INTO subscriptions (user_id, active)
SELECT id, FALSE FROM users
WHERE id NOT IN (SELECT user_id FROM subscriptions);

-- +goose Down

ALTER TABLE subscriptions DROP CONSTRAINT subscriptions_pkey;
ALTER TABLE subscriptions ADD COLUMN id UUID NOT NULL DEFAULT gen_random_uuid();
ALTER TABLE subscriptions ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE subscriptions ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE subscriptions ADD PRIMARY KEY (id);
ALTER TABLE subscriptions ADD UNIQUE(user_id);

ALTER TABLE user_permissions DROP CONSTRAINT user_permissions_pkey;
ALTER TABLE user_permissions ADD COLUMN id UUID NOT NULL DEFAULT gen_random_uuid();
ALTER TABLE user_permissions ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE user_permissions ADD COLUMN updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW();
ALTER TABLE user_permissions ADD PRIMARY KEY (id);
ALTER TABLE user_permissions ADD UNIQUE(user_id);
