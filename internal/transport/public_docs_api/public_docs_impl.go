package public_docs_api

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/transport"
	"google.golang.org/grpc"
)

type PublicDocsImpl struct {
	pb.UnimplementedPublicDocsAPIServer
	publicDocsSvc service.PublicDocsService
	cookieSecure  bool
}

func New(publicDocsSvc service.PublicDocsService, cookieSecure bool) *PublicDocsImpl {
	return &PublicDocsImpl{publicDocsSvc: publicDocsSvc, cookieSecure: cookieSecure}
}

func (p *PublicDocsImpl) Register(srv grpc.ServiceRegistrar) {
	pb.RegisterPublicDocsAPIServer(srv, p)
}

func (p *PublicDocsImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := transport.NewGatewayMux(p.cookieSecure)

	err := pb.RegisterPublicDocsAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering public docs grpc-gateway handler")
	}

	return "/api/public-docs/", gwMux
}
