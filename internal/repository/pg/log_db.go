package pg

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
)

var pgTracer = otel.Tracer("artel/pg")

type loggingDB struct {
	db *sql.DB
}

func newLoggingDB(db *sql.DB) *loggingDB {
	return &loggingDB{db: db}
}

func (l *loggingDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	q := trimQuery(query)
	ctx, span := pgTracer.Start(ctx, "pg.exec")
	defer span.End()

	start := time.Now()
	res, err := l.db.ExecContext(ctx, query, args...)
	elapsed := time.Since(start)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		log.Ctx(ctx).Error().Err(err).Str("query", q).Dur("dur", elapsed).Msg("pg")
	} else {
		log.Ctx(ctx).Debug().Str("query", q).Dur("dur", elapsed).Msg("pg")
	}
	return res, err
}

func (l *loggingDB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return l.db.PrepareContext(ctx, query)
}

func (l *loggingDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	q := trimQuery(query)
	ctx, span := pgTracer.Start(ctx, "pg.query")
	defer span.End()

	start := time.Now()
	rows, err := l.db.QueryContext(ctx, query, args...)
	elapsed := time.Since(start)

	if err != nil {
		span.SetStatus(codes.Error, err.Error())
		log.Ctx(ctx).Error().Err(err).Str("query", q).Dur("dur", elapsed).Msg("pg")
	} else {
		log.Ctx(ctx).Debug().Str("query", q).Dur("dur", elapsed).Msg("pg")
	}
	return rows, err
}

func (l *loggingDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	q := trimQuery(query)
	ctx, span := pgTracer.Start(ctx, "pg.query_row")
	defer span.End()

	start := time.Now()
	row := l.db.QueryRowContext(ctx, query, args...)
	elapsed := time.Since(start)
	log.Ctx(ctx).Debug().Str("query", q).Dur("dur", elapsed).Msg("pg")
	return row
}

func trimQuery(q string) string {
	q = strings.Join(strings.Fields(q), " ")
	if len(q) > 120 {
		q = q[:120]
	}
	return q
}
