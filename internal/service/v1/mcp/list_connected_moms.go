package mcp

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"go.redsock.ru/rerrors"
)

func (s *McpServiceImpl) listConnectedMoms(ctx context.Context, keyUuid uuid.UUID) (domain.ToolExecResult, error) {
	connectors, err := s.mcpConnectors.ListByKey(ctx, keyUuid)
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "list mcp connectors")
	}

	moms := make([]map[string]string, 0, len(connectors))

	for _, c := range connectors {
		def, err := s.mcpDefinitions.Get(ctx, c.McpName)
		if err != nil {
			return domain.ToolExecResult{}, rerrors.Wrap(err, "get mcp definition")
		}

		if !def.Valid {
			continue
		}

		moms = append(moms, map[string]string{
			"name":        def.V.Name,
			"author":      def.V.Author,
			"description": def.V.Description,
		})
	}

	text, err := json.MarshalIndent(moms, "", "  ")
	if err != nil {
		return domain.ToolExecResult{}, rerrors.Wrap(err, "marshal connected moms")
	}

	return domain.ToolExecResult{Text: string(text)}, nil
}
