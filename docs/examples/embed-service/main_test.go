package main

import (
	"io"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/ruf-dev/artel/pkg/app"
	"go.redsock.ru/toolbox/closer"
)

// TestEmbeddedApp proves the full embedding lifecycle described in
// README.md: app.New() builds and wires an App from this directory's
// config/config.yaml plus env-var overrides, registerExampleHandler (the
// exact function main() calls) adds an extra HTTP route on top of artel's
// own, a.Start() actually serves it, and closer.Close() shuts everything
// down cleanly.
//
// This is the only test in the package/binary: config.Init() (called by
// app.New()) is a process-wide singleton that errors on a second call, so
// app.New() must run exactly once per test binary.
func TestEmbeddedApp(t *testing.T) {
	t.Setenv("DATA_SOURCES_POSTGRES_HOST", "localhost")
	t.Setenv("DATA_SOURCES_POSTGRES_PORT", "15436")
	t.Setenv("DATA_SOURCES_POSTGRES_USER", "artel")
	t.Setenv("DATA_SOURCES_POSTGRES_PWD", "")
	t.Setenv("DATA_SOURCES_POSTGRES_DB_NAME", "artel_db")
	t.Setenv("DATA_SOURCES_POSTGRES_SSL_MODE", "disable")
	t.Setenv("DATA_SOURCES_POSTGRES_MIGRATIONS_FOLDER", "../../../migrations")

	// Port 0 asks the OS for a free port so this test doesn't collide with
	// anything else already listening; the real port is read back from
	// a.MASTER.Addr() below.
	t.Setenv("SERVERS_MASTER_PORT", "0")

	t.Setenv("ENVIRONMENT_CREDS_ENCRYPTION_KEY", "")
	t.Setenv("ENVIRONMENT_GOOGLE_API_KEY", "")
	t.Setenv("ENVIRONMENT_GOOGLE_CLIENT_ID", "")
	t.Setenv("ENVIRONMENT_GOOGLE_CLIENT_SECRET", "")
	t.Setenv("ENVIRONMENT_LOG_FORMAT", "TEXT")
	t.Setenv("ENVIRONMENT_LOG_LEVEL", "Info")
	t.Setenv("ENVIRONMENT_NO_AUTH_ENABLED", "true")
	t.Setenv("ENVIRONMENT_OTEL_ENDPOINT", "")
	t.Setenv("ENVIRONMENT_SUBSCRIPTIONS_ENABLED", "false")
	t.Setenv("ENVIRONMENT_TELEGRAM_CLIENT_ID", "")

	a, err := app.New()
	if err != nil {
		t.Fatalf("app.New() failed (is Postgres running? try: docker compose up -d): %v", err)
	}

	registerExampleHandler(a)

	startErrC := make(chan error, 1)
	go func() {
		startErrC <- a.Start()
	}()

	tcpAddr, ok := a.MASTER.Addr().(*net.TCPAddr)
	if !ok {
		t.Fatalf("expected *net.TCPAddr, got %T", a.MASTER.Addr())
	}
	url := "http://127.0.0.1:" + strconv.Itoa(tcpAddr.Port) + "/example/hello"

	body, err := pollUntilServed(url, 10*time.Second)
	if err != nil {
		t.Fatalf("embedded handler never became reachable at %s: %v", url, err)
	}

	const want = `{"message":"hello from the embedded example"}`
	if body != want {
		t.Fatalf("unexpected response body: got %q, want %q", body, want)
	}

	err = closer.Close()
	if err != nil {
		t.Logf("closer.Close() returned an error during shutdown: %v", err)
	}

	select {
	case startErr := <-startErrC:
		if startErr != nil {
			t.Fatalf("a.Start() returned an error after shutdown: %v", startErr)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("a.Start() did not return within 10s of closer.Close()")
	}
}

// pollUntilServed retries the GET against url until it succeeds or timeout
// elapses. The embedded transport only starts accepting connections once
// a.Start() actually runs (not when app.New() returns), so the caller can't
// assume the server is ready the instant the goroutine is launched.
func pollUntilServed(url string, timeout time.Duration) (string, error) {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: time.Second}

	var lastErr error
	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err != nil {
			lastErr = err
			time.Sleep(100 * time.Millisecond)
			continue
		}

		defer resp.Body.Close()

		data, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return "", readErr
		}

		return string(data), nil
	}

	return "", lastErr
}
