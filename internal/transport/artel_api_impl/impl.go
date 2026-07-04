package artel_api_impl

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"github.com/ruf-dev/artel/internal/api/server/artel_api"
	"google.golang.org/grpc"
)

type Impl struct {
	artel_api.UnimplementedVaultsAPIServer
}

func New() *Impl {
	return &Impl{}
}

func (impl *Impl) Register(server grpc.ServiceRegistrar) {
	artel_api.RegisterVaultsAPIServer(server, impl)
}

func (impl *Impl) Gateway(
	ctx context.Context,
	endpoint string,
	opts ...grpc.DialOption,
) (route string, handler http.Handler) {
	gwHttpMux := runtime.NewServeMux()

	err := artel_api.RegisterVaultsAPIHandlerFromEndpoint(
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
