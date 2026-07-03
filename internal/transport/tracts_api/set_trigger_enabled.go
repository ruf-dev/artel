package tracts_api

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func (t *TractsImpl) SetTriggerEnabled(ctx context.Context, req *pb.SetTriggerEnabled_Request) (*pb.SetTriggerEnabled_Response, error) {
	id, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing trigger uuid")
	}

	err = t.tractSvc.SetTriggerEnabled(ctx, id, req.Enabled)
	if err != nil {
		return nil, rerrors.Wrap(err, "error setting trigger enabled")
	}

	return &pb.SetTriggerEnabled_Response{}, nil
}
