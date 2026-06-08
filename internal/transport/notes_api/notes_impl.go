package notes_api

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
)

type NotesImpl struct {
	pb.UnimplementedNotesAPIServer
	noteSvc service.NotesService
}

func NewNotesImpl(noteSvc service.NotesService) *NotesImpl {
	return &NotesImpl{noteSvc: noteSvc}
}

func (n *NotesImpl) Register(srv grpc.ServiceRegistrar) {
	pb.RegisterNotesAPIServer(srv, n)
}

func (n *NotesImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := runtime.NewServeMux()

	err := pb.RegisterNotesAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering notes grpc-gateway handler")
	}

	return "/api/notes/", gwMux
}
