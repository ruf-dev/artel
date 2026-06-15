package mcpconnectors

import (
	"context"
	"database/sql"
	"errors"

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

func (r *Repo) Insert(ctx context.Context, connector domain.McpConnector) (domain.McpConnector, error) {
	params := artel_q.InsertMcpConnectorParams{
		McpKeyID:             connector.McpKeyUuid,
		McpName:              connector.McpName,
		ExternalConnectionID: connector.ExternalConnectionUuid,
	}

	row, err := r.q.InsertMcpConnector(ctx, params)
	if err != nil {
		return domain.McpConnector{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error inserting mcp connector")
	}

	return toDomain(row), nil
}

func (r *Repo) Get(ctx context.Context, mcpKeyUuid uuid.UUID, mcpName string) (sql.Null[domain.McpConnector], error) {
	params := artel_q.GetMcpConnectorParams{
		McpKeyID: mcpKeyUuid,
		McpName:  mcpName,
	}

	row, err := r.q.GetMcpConnector(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sql.Null[domain.McpConnector]{}, nil
		}
		return sql.Null[domain.McpConnector]{}, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error getting mcp connector")
	}

	result := sql.Null[domain.McpConnector]{V: toDomain(row), Valid: true}
	return result, nil
}

func (r *Repo) ListByKey(ctx context.Context, mcpKeyUuid uuid.UUID) ([]domain.McpConnector, error) {
	rows, err := r.q.ListMcpConnectorsByKey(ctx, mcpKeyUuid)
	if err != nil {
		return nil, rerrors.Wrap(pg_err.UnwrapPgErr(err), "error listing mcp connectors")
	}

	connectors := make([]domain.McpConnector, len(rows))
	for i, row := range rows {
		connectors[i] = toDomain(row)
	}
	return connectors, nil
}

func (r *Repo) Delete(ctx context.Context, mcpKeyUuid uuid.UUID, mcpName string) error {
	params := artel_q.DeleteMcpConnectorParams{
		McpKeyID: mcpKeyUuid,
		McpName:  mcpName,
	}

	err := r.q.DeleteMcpConnector(ctx, params)
	if err != nil {
		return rerrors.Wrap(pg_err.UnwrapPgErr(err), "error deleting mcp connector")
	}
	return nil
}

func toDomain(row artel_q.McpConnector) domain.McpConnector {
	return domain.McpConnector{
		Uuid:                   row.ID,
		McpKeyUuid:             row.McpKeyID,
		McpName:                row.McpName,
		ExternalConnectionUuid: row.ExternalConnectionID,
		CreatedAt:              row.CreatedAt,
	}
}
