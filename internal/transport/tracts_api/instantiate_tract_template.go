package tracts_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (t *TractsImpl) InstantiateTractTemplate(
	ctx context.Context,
	req *pb.InstantiateTractTemplate_Request,
) (*pb.InstantiateTractTemplate_Response, error) {
	templateUuid, err := uuid.Parse(req.TemplateUuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing template uuid")
	}

	connections, err := connectionsFromProto(req.Connections)
	if err != nil {
		return nil, err
	}

	created, warnings, err := t.tractSvc.InstantiateTemplate(ctx, templateUuid, req.Name, req.Description, connections)
	if err != nil {
		return nil, rerrors.Wrap(err, "error instantiating tract template")
	}

	item := tractToProto(created)

	resp := &pb.InstantiateTractTemplate_Response{
		Tract:    item,
		Warnings: warnings,
	}

	return resp, nil
}

// connectionsFromProto parses the incoming mcp-name -> connection_uuid string map into
// map[string]uuid.UUID, rejecting any malformed uuid value with a clean InvalidArgument error.
func connectionsFromProto(raw map[string]string) (map[string]uuid.UUID, error) {
	if len(raw) == 0 {
		return nil, nil
	}

	connections := make(map[string]uuid.UUID, len(raw))

	for mcpName, rawUuid := range raw {
		connUuid, err := uuid.Parse(rawUuid)
		if err != nil {
			return nil, rerrors.Wrap(user_errors.TractConnectionUuidInvalid, mcpName)
		}

		connections[mcpName] = connUuid
	}

	return connections, nil
}
