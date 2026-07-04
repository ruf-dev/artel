package task_trackers_api

import (
	"context"

	"github.com/google/uuid"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service/user_errors"
	"go.redsock.ru/rerrors"
)

func (t *TaskTrackersImpl) ListTrelloBoards(
	ctx context.Context,
	req *pb.ListTrelloBoards_Request,
) (*pb.ListTrelloBoards_Response, error) {
	trackerUuid, err := uuid.Parse(req.TrackerId)
	if err != nil {
		return nil, rerrors.Wrap(user_errors.NotFound, "error parsing tracker id")
	}

	boards, err := t.trackerSvc.ListTrelloBoards(ctx, trackerUuid)
	if err != nil {
		return nil, rerrors.Wrap(err, "error listing trello boards")
	}

	return &pb.ListTrelloBoards_Response{Boards: boardsToProto(boards)}, nil
}
