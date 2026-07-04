package mcp_keys_api

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (m *McpKeysImpl) ExecuteMomTool(ctx context.Context, req *pb.ExecuteMomTool_Request) (*pb.ExecuteMomTool_Response, error) {
	exConnUuid, err := uuid.Parse(req.ExternalConnectionId)
	if err != nil {
		return nil, rerrors.Wrap(err, "error parsing external connection id")
	}

	params := map[string]interface{}{}
	if req.ParamsJson != "" {
		err = json.Unmarshal([]byte(req.ParamsJson), &params)
		if err != nil {
			return nil, rerrors.Wrap(err, "error unmarshaling params json")
		}
	}

	result, err := m.momSvc.ExecuteToolForUserConnection(ctx, exConnUuid, req.McpName, req.ToolName, params)
	if err != nil {
		return nil, rerrors.Wrap(err, "error executing mom tool")
	}

	resp := &pb.ExecuteMomTool_Response{Result: result}

	return resp, nil
}
