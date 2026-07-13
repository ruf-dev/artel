package mailserversuggestions

import (
	"context"
	"database/sql"

	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/pg_err"
	"go.redsock.ru/rerrors"
)

type Repo struct {
	q *artel_q.Queries
}

func New(q *artel_q.Queries) *Repo {
	return &Repo{q: q}
}

func (r *Repo) ListByDomain(ctx context.Context, domainPrefix string) ([]domain.MailServerSuggestion, error) {
	rows, err := r.q.ListMailServerSuggestions(ctx, sql.NullString{String: domainPrefix, Valid: true})
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "list mail server suggestions")
	}

	suggestions := make([]domain.MailServerSuggestion, len(rows))
	for i, row := range rows {
		suggestions[i] = domain.MailServerSuggestion{
			Domain:         row.Domain,
			Smtp:           row.Smtp,
			SmtpPort:       int(row.SmtpPort),
			Imap:           row.Imap,
			ImapPort:       int(row.ImapPort),
			AppPasswordUrl: row.AppPasswordUrl.String,
		}
	}

	return suggestions, nil
}
