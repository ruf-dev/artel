package mcp

import (
	"context"

	"github.com/google/uuid"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/middleware/user_context"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (s *McpServiceImpl) AddConnector(ctx context.Context, keyID uuid.UUID, mcpName string, externalConnectionID uuid.UUID) (domain.McpConnector, error) {
	uc, ok := user_context.GetUserContext(ctx)
	if !ok {
		return domain.McpConnector{}, user_errors.Unauthenticated
	}

	key, err := s.mcpKeys.GetMcpKeyByID(ctx, keyID)
	if err != nil {
		return domain.McpConnector{}, rerrors.Wrap(err, "get mcp key")
	}

	if key.UserUuid != uc.UserUuid {
		return domain.McpConnector{}, rerrors.Wrap(user_errors.NotFound, "mcp key not found")
	}

	def, err := s.mcpDefinitions.Get(ctx, mcpName)
	if err != nil {
		return domain.McpConnector{}, rerrors.Wrap(err, "get mcp definition")
	}

	if !def.Valid {
		return domain.McpConnector{}, rerrors.Wrap(user_errors.NotFound, "mcp definition not found")
	}

	conn, err := s.externalConnections.GetByID(ctx, externalConnectionID)
	if err != nil {
		return domain.McpConnector{}, rerrors.Wrap(err, "get external connection")
	}

	if conn.UserUuid != uc.UserUuid {
		return domain.McpConnector{}, rerrors.Wrap(user_errors.NotFound, "external connection not found")
	}

	connector, err := s.mcpConnectors.Insert(ctx, domain.McpConnector{
		McpKeyUuid:             keyID,
		McpName:                mcpName,
		ExternalConnectionUuid: externalConnectionID,
	})
	if err != nil {
		return domain.McpConnector{}, rerrors.Wrap(err, "insert mcp connector")
	}

	return connector, nil
}
