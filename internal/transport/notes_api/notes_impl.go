package notes_api

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/transport"
	"google.golang.org/grpc"
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
	gwMux := transport.NewGatewayMux()

	err := pb.RegisterNotesAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering notes grpc-gateway handler")
	}

	return "/api/notes/", gwMux
}
