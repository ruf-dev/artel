package task_trackers_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"go.redsock.ru/rerrors"
)

// AddTaskTracker only ever adds a trello tracker today — req.Type is accepted for wire-shape
// compatibility with the frontend but unused, matching AddTracker taking raw credentials rather
// than a domain.TaskTracker (there's nothing else meaningful to set before the connection is
// validated and persisted).
func (t *TaskTrackersImpl) AddTaskTracker(
	ctx context.Context,
	req *pb.AddTaskTracker_Request,
) (*pb.AddTaskTracker_Response, error) {
	created, boards, err := t.trackerSvc.AddTracker(ctx, req.ApiKey, req.ApiToken)
	if err != nil {
		return nil, rerrors.Wrap(err, "error adding task tracker")
	}

	resp := &pb.AddTaskTracker_Response{
		Tracker: trackerToProto(created),
		Boards:  boardsToProto(boards),
	}

	return resp, nil
}
