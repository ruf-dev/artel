package migrations

import (
	"database/sql"
	"embed"

	"github.com/pressly/goose/v3"
	"go.redsock.ru/rerrors"
)

//go:embed *.sql
var FS embed.FS

// ApplyMigration runs all pending goose migrations embedded in this package
// against conn. Safe to call on every startup — goose tracks applied
// versions in the goose_db_version table and is a no-op once up to date.
func ApplyMigration(conn *sql.DB) error {
	goose.SetLogger(goose.NopLogger())

	err := goose.SetDialect("postgres")
	if err != nil {
		return rerrors.Wrap(err, "error setting dialect")
	}

	goose.SetBaseFS(FS)

	err = goose.Up(conn, ".")
	if err != nil {
		return rerrors.Wrap(err, "error applying migrations")
	}

	return nil
}
