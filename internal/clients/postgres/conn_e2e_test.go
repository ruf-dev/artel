//go:build e2e
// +build e2e

package postgres_test

import (
	"database/sql"
	"os"
	"strconv"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"

	"github.com/ruf-dev/artel/internal/clients/postgres"
)

// envOrDefault mirrors tests/harness's helper of the same name, kept local here so this file
// stays self-contained and doesn't reach into harness's unexported helpers.
func envOrDefault(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	return val
}

// testPostgresConfig builds a *resources.Postgres pointed at the shared e2e Postgres instance
// (tests/docker-compose.yaml's test-postgres service), matching tests/harness.go's PG_HOST/
// PG_PORT/PG_ADMIN_DATABASE/PG_ADMIN_USER/PG_ADMIN_PASSWORD conventions. AdminUser/AdminPwd are
// left empty so postgres.Migrate's AdminConnectionString() falls back to User/Pwd.
func testPostgresConfig(t *testing.T, schema string) *resources.Postgres {
	t.Helper()

	port, err := strconv.ParseUint(envOrDefault("PG_PORT", "15434"), 10, 64)
	require.NoError(t, err, "parse PG_PORT")

	cfg := &resources.Postgres{
		Host:    envOrDefault("PG_HOST", "localhost"),
		Port:    port,
		User:    envOrDefault("PG_ADMIN_USER", "artel"),
		Pwd:     envOrDefault("PG_ADMIN_PASSWORD", "artel_db"),
		DbName:  envOrDefault("PG_ADMIN_DATABASE", "artel_db"),
		SslMode: "disable",
		Schema:  schema,
	}

	return cfg
}

// rawConn opens a direct connection to the shared e2e Postgres instance, independent of the
// package under test, for setup/teardown and assertions.
func rawConn(t *testing.T) *sql.DB {
	t.Helper()

	cfg := testPostgresConfig(t, "")

	db, err := sql.Open("postgres", cfg.ConnectionString())
	require.NoError(t, err, "open raw postgres connection")

	t.Cleanup(func() {
		_ = db.Close()
	})

	return db
}

// schemaExists reports whether schema is registered in pg_catalog.pg_namespace.
func schemaExists(t *testing.T, db *sql.DB, schema string) bool {
	t.Helper()

	var exists bool

	err := db.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_namespace WHERE nspname = $1)`, schema,
	).Scan(&exists)
	require.NoError(t, err, "check schema existence")

	return exists
}

// TestMigrate_NonexistentSchema_FailsFast confirms Migrate refuses to run migrations against a
// schema that doesn't exist, and does not auto-create it either.
func TestMigrate_NonexistentSchema_FailsFast(t *testing.T) {
	db := rawConn(t)

	schema := "nonexistent_" + uuid.NewString()

	cfg := testPostgresConfig(t, schema)

	err := postgres.Migrate(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), schema)

	require.False(t, schemaExists(t, db, schema), "Migrate must not auto-create the missing schema")
}

// TestMigrate_DefaultPublicSchema_Succeeds confirms Migrate still works against the implicit
// "public" schema, the same schema tests/bootstrap already migrates. goose is idempotent, so
// running it again here is safe.
func TestMigrate_DefaultPublicSchema_Succeeds(t *testing.T) {
	cfg := testPostgresConfig(t, "")

	err := postgres.Migrate(cfg)
	require.NoError(t, err)
}

// TestMigrate_ExplicitPreCreatedSchema_TablesLandThere confirms Migrate against an explicitly
// pre-created, non-default schema succeeds and that migration tables land there, not in public.
func TestMigrate_ExplicitPreCreatedSchema_TablesLandThere(t *testing.T) {
	db := rawConn(t)

	schema := "artel_test_" + uuid.NewString()[:8]

	_, err := db.Exec(`CREATE SCHEMA IF NOT EXISTS "` + schema + `"`)
	require.NoError(t, err, "create throwaway schema")

	t.Cleanup(func() {
		_, dropErr := db.Exec(`DROP SCHEMA IF EXISTS "` + schema + `" CASCADE`)
		require.NoError(t, dropErr, "drop throwaway schema")
	})

	cfg := testPostgresConfig(t, schema)

	err = postgres.Migrate(cfg)
	require.NoError(t, err)

	var gooseTableExists bool

	err = db.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = $1 AND table_name = 'goose_db_version')`,
		schema,
	).Scan(&gooseTableExists)
	require.NoError(t, err, "check goose_db_version table location")
	require.True(t, gooseTableExists, "goose_db_version should have been created in the throwaway schema")

	var usersTableExists bool

	err = db.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_schema = $1 AND table_name = 'users')`,
		schema,
	).Scan(&usersTableExists)
	require.NoError(t, err, "check users table location")
	require.True(t, usersTableExists, "users table should have been created in the throwaway schema")
}
