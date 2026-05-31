package pg

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type loggingDB struct {
	db *sql.DB
}

func newLoggingDB(db *sql.DB) *loggingDB {
	return &loggingDB{db: db}
}

func (l *loggingDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	res, err := l.db.ExecContext(ctx, query, args...)
	elapsed := time.Since(start)
	q := trimQuery(query)
	if err != nil {
		log.Error().Err(err).Str("query", q).Dur("dur", elapsed).Msg("pg")
	} else {
		log.Debug().Str("query", q).Dur("dur", elapsed).Msg("pg")
	}
	return res, err
}

func (l *loggingDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return l.db.PrepareContext(ctx, query)
}

func (l *loggingDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := l.db.QueryContext(ctx, query, args...)
	elapsed := time.Since(start)
	q := trimQuery(query)
	if err != nil {
		log.Error().Err(err).Str("query", q).Dur("dur", elapsed).Msg("pg")
	} else {
		log.Debug().Str("query", q).Dur("dur", elapsed).Msg("pg")
	}
	return rows, err
}

func (l *loggingDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	start := time.Now()
	row := l.db.QueryRowContext(ctx, query, args...)
	elapsed := time.Since(start)
	log.Debug().Str("query", trimQuery(query)).Dur("dur", elapsed).Msg("pg")
	return row
}

func trimQuery(q string) string {
	q = strings.Join(strings.Fields(q), " ")
	if len(q) > 120 {
		q = q[:120]
	}
	return q
}
