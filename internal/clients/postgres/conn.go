package postgres

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"
	"go.redsock.ru/toolbox/closer"
	"go.vervstack.ru/matreshka/pkg/matreshka/resources"

	"github.com/pressly/goose/v3"

	"github.com/ruf-dev/artel/internal/utils"
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

	goose.SetLogger(sqlLogger{})

	err = goose.SetDialect("postgres")
	if err != nil {
		return rerrors.Wrap(err, "error setting dialect")
	}

	mig := cfg.MigrationFolder()
	if mig == "" {
		mig = "./migrations"
	}

	err = goose.Up(conn, mig)
	if err != nil {
		return rerrors.Wrap(err, "error performing up")
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

type sqlLogger struct{}

func (s sqlLogger) Fatalf(format string, v ...interface{}) {
	log.Fatal().Msgf(format, v...)
}

func (s sqlLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
