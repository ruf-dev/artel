package tracts_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (t *TractsImpl) LinkTrigger(ctx context.Context, req *pb.LinkTrigger_Request) (*pb.LinkTrigger_Response, error) {
	triggerUuid, err := uuid.Parse(req.TriggerUuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing trigger uuid")
	}

	tractUuid, err := uuid.Parse(req.TractUuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing tract uuid")
	}

	filters, err := filtersFromJSON(req.Filters)
	if err != nil {
		return nil, err
	}

	err = t.tractSvc.LinkTrigger(ctx, triggerUuid, tractUuid, filters)
	if err != nil {
		return nil, rerrors.Wrap(err, "error linking trigger to tract")
	}

	return &pb.LinkTrigger_Response{}, nil
}
