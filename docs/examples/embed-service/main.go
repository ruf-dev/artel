// Command embed-service demonstrates embedding artel as a library inside a
// separate Go module. See README.md in this directory for the full writeup,
// and docs/embedding.md in the artel repo for the general guide.
package main

import (
	"log"
	"net/http"

	"github.com/ruf-dev/artel/pkg/app"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	registerExampleHandler(a)

	err = a.Start()
	if err != nil {
		log.Fatalf("failed to start app: %v", err)
	}
}

// registerExampleHandler wires this example's own HTTP handler onto the
// embedded app's transport, on top of everything artel registers itself
// (see internal/app/custom.go). It must be called after app.New() returns
// (that's when a.Custom.Transport is populated) and before a.Start() is
// called (that's when the transport actually starts accepting connections).
func registerExampleHandler(a app.App) {
	a.Custom.Transport.AddHttpHandler("/example/hello", http.HandlerFunc(helloHandler))
}

// helloHandler is exercised by both main() (as the real embedded handler)
// and main_test.go (as the thing under test), so there is exactly one
// implementation of the behavior being demonstrated/tested.
func helloHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"message":"hello from the embedded example"}`))
}
