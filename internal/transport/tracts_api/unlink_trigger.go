package tracts_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (t *TractsImpl) UnlinkTrigger(
	ctx context.Context,
	req *pb.UnlinkTrigger_Request,
) (*pb.UnlinkTrigger_Response, error) {
	triggerUuid, err := uuid.Parse(req.TriggerUuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing trigger uuid")
	}

	tractUuid, err := uuid.Parse(req.TractUuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing tract uuid")
	}

	err = t.tractSvc.UnlinkTrigger(ctx, triggerUuid, tractUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error unlinking trigger from tract")
	}

	return &pb.UnlinkTrigger_Response{}, nil
}
