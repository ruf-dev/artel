-- +goose Up

ALTER TABLE vaults
    ADD COLUMN s3_instance_id UUID REFERENCES s3_instances (id),
    ADD COLUMN s3_bucket_name TEXT,
    ADD CONSTRAINT vaults_s3_pair_chk CHECK ((s3_instance_id IS NULL) = (s3_bucket_name IS NULL));

-- +goose Down

ALTER TABLE vaults
    DROP CONSTRAINT IF EXISTS vaults_s3_pair_chk,
    DROP COLUMN IF EXISTS s3_bucket_name,
    DROP COLUMN IF EXISTS s3_instance_id;
