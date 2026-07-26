-- +goose Up

ALTER TABLE docker_hosts ADD COLUMN ca_cert_enc     BYTEA;
ALTER TABLE docker_hosts ADD COLUMN client_cert_enc BYTEA;
ALTER TABLE docker_hosts ADD COLUMN client_key_enc  BYTEA;

-- +goose Down

ALTER TABLE docker_hosts DROP COLUMN ca_cert_enc;
ALTER TABLE docker_hosts DROP COLUMN client_cert_enc;
ALTER TABLE docker_hosts DROP COLUMN client_key_enc;
