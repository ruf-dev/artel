package tracts_api

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
)

func (t *TractsImpl) DeleteTract(ctx context.Context, req *pb.DeleteTract_Request) (*pb.DeleteTract_Response, error) {
	id, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing tract uuid")
	}

	err = t.tractSvc.DeleteTract(ctx, id)
	if err != nil {
		return nil, rerrors.Wrap(err, "error deleting tract")
	}

	return &pb.DeleteTract_Response{}, nil
}
