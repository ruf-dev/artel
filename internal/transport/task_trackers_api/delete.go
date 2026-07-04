package task_trackers_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (t *TaskTrackersImpl) DeleteTaskTracker(ctx context.Context, req *pb.DeleteTaskTracker_Request) (*pb.DeleteTaskTracker_Response, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing tracker id")
	}

	err = t.trackerSvc.DeleteTracker(ctx, id)
	if err != nil {
		return nil, rerrors.Wrap(err, "error deleting task tracker")
	}

	return &pb.DeleteTaskTracker_Response{}, nil
}
