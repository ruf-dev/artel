package tracts_api

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
	pb "github.com/ruf-dev/artel/internal/api/server/artel_api"
	"github.com/ruf-dev/artel/internal/service"
	"github.com/ruf-dev/artel/internal/transport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// TractsImpl implements pb.TractsAPIServer. RunTract fires the engine as
// `go s.tractSvc.StartRun(s.baseCtx, ...)` — baseCtx is the server-lifecycle context (set at
// construction, see internal/app/custom.go), not the per-request context, which is cancelled
// the moment the RPC returns.
type TractsImpl struct {
	pb.UnimplementedTractsAPIServer
	tractSvc service.TractService
	baseCtx  context.Context
}

func New(baseCtx context.Context, tractSvc service.TractService) *TractsImpl {
	return &TractsImpl{
		tractSvc: tractSvc,
		baseCtx:  baseCtx,
	}
}

func (t *TractsImpl) Register(srv grpc.ServiceRegistrar) {
	pb.RegisterTractsAPIServer(srv, t)
}

func (t *TractsImpl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (string, http.Handler) {
	gwMux := transport.NewGatewayMux()

	err := pb.RegisterTractsAPIHandlerFromEndpoint(ctx, gwMux, endpoint, opts)
	if err != nil {
		log.Error().Err(err).Msg("error registering tracts grpc-gateway handler")
	}

	return "/api/tracts/", gwMux
}

// webhookBaseURL derives the scheme+host to prefix a webhook path with, from the grpc-gateway
// metadata forwarded for this call. Best-effort: grpc-gateway does not guarantee "x-forwarded-
// host"/":authority" survive to the backend for every deployment topology (reverse proxies
// vary), so an empty string falls back to a host-relative path.
func webhookBaseURL(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	host := firstMetadataValue(md, "x-forwarded-host")
	if host == "" {
		host = firstMetadataValue(md, ":authority")
	}

	if host == "" {
		return ""
	}

	scheme := firstMetadataValue(md, "x-forwarded-proto")
	if scheme == "" {
		scheme = "https"
	}

	return scheme + "://" + host
}

func firstMetadataValue(md metadata.MD, key string) string {
	values := md.Get(key)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}
