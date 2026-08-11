package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox/closer"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"

	"github.com/ruf-dev/artel/internal/utils"
	"github.com/ruf-dev/artel/migrations"
)

func Connect(cfg *resources.Postgres) (*sql.DB, error) {
	conn, err := sql.Open("postgres", cfg.ConnectionString())
	if err != nil {
		return nil, rerrors.Wrap(err, "error checking connection to postgres")
	}

	closer.Add(func() error {
		return conn.Close()
	})

	return conn, nil
}

func Migrate(cfg *resources.Postgres) error {
	conn, err := sql.Open("postgres", cfg.AdminConnectionString())
	if err != nil {
		return rerrors.Wrap(err, "error checking connection to postgres")
	}
	defer utils.CloseWithLog(conn, "postgres connection")

	schema := cfg.Schema
	if schema == "" {
		schema = "public"
	}

	err = ensureSchemaExists(conn, schema)
	if err != nil {
		return rerrors.Wrap(err, "error ensuring postgres schema exists")
	}

	return migrations.ApplyMigration(conn)
}

func ensureSchemaExists(conn *sql.DB, schema string) error {
	var exists bool

	err := conn.QueryRow(
		`SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_namespace WHERE nspname = $1)`,
		schema,
	).Scan(&exists)
	if err != nil {
		return rerrors.Wrap(err, "error checking postgres schema existence")
	}

	if !exists {
		return rerrors.New(fmt.Sprintf("postgres schema %q does not exist — create it before starting artel (e.g. CREATE SCHEMA %q)", schema, schema))
	}

	return nil
}

type DB interface {
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)

	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)

	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)

	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}
