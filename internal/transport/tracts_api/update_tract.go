package tracts_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (t *TractsImpl) UpdateTract(ctx context.Context, req *pb.UpdateTract_Request) (*pb.UpdateTract_Response, error) {
	id, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing tract uuid")
	}

	def, err := definitionFromProto(req.Definition)
	if err != nil {
		return nil, err
	}

	updated, warnings, err := t.tractSvc.UpdateTract(ctx, id, req.Name, req.Description, def)
	if err != nil {
		return nil, rerrors.Wrap(err, "error updating tract")
	}

	item := tractToProto(updated)

	resp := &pb.UpdateTract_Response{
		Tract:    item,
		Warnings: warnings,
	}

	return resp, nil
}
