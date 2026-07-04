package prompts

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"

	"github.com/ruf-dev/artel/internal/clients/sqldb"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/repository"
	"github.com/ruf-dev/artel/internal/utils"

	"go.redsock.ru/rerrors"
)

type Repo struct {
	db sqldb.DB
}

func New(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (r *Repo) List(ctx context.Context, params repository.ListPromptsParams) ([]domain.Prompt, int64, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	base := psql.Select("id::text", "prompt").From("prompts")
	countBase := psql.Select("COUNT(*)").From("prompts")

	if len(params.Ids) > 0 {
		base = base.Where(sq.Eq{"id::text": params.Ids})
		countBase = countBase.Where(sq.Eq{"id::text": params.Ids})
	}

	if params.Limit > 0 {
		base = base.Limit(uint64(params.Limit))
	}

	base = base.Offset(uint64(params.Offset))

	countSQL, countArgs, err := countBase.ToSql()
	if err != nil {
		return nil, 0, rerrors.Wrap(err, "build count query")
	}

	var total int64

	err = r.db.QueryRowContext(ctx, countSQL, countArgs...).Scan(&total)
	if err != nil {
		return nil, 0, rerrors.Wrap(err, "count prompts")
	}

	listSQL, listArgs, err := base.ToSql()
	if err != nil {
		return nil, 0, rerrors.Wrap(err, "build list query")
	}

	rows, err := r.db.QueryContext(ctx, listSQL, listArgs...)
	if err != nil {
		return nil, 0, rerrors.Wrap(err, "list prompts")
	}
	defer utils.CloseWithLog(rows, "error closing list query (prompts)")

	var result []domain.Prompt

	for rows.Next() {
		var p domain.Prompt

		err = rows.Scan(&p.Id, &p.Text)
		if err != nil {
			return nil, 0, rerrors.Wrap(err, "scan prompt")
		}

		result = append(result, p)
	}

	return result, total, nil
}
