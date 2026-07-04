package mcp_keys_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (m *McpKeysImpl) AddMcpConnector(
	ctx context.Context,
	req *pb.AddMcpConnector_Request,
) (*pb.AddMcpConnector_Response, error) {
	keyID, err := uuid.Parse(req.KeyId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse key id")
	}

	externalConnectionID, err := uuid.Parse(req.ExternalConnectionId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse external connection id")
	}

	connector, err := m.mcpSvc.AddConnector(ctx, keyID, req.McpName, externalConnectionID)
	if err != nil {
		return nil, rerrors.Wrap(err, "add mcp connector")
	}

	return &pb.AddMcpConnector_Response{
		Connector: &pb.McpConnectorInfo{
			McpKeyId:             connector.McpKeyUuid.String(),
			McpName:              connector.McpName,
			ExternalConnectionId: connector.ExternalConnectionUuid.String(),
			CreatedAt:            connector.CreatedAt.String(),
		},
	}, nil
}
