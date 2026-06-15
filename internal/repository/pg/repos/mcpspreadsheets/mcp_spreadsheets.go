package mcpspreadsheets

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	"github.com/ruf-dev/artel/internal/domain"
	artel_q "github.com/ruf-dev/artel/internal/repository/pg/generated"
	"github.com/ruf-dev/artel/internal/repository/pg/pg_err"
)

type Repo struct {
	q *artel_q.Queries
}

func New(q *artel_q.Queries) *Repo {
	return &Repo{q: q}
}

func (r *Repo) Insert(ctx context.Context, spreadsheet domain.McpSpreadsheet) (domain.McpSpreadsheet, error) {
	params := artel_q.InsertMcpSpreadsheetParams{
		UserID:               spreadsheet.UserUuid,
		ExternalConnectionID: spreadsheet.ExternalConnectionId,
		SpreadsheetID:        spreadsheet.SpreadsheetId,
		Name:                 spreadsheet.Name,
	}

	row, err := r.q.InsertMcpSpreadsheet(ctx, params)
	if err != nil {
		return domain.McpSpreadsheet{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error inserting mcp spreadsheet")
	}

	return toDomain(row), nil
}

func (r *Repo) ListByUser(ctx context.Context, userUuid uuid.UUID) ([]domain.McpSpreadsheet, error) {
	rows, err := r.q.ListMcpSpreadsheetsByUser(ctx, userUuid)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing mcp spreadsheets")
	}

	spreadsheets := make([]domain.McpSpreadsheet, len(rows))
	for i, row := range rows {
		spreadsheets[i] = toDomain(row)
	}
	return spreadsheets, nil
}

func (r *Repo) Delete(ctx context.Context, userUuid uuid.UUID, spreadsheetId string) error {
	params := artel_q.DeleteMcpSpreadsheetParams{
		UserID:        userUuid,
		SpreadsheetID: spreadsheetId,
	}

	err := r.q.DeleteMcpSpreadsheet(ctx, params)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error deleting mcp spreadsheet")
	}
	return nil
}

func toDomain(row artel_q.McpSpreadsheet) domain.McpSpreadsheet {
	return domain.McpSpreadsheet{
		Uuid:                 row.ID,
		UserUuid:             row.UserID,
		ExternalConnectionId: row.ExternalConnectionID,
		SpreadsheetId:        row.SpreadsheetID,
		Name:                 row.Name,
		CreatedAt:            row.CreatedAt,
	}
}
