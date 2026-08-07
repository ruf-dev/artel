-- +goose Up

ALTER TABLE couch_instances ADD COLUMN owner_user_id UUID REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE s3_instances ADD COLUMN owner_user_id UUID REFERENCES users(id) ON DELETE CASCADE;

-- +goose Down

ALTER TABLE couch_instances DROP COLUMN owner_user_id;
ALTER TABLE s3_instances DROP COLUMN owner_user_id;
