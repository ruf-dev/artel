package tracts_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (t *TractsImpl) ListRuns(ctx context.Context, req *pb.ListRuns_Request) (*pb.ListRuns_Response, error) {
	tractUuid, err := uuid.Parse(req.TractUuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing tract uuid")
	}

	runs, err := t.tractSvc.ListRuns(ctx, tractUuid, req.Limit)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing tract runs")
	}

	resp := &pb.ListRuns_Response{Runs: runsToProto(runs)}

	return resp, nil
}
