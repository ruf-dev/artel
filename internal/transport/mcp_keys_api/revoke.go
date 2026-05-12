package mcp_keys_api

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

func (m *McpKeysImpl) RevokeMcpKey(ctx context.Context, req *pb.RevokeMcpKey_Request) (*pb.RevokeMcpKey_Response, error) {
	keyID, err := uuid.Parse(req.KeyId)
	if err != nil {
		return nil, rerrors.Wrap(err, "parse key id")
	}

	err = m.mcpSvc.RevokeKey(ctx, keyID)
	if err != nil {
		return nil, rerrors.Wrap(err, "revoke mcp key")
	}

	resp := &pb.RevokeMcpKey_Response{}
	return resp, nil
}
