package couchdb

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// transientUnauthorizedMaxAttempts and transientUnauthorizedRetryDelay bound loggingTransport's
// retry of a 401 response — see RoundTrip's doc comment for why every request through this client
// needs this, not just the ones that happened to be observed failing.
const (
	transientUnauthorizedMaxAttempts = 5
	transientUnauthorizedRetryDelay  = 300 * time.Millisecond
)

type loggingTransport struct {
	next http.RoundTripper
}

func newLoggingTransport() *loggingTransport {
	return &loggingTransport{next: otelhttp.NewTransport(http.DefaultTransport)}
}

// RoundTrip retries a 401 response up to transientUnauthorizedMaxAttempts times. CouchDB can
// answer any authenticated request — not just a specific one, every call site through this client
// hits the same shared _users-backed admin session — with a transient "Unauthorized: Name or
// password is incorrect" while admin credentials it recently wrote (via /_cluster_setup's
// enable_single_node, or a concurrent _users write) are still propagating internally; a retry
// shortly after succeeds against the identical request. Retrying here, once, covers every request
// this client makes instead of needing the same retry loop reimplemented at each call site.
func (t *loggingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	start := time.Now()

	var reqBody []byte
	if req.Body != nil {
		reqBody, _ = io.ReadAll(req.Body)
		req.Body = io.NopCloser(bytes.NewReader(reqBody))
	}

	var resp *http.Response
	var err error

	for attempt := 0; attempt < transientUnauthorizedMaxAttempts; attempt++ {
		if req.Body != nil {
			req.Body = io.NopCloser(bytes.NewReader(reqBody))
		}

		resp, err = t.next.RoundTrip(req)
		if err != nil || resp.StatusCode != http.StatusUnauthorized {
			break
		}

		if attempt < transientUnauthorizedMaxAttempts-1 {
			time.Sleep(transientUnauthorizedRetryDelay)
		}
	}

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
