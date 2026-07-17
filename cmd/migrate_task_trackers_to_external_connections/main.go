// Command migrate_task_trackers_to_external_connections is a one-off, manually-run data
// migration. It backfills every row of the (about to be dropped) task_trackers table into
// external_connections (provider = "trello"), decrypting each row's api_key_enc/api_token_enc
// with cryptoutil and re-encrypting them through the normal external connections repo so the
// result is indistinguishable from a connection added via AddTrelloConnection.
//
// Run this against the target database BEFORE deploying migration
// 050_drop_task_trackers.sql — that migration drops task_trackers unconditionally and does not
// migrate rows itself (see its header comment for why: credentials_enc is Go-side AES-GCM
// ciphertext, so the copy can't be expressed in pure SQL).
//
// This is deliberately not wired into goose: it must run exactly once, by a human, against one
// database, before the schema-dropping migration is applied — not automatically on every
// process start the way migrations/*.sql are.
//
// Usage: go run ./cmd/migrate_task_trackers_to_external_connections
// Reads the same config file/env the main service reads (internal/config.Init) for the DB
// connection string and CredsEncryptionKey, but never calls goose — it deliberately opens the
// DB connection itself so it cannot accidentally trigger a pending goose migration (in
// particular the one that drops task_trackers) before it has read the table.
package main

import (
	"context"
	"database/sql"
	"encoding/hex"
	"encoding/json"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"

	"github.com/ruf-dev/artel/internal/config"
	"github.com/ruf-dev/artel/internal/cryptoutil"
	"github.com/ruf-dev/artel/internal/domain"
	repopg "github.com/ruf-dev/artel/internal/repository/pg"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/utils"
	"go.redsock.ru/rerrors"
)

// legacyTaskTracker mirrors the row shape of the about-to-be-dropped task_trackers table.
// Kept local to this command (rather than importing the deleted
// internal/repository/pg/repos/tasktrackers package) since that package is removed as dead code
// in the same change that adds this command.
type legacyTaskTracker struct {
	uuid        uuid.UUID
	userUuid    uuid.UUID
	name        string
	apiKeyEnc   []byte
	apiTokenEnc []byte
}

// trelloConnectionMeta mirrors the unexported type of the same name in
// internal/service/v1/externalconnections/external_connections.go — duplicated here since that
// type isn't exported, but the two must stay in sync: this is the metadata shape
// AddTrelloConnection writes, and ListTrackers expects to read back.
type trelloConnectionMeta struct {
	FullName string `json:"full_name"`
}

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("error reading config")
	}

	encKey, err := hex.DecodeString(cfg.Environment.CredsEncryptionKey)
	if err != nil {
		log.Fatal().Err(err).Msg("error decoding creds_encryption_key")
	}

	// Deliberately sql.Open, not sqldb.New: sqldb.New runs goose.Up as a side effect of
	// connecting, which would apply 050_drop_task_trackers.sql (if already present in
	// migrations/) before this command gets a chance to read the table.
	db, err := sql.Open(cfg.DataSources.Postgres.SqlDialect(), cfg.DataSources.Postgres.ConnectionString())
	if err != nil {
		log.Fatal().Err(err).Msg("error connecting to postgres")
	}

	ctx := context.Background()

	err = run(ctx, db, encKey)
	if err != nil {
		log.Fatal().Err(err).Msg("error migrating task trackers")
	}

	log.Info().Msg("task_trackers migration complete")
}

func run(ctx context.Context, db *sql.DB, encKey []byte) error {
	trackers, err := listLegacyTaskTrackers(ctx, db)
	if err != nil {
		return rerrors.Wrap(err, "error listing legacy task trackers")
	}

	repos := repopg.New(db, encKey)
	connections := repos.ExternalConnections()

	migrated := 0

	for _, tracker := range trackers {
		conn, err := toExternalConnection(encKey, tracker)
		if err != nil {
			return rerrors.Wrap(err, "error decrypting legacy task tracker")
		}

		_, err = connections.Insert(ctx, conn)
		if err != nil {
			return rerrors.Wrap(err, "error inserting external connection")
		}

		migrated++
	}

	log.Info().Int("count", migrated).Msg("migrated task_trackers rows to external_connections")

	return nil
}

func listLegacyTaskTrackers(ctx context.Context, db *sql.DB) ([]legacyTaskTracker, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT id, user_id, name, api_key_enc, api_token_enc
		FROM task_trackers
	`)
	if err != nil {
		return nil, rerrors.Wrap(err, "error querying task_trackers")
	}
	defer utils.CloseWithLog(rows, "task_trackers migration rows")

	var trackers []legacyTaskTracker

	for rows.Next() {
		var tracker legacyTaskTracker

		err = rows.Scan(&tracker.uuid, &tracker.userUuid, &tracker.name, &tracker.apiKeyEnc, &tracker.apiTokenEnc)
		if err != nil {
			return nil, rerrors.Wrap(err, "error scanning task_trackers row")
		}

		trackers = append(trackers, tracker)
	}

	err = rows.Err()
	if err != nil {
		return nil, rerrors.Wrap(err, "error iterating task_trackers rows")
	}

	return trackers, nil
}

func toExternalConnection(encKey []byte, tracker legacyTaskTracker) (domain.ExternalConnection, error) {
	apiKey, err := cryptoutil.Decrypt(encKey, tracker.apiKeyEnc)
	if err != nil {
		return domain.ExternalConnection{}, rerrors.Wrap(err, "error decrypting api key")
	}

	apiToken, err := cryptoutil.Decrypt(encKey, tracker.apiTokenEnc)
	if err != nil {
		return domain.ExternalConnection{}, rerrors.Wrap(err, "error decrypting api token")
	}

	creds := domain.APIKeyCredentials{
		APIKey:   string(apiKey),
		APIToken: string(apiToken),
	}

	credJSON, err := json.Marshal(creds)
	if err != nil {
		return domain.ExternalConnection{}, rerrors.Wrap(err, "error marshaling credentials")
	}

	meta := trelloConnectionMeta{FullName: tracker.name}

	metaJSON, err := json.Marshal(meta)
	if err != nil {
		return domain.ExternalConnection{}, rerrors.Wrap(err, "error marshaling metadata")
	}

	conn := domain.ExternalConnection{
		UserUuid:        tracker.userUuid,
		Provider:        domain.ProviderTrello,
		ProviderType:    artel_q.ExternalProviderTypeApiKey,
		CredentialsJSON: credJSON,
		Metadata:        metaJSON,
	}

	return conn, nil
}
