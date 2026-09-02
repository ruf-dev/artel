-- +goose Up

ALTER TABLE user_settings ADD COLUMN liked_openrouter_models TEXT[] NOT NULL DEFAULT '{}';
ALTER TABLE user_settings ADD COLUMN last_used_model TEXT NOT NULL DEFAULT '';

-- +goose Down

ALTER TABLE user_settings DROP COLUMN last_used_model;
ALTER TABLE user_settings DROP COLUMN liked_openrouter_models;
