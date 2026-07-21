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
	internalapp "github.com/ruf-dev/artel/internal/app"
)

// App is artel's application instance: config, data sources, network
// listeners, and wired services/transports. It is an alias for
// internal/app.App, so its exported fields and methods (Start, ...) are
// usable from outside this module without importing internal/app directly.
type App = internalapp.App

// New builds an App: loads config, opens data sources and network listeners,
// and wires all services/transports. Call Start on the result to run it.
func New() (App, error) {
	return internalapp.New()
}
