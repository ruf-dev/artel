-- +goose Up

ALTER TABLE system_settings ADD COLUMN default_docs_source TEXT NOT NULL DEFAULT 'vault'
    CHECK (default_docs_source IN ('vault', 'github'));

-- +goose Down

ALTER TABLE system_settings DROP COLUMN default_docs_source;
