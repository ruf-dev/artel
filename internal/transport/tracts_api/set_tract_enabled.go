package tracts_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (t *TractsImpl) SetTractEnabled(
	ctx context.Context,
	req *pb.SetTractEnabled_Request,
) (*pb.SetTractEnabled_Response, error) {
	id, err := uuid.Parse(req.Uuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing tract uuid")
	}

	err = t.tractSvc.SetTractEnabled(ctx, id, req.Enabled)
	if err != nil {
		return nil, rerrors.Wrap(err, "error setting tract enabled")
	}

	return &pb.SetTractEnabled_Response{}, nil
}
