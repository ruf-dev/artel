package mcp_keys_api

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (m *McpKeysImpl) ListMcpConnectors(ctx context.Context, req *pb.ListMcpConnectors_Request) (*pb.ListMcpConnectors_Response, error) {
	keyID, err := uuid.Parse(req.KeyId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse key id")
	}

	connectors, err := m.mcpSvc.ListConnectors(ctx, keyID)
	if err != nil {
		return nil, rerrors.Wrap(err, "list mcp connectors")
	}

	infos := make([]*pb.McpConnectorInfo, 0, len(connectors))
	for _, c := range connectors {
		infos = append(infos, &pb.McpConnectorInfo{
			McpKeyId:             c.McpKeyUuid.String(),
			McpName:              c.McpName,
			ExternalConnectionId: c.ExternalConnectionUuid.String(),
			CreatedAt:            c.CreatedAt.String(),
		})
	}

	return &pb.ListMcpConnectors_Response{Connectors: infos}, nil
}
