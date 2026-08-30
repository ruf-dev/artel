// Package app is the public entry point for embedding artel as a library.
//
// internal/app cannot be imported outside this module (Go's internal-package
// rule), so an external repo that wants to build artel from source without
// forking it depends on github.com/ruf-dev/artel and imports this package
// instead:
//
//	import "github.com/ruf-dev/artel/pkg/app"
//
//	func main() {
//		a, err := app.New()
//		if err != nil {
//			log.Fatal().Err(err).Msg("Failed to create application")
//		}
//
//		err = a.Start()
//		if err != nil {
//			log.Fatal().Err(err).Msg("Failed to start application")
//		}
//	}
//
// The resulting binary still needs the config/ directory (or -config) next
// to it at runtime, same as cmd/service.
package app

import (
	"github.com/rs/zerolog/log"
	"go.redsock.ru/rerrors"

	internalapp "github.com/ruf-dev/artel/internal/app"
	svcv1 "github.com/ruf-dev/artel/internal/service/v1"
	"github.com/ruf-dev/artel/pkg/momhttp"
)

// App is artel's application instance: config, data sources, network
// listeners, and wired services/transports. It is an alias for
// internal/app.App, so its exported fields and methods (Start, ...) are
// usable from outside this module without importing internal/app directly.
type App = internalapp.App

// Option configures the App built by New.
type Option func(*options)

type options struct {
	svcOpts []svcv1.Option
}

// WithToolHttpMiddleware scopes an http.RoundTripper middleware to a single MoM tool
// (momhttp.McpTool*, e.g. momhttp.McpToolTrelloCreateCard) on the MoM http executor.
func WithToolHttpMiddleware(tool momhttp.McpTool, mw momhttp.Middleware) Option {
	return func(o *options) {
		o.svcOpts = append(o.svcOpts, svcv1.WithMomToolHttpMiddleware(tool, mw))
	}
}

// WithMcpHttpMiddleware scopes an http.RoundTripper middleware to every tool of one MoM
// (momhttp.McpName*, e.g. momhttp.McpNameTrello) on the MoM http executor.
func WithMcpHttpMiddleware(mcp momhttp.McpName, mw momhttp.Middleware) Option {
	return func(o *options) {
		o.svcOpts = append(o.svcOpts, svcv1.WithMomMcpHttpMiddleware(mcp, mw))
	}
}

// New builds an App: loads config, opens data sources and network listeners, wires all
// services/transports, and applies opts to the MoM http executor wiring. Call Start on the
// result to run it.
//
// It inlines internal/app.New's sequence (InitConfig -> InitServers -> Custom.Init) rather than
// delegating to it, so it can set Custom.SvcOptions from opts before Custom.Init wires the
// services. The generated internal/app.New() (DO NOT EDIT, no args) stays authoritative for the
// no-option path — cmd/service still calls it — and is simply unused here. app.New() with no
// options is behaviourally identical to it.
func New(opts ...Option) (App, error) {
	log.Info().Msg("starting app")

	o := &options{}

	for _, opt := range opts {
		opt(o)
	}

	var a App

	err := a.InitConfig()
	if err != nil {
		return App{}, rerrors.Wrap(err, "error initializing config")
	}

	err = a.InitServers()
	if err != nil {
		return App{}, rerrors.Wrap(err, "error during network listeners initialization")
	}

	a.Custom.SvcOptions = o.svcOpts

	err = a.Custom.Init(&a)
	if err != nil {
		return App{}, rerrors.Wrap(err, "error initializing custom app properties")
	}

	return a, nil
}
