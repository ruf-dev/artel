package task_trackers_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

func (t *TaskTrackersImpl) ListTaskTrackers(
	ctx context.Context,
	_ *pb.ListTaskTrackers_Request,
) (*pb.ListTaskTrackers_Response, error) {
	trackers, err := t.trackerSvc.ListTrackers(ctx)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing task trackers")
	}

	items := make([]*pb.TaskTrackerInfo, len(trackers))
	for i, tr := range trackers {
		items[i] = trackerToProto(tr)
	}

	return &pb.ListTaskTrackers_Response{Trackers: items}, nil
}
