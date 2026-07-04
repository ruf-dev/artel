package mcp_keys_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (m *McpKeysImpl) RemoveMcpConnector(ctx context.Context, req *pb.RemoveMcpConnector_Request) (*pb.RemoveMcpConnector_Response, error) {
	keyID, err := uuid.Parse(req.KeyId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse key id")
	}

	err = m.mcpSvc.RemoveConnector(ctx, keyID, req.McpName)
	if err != nil {
		return nil, rerrors.Wrap(err, "remove mcp connector")
	}

	return &pb.RemoveMcpConnector_Response{}, nil
}
