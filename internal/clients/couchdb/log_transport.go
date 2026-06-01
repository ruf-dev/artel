package couchdb

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type loggingTransport struct {
	next http.RoundTripper
}

func newLoggingTransport() *loggingTransport {
	return &loggingTransport{next: otelhttp.NewTransport(http.DefaultTransport)}
}

func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	var reqBody []byte
	if req.Body != nil {
		reqBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(reqBody))
	}

	resp, err := t.next.RoundTrip(req)
	elapsed := time.Since(start)

	if err != nil {
		log.Ctx(req.Context()).Error().Err(err).
			Str("method", req.Method).
			Str("url", req.URL.String()).
			Bytes("req_body", reqBody).
			Dur("dur", elapsed).
			Msg("couch db returned error")
		return resp, err
	}

	var respBody []byte
	if resp.Body != nil {
		respBody, _ = io.ReadAll(resp.Body)
		resp.Body = io.NopCloser(bytes.NewReader(respBody))
	}

	log.Ctx(req.Context()).Debug().
		Str("method", req.Method).
		Str("url", req.URL.String()).
		Int("status", resp.StatusCode).
		Bytes("req_body", reqBody).
		Bytes("resp_body", respBody).
		Dur("dur", elapsed).
		Msg("couch db request")

	return resp, nil
}
