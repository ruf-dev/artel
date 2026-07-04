package task_trackers_api

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
	"google.golang.org/grpc"
)

type TaskTrackersImpl struct {
	pb.UnimplementedTaskTrackersAPIServer
	trackerSvc service.TaskTrackerService
}

func New(trackerSvc service.TaskTrackerService) *TaskTrackersImpl {
	return &TaskTrackersImpl{trackerSvc: trackerSvc}
}

func (t *TaskTrackersImpl) Register(srv grpc.ServiceRegistrar) {
	pb.RegisterTaskTrackersAPIServer(srv, t)
}

func (t *TaskTrackersImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := runtime.NewServeMux()

	err := pb.RegisterTaskTrackersAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering task trackers grpc-gateway handler")
	}

	return "/api/task-trackers/", gwMux
}
