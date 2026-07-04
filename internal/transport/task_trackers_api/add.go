package task_trackers_api

import (
	"context"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/clients/trello"
	"github.com/ruf-dev/artel/internal/domain"
	"go.redsock.ru/rerrors"
)

func (t *TaskTrackersImpl) AddTaskTracker(ctx context.Context, req *pb.AddTaskTracker_Request) (*pb.AddTaskTracker_Response, error) {
	tracker := domain.TaskTracker{
		Type: req.Type,
	}

	creds := trello.TrelloCredentials{
		ApiKey:   req.ApiKey,
		ApiToken: req.ApiToken,
	}

	created, boards, err := t.trackerSvc.AddTracker(ctx, tracker, creds)
	if err != nil {
		return nil, rerrors.Wrap(err, "error adding task tracker")
	}

	resp := &pb.AddTaskTracker_Response{
		Tracker: trackerToProto(created),
		Boards:  boardsToProto(boards),
	}

	return resp, nil
}
