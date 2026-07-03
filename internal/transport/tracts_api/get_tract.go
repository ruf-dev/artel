package tracts_api

import (
	"context"

	"github.com/google/uuid"
	"go.redsock.ru/rerrors"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
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

	item, err := tractToProto(found)
	if err != nil {
		return nil, err
	}

	resp := &pb.GetTract_Response{Tract: item}
	return resp, nil
}
