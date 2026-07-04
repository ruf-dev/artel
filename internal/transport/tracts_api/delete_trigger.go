package tracts_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (t *TractsImpl) DeleteTrigger(ctx context.Context, req *pb.DeleteTrigger_Request) (*pb.DeleteTrigger_Response, error) {
	id, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing trigger uuid")
	}

	err = t.tractSvc.DeleteTrigger(ctx, id)
	if err != nil {
		return nil, rerrors.Wrap(err, "error deleting trigger")
	}

	return &pb.DeleteTrigger_Response{}, nil
}
