package trello

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var httpClient = &http.Client{Transport: newLoggingTransport()}

type loggingTransport struct {
	next http.RoundTripper
}

func newLoggingTransport() *loggingTransport {
	return &loggingTransport{next: otelhttp.NewTransport(http.DefaultTransport)}
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := t.next.RoundTrip(req)
	elapsed := time.Since(start)

	if err != nil {
		log.Ctx(req.Context()).
			Error().
			Err(err).
			Str("method", req.Method).
			Str("url", req.URL.String()).
			Dur("dur", elapsed).
			Msg("trello returned error")

		return resp, err
	}

	log.Ctx(req.Context()).
		Debug().
		Str("method", req.Method).
		Str("url", req.URL.String()).
		Int("status", resp.StatusCode).
		Dur("dur", elapsed).
		Msg("trello request")

	return resp, nil
}
