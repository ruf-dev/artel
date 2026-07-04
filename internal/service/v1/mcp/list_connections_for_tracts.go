package mcp

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"go.redsock.ru/rerrors"
)

func (s *McpServiceImpl) listConnectionsForTracts(ctx context.Context, userUuid uuid.UUID) (domain.ToolExecResult, error) {
	conns, err := s.externalConnections.ListByUser(ctx, userUuid)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error listing external connections")
	}

	rows := make([]map[string]string, 0, len(conns))

	for _, c := range conns {
		row := map[string]string{
			"uuid":     c.Uuid.String(),
			"provider": c.Provider,
		}
		rows = append(rows, row)
	}

	text, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "error marshaling connections")
	}

	return domain.ToolExecResult{Text: string(text)}, nil
}
