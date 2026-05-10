package artel_api_impl

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	artel_api "github.com/ruf-dev/artel/internal/api/server/artel_api"
)

type Impl struct {
    artel_api.UnimplementedArtelAPIServer
}

func New() *Impl {
    return &Impl{}
}

func (impl *Impl) Register(server grpc.ServiceRegistrar) {
    artel_api.RegisterArtelAPIServer(server, impl)
}

func (impl *Impl) Gateway(ctx context.Context, endpoint string, opts ...grpc.DialOption) (route string, handler http.Handler) {
    gwHttpMux := runtime.NewServeMux()

    err := artel_api.RegisterArtelAPIHandlerFromEndpoint(
        ctx,
        gwHttpMux,
        endpoint,
        opts,
    )
    if err != nil {
        log.Error().Err(err).Msg("error registering grpc2http handler")
    }

    return "/api/", gwHttpMux
}
