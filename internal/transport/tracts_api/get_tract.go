package tracts_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (t *TractsImpl) GetTract(ctx context.Context, req *pb.GetTract_Request) (*pb.GetTract_Response, error) {
	id, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing tract uuid")
	}

	found, err := t.tractSvc.GetTract(ctx, id)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting tract")
	}

	item := tractToProto(found)

	resp := &pb.GetTract_Response{Tract: item}

	return resp, nil
}
