package tracts_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (t *TractsImpl) GetRun(ctx context.Context, req *pb.GetRun_Request) (*pb.GetRun_Response, error) {
	runUuid, err := uuid.Parse(req.RunUuid)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing run uuid")
	}

	run, steps, err := t.tractSvc.GetRun(ctx, runUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error getting tract run")
	}

	resp := &pb.GetRun_Response{
		Run:   runToProto(run),
		Steps: runStepsToProto(steps),
	}

	return resp, nil
}
