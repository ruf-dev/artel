package couchdb

import (
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type loggingTransport struct {
	next http.RoundTripper
}

func newLoggingTransport() *loggingTransport {
	return &loggingTransport{next: http.DefaultTransport}
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()
	resp, err := t.next.RoundTrip(req)
	elapsed := time.Since(start)
	if err != nil {
		log.Error().Err(err).Str("method", req.Method).Str("url", req.URL.String()).Dur("dur", elapsed).Msg("couch db returned error")
		return resp, err
	}
	log.Debug().Str("method", req.Method).Str("url", req.URL.String()).Int("status", resp.StatusCode).Dur("dur", elapsed).Msg("couch db request")
	return resp, nil
}
