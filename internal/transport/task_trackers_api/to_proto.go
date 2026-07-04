package task_trackers_api

import (
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/domain"
)

func trackerToProto(t domain.TaskTracker) *pb.TaskTrackerInfo {
	return &pb.TaskTrackerInfo{
		Id:        t.Uuid.String(),
		Type:      t.Type,
		Name:      t.Name,
		CreatedAt: t.CreatedAt.UTC().Format("2006-01-02T15:04:05Z"),
	}
}

func boardToProto(b domain.TrelloBoard) *pb.TrelloBoardInfo {
	return &pb.TrelloBoardInfo{
		Id:   b.Id,
		Name: b.Name,
	}
}

func boardsToProto(boards []domain.TrelloBoard) []*pb.TrelloBoardInfo {
	result := make([]*pb.TrelloBoardInfo, len(boards))
	for i, b := range boards {
		result[i] = boardToProto(b)
	}

	return result
}
