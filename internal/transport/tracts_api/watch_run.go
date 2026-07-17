package tracts_api

import (
	"time"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

const watchRunPollInterval = 700 * time.Millisecond

// WatchRun streams periodic snapshots of a run's status and steps until the run reaches a
// terminal state (done/failed) or the client disconnects. It is a thin polling wrapper around
// GetRun — no pubsub/event hook exists in the engine, and this is a transport-layer concern only.
func (t *TractsImpl) WatchRun(req *pb.WatchRun_Request, stream pb.TractsAPI_WatchRunServer) error {
	runUuid, err := uuid.Parse(req.RunUuid)
	if err != nil {
		return rerrors.Wrap(user_errors.NotFound, "error parsing run uuid")
	}

	for {
		run, steps, err := t.tractSvc.GetRun(stream.Context(), runUuid)
		if err != nil {
			return rerrors.Wrap(err, "error getting tract run")
		}

		resp := &pb.WatchRun_Response{
			Run:   runToProto(run),
			Steps: runStepsToProto(steps),
		}

		err = stream.Send(resp)
		if err != nil {
			return rerrors.Wrap(err, "error sending run snapshot")
		}

		if run.Status != domain.TractRunRunning {
			return nil
		}

		select {
		case <-stream.Context().Done():
			return nil
		case <-time.After(watchRunPollInterval):
		}
	}
}
