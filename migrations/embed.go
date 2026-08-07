package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"
)

//go:embed *.sql
var FS embed.FS

// ApplyMigration runs all pending goose migrations embedded in this package
// against conn. Safe to call on every startup — goose tracks applied
// versions in the goose_db_version table and is a no-op once up to date.
func ApplyMigration(conn *sql.DB) error {
	goose.SetLogger(_logger{})

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

type _logger struct {
}

func (l _logger) Fatalf(format string, v ...interface{}) {
	log.Fatal().
		Str("msg", fmt.Sprintf(format, v...)).
		Msg("Fatal")
}

func (l _logger) Printf(format string, v ...interface{}) {
	log.Info().
		Str("msg", fmt.Sprintf(format, v...)).
		Msg("Print")
}
