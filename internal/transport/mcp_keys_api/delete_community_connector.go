package mcp_keys_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (m *McpKeysImpl) DeleteCommunityConnector(
	ctx context.Context,
	req *pb.DeleteCommunityConnector_Request,
) (*pb.DeleteCommunityConnector_Response, error) {
	err := m.mcpSvc.DeleteCommunityConnector(ctx, req.Name)
	if err != nil {
		return nil, rerrors.Wrap(err, "delete community connector")
	}

	return &pb.DeleteCommunityConnector_Response{}, nil
}
